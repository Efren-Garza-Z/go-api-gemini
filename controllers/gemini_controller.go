package controllers

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/Efren-Garza-Z/go-api-gemini/db"
	"github.com/Efren-Garza-Z/go-api-gemini/gemini"
	"github.com/Efren-Garza-Z/go-api-gemini/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "google.golang.org/api/option"
)

// In-memory store para las tareas. En producción, usarías una base de datos.
//var GeminiProcessing = make(map[string]models.GeminiProcessingResponse)
//var mu sync.Mutex // Mutex para evitar race conditions en el map

// ProcessPrompt @Summary Iniciar tarea asíncrona de Gemini
// @Description Inicia una tarea en segundo plano para procesar un prompt con la API de Gemini.
// @Tags gemini
// @Accept  json
// @Produce  json
// @Param   requestBody body models.PromptRequest true "Prompt a procesar"
// @Success 202 {object} models.GeminiProcessingIDResponse "Solicitud aceptada y procesando"
// @Failure 400 {object} map[string]string "JSON de solicitud inválido"
// @Router /gemini/process [post]
func ProcessPrompt(c *gin.Context) {
	var requestBody models.PromptRequest
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON de solicitud inválido"})
		return
	}

	GeminiProcessingID := uuid.New().String()

	// Crear registro inicial en DB
	newGeminiProcessing := models.GeminiProcessingDB{
		ID:     GeminiProcessingID,
		Status: models.StatusPending,
		Prompt: requestBody.Prompt,
	}
	if err := db.DB.Create(&newGeminiProcessing).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear la tarea"})
		return
	}

	// Procesar en segundo plano
	go func(id string, prompt string) {
		service := &gemini.Service{}

		// Actualizar estado a en_proceso
		db.DB.Model(&models.GeminiProcessingDB{}).
			Where("id = ?", id).
			Update("status", models.StatusProcessing)

		result, err := service.GenerateContent(prompt)
		if err != nil {
			db.DB.Model(&models.GeminiProcessingDB{}).
				Where("id = ?", id).
				Updates(map[string]interface{}{
					"status": models.StatusError,
					"error":  err.Error(),
				})
		} else {
			db.DB.Model(&models.GeminiProcessingDB{}).
				Where("id = ?", id).
				Updates(map[string]interface{}{
					"status": models.StatusCompleted,
					"result": result,
				})
		}
	}(GeminiProcessingID, requestBody.Prompt)

	c.JSON(http.StatusAccepted, models.GeminiProcessingIDResponse{GeminiProcessingID: GeminiProcessingID})
}

// GetTaskStatus @Summary Obtener estado de la tarea de Gemini
// @Description Consulta el estado y el resultado de una tarea por su ID.
// @Tags gemini
// @Accept  json
// @Produce  json
// @Param   gemini_processing_id path string true "ID del proceso"
// @Success 200 {object} models.GeminiProcessingResponse "Estado del proceso y resultado"
// @Failure 400 {object} map[string]string "ID de proceso inválido"
// @Router /gemini/status/{gemini_processing_id} [get]
func GetTaskStatus(c *gin.Context) {
	geminiProcessingID := c.Param("gemini_processing_id")
	var geminiProcessing models.GeminiProcessingDB

	if err := db.DB.First(&geminiProcessing, "id = ?", geminiProcessingID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ID de proceso no encontrado"})
		return
	}

	// Convertir TaskDB → TaskResponse para no exponer Prompt ni timestamps
	response := models.GeminiProcessingResponse{
		ID:     geminiProcessing.ID,
		Status: geminiProcessing.Status,
		Result: geminiProcessing.Result,
		Error:  geminiProcessing.Error,
	}
	c.JSON(http.StatusOK, response)
}

// GenerateWithFileController @Summary Generar contenido de Gemini con un archivo (asíncrono)
// @Description Procesa un prompt y un archivo de forma asíncrona, guarda los datos y retorna un ID de proceso.
// @Tags gemini
// @Accept  multipart/form-data
// @Produce  json
// @Param   prompt formData string true "Texto del prompt"
// @Param   file formData file true "Archivo (PDF, PNG, JPEG)"
// @Success 202 {object} models.GeminiProcessingFileIDResponse "Proceso en cola"
// @Failure 400 {object} map[string]string "Solicitud inválida"
// @Router /gemini/process/file [post]
func GenerateWithFileController(c *gin.Context) {
	prompt := c.PostForm("prompt")
	if prompt == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El prompt es obligatorio"})
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Se requiere un archivo"})
		return
	}

	fileType := fileHeader.Header.Get("Content-Type")
	switch fileType {
	case "image/jpeg", "image/png", "application/pdf":
		// Tipo de archivo válido
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Tipo de archivo no soportado: %s", fileType)})
		return
	}

	// Leer el contenido del archivo de forma segura en memoria
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo abrir el archivo"})
		return
	}
	fileContent, err := io.ReadAll(file)
	file.Close() // Cerrar el archivo después de leerlo
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo leer el archivo"})
		return
	}

	// 1. Crear un ID único y guardar el registro inicial con el archivo
	processID := uuid.New().String()
	geminiProcess := models.GeminiProcessingFileDB{
		ID:     processID,
		Status: models.StatusProcessing,
		Prompt: prompt,
		File:   fileContent, // Guardamos el contenido del archivo
	}

	if err := db.DB.Create(&geminiProcess).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar el registro en la base de datos"})
		return
	}

	// 2. Responder inmediatamente con el ID de la tarea
	c.JSON(http.StatusAccepted, models.GeminiProcessingFileIDResponse{
		GeminiProcessingFileID: processID,
	})

	// 3. Iniciar el procesamiento de Gemini en una goroutine
	go func() {
		service := &gemini.Service{}
		// Convertimos el slice de bytes a un io.Reader para la función GenerateWithFile.
		fileReader := bytes.NewReader(fileContent)

		result, err := service.GenerateWithFile(fileReader, fileHeader.Filename, fileType, prompt)

		// 4. Actualizar el registro en la base de datos
		if err != nil {
			log.Printf("Error procesando tarea %s con Gemini: %v", processID, err)
			db.DB.Model(&models.GeminiProcessingFileDB{}).
				Where("id = ?", processID).
				Updates(map[string]interface{}{
					"status": models.StatusError,
					"error":  err.Error(),
				})
		} else {
			db.DB.Model(&models.GeminiProcessingFileDB{}).
				Where("id = ?", processID).
				Updates(map[string]interface{}{
					"status": models.StatusCompleted,
					"result": result,
				})
		}
	}()
}

// GetGeminiProcessStatus @Summary Obtener estado de la tarea de Gemini
// @Description Consulta el estado y el resultado de una tarea por su ID.
// @Tags gemini
// @Accept  json
// @Produce  json
// @Param   gemini_processing_id path string true "ID del proceso"
// @Success 200 {object} models.GeminiProcessingFileResponse "Estado del proceso y resultado"
// @Failure 400 {object} map[string]string "ID de proceso inválido"
// @Router /gemini/status-file/{gemini_processing_id} [get]
func GetGeminiProcessStatus(c *gin.Context) {
	geminiProcessingFileID := c.Param("gemini_processing_id")
	var geminiProcessingFile models.GeminiProcessingFileDB

	if err := db.DB.First(&geminiProcessingFile, "id = ?", geminiProcessingFileID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ID de proceso no encontrado"})
		return
	}

	// Convertir TaskDB → TaskResponse para no exponer Prompt ni timestamps
	response := models.GeminiProcessingFileResponse{
		ID:     geminiProcessingFile.ID,
		Status: geminiProcessingFile.Status,
		Result: geminiProcessingFile.Result,
		Error:  geminiProcessingFile.Error,
	}
	c.JSON(http.StatusOK, response)
}
