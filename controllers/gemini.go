package controllers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/Efren-Garza-Z/go-api-gemini/db"
	"github.com/Efren-Garza-Z/go-api-gemini/models" // Importa el paquete models
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "google.golang.org/api/option"
	"google.golang.org/genai"
)

// In-memory store para las tareas. En producción, usarías una base de datos.
var tasks = make(map[string]models.TaskResponse)
var mu sync.Mutex // Mutex para evitar race conditions en el map

// GeminiService es una abstracción del servicio de la API de Gemini.
type GeminiService struct{}

func (s *GeminiService) GenerateContent(prompt string) (string, error) {
	// Carga variables del .env si existe (no revienta si no está)
	_ = godotenv.Load()

	ctx := context.Background()

	// Ayuda para depurar configuración:
	// - Para Gemini API: export GOOGLE_GENAI_USE_VERTEXAI=false y export GOOGLE_API_KEY=...
	// - Para Vertex AI:  export GOOGLE_GENAI_USE_VERTEXAI=true, GOOGLE_CLOUD_PROJECT y GOOGLE_CLOUD_LOCATION
	if os.Getenv("GOOGLE_GENAI_USE_VERTEXAI") != "true" && os.Getenv("GEMINI_API_KEY") == "" {
		return "", fmt.Errorf("falta configurar GOOGLE_API_KEY o habilitar VertexAI (GOOGLE_GENAI_USE_VERTEXAI=true)")
	}

	// La API nueva toma la configuración desde variables de entorno si pasas nil.
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("error creando cliente genai: %w", err)
	}
	// Nota: el cliente nuevo no expone Close()

	// Configuración de generación (opcional)
	cfg := &genai.GenerateContentConfig{
		Temperature: genai.Ptr[float32](0.5),
	}

	// Crea una sesión de chat con el modelo (usa el que prefieras).
	// Ejemplos oficiales recomiendan "gemini-2.0-flash".
	chat, err := client.Chats.Create(ctx, "gemini-2.0-flash", cfg, nil)
	if err != nil {
		return "", fmt.Errorf("error creando chat: %w", err)
	}

	// Envía el prompt (la API nueva usa genai.Part{Text: ...}, ya no genai.Text)
	res, err := chat.SendMessage(ctx, genai.Part{Text: prompt})
	if err != nil {
		return "", fmt.Errorf("error enviando mensaje: %w", err)
	}

	// Toma el texto de la respuesta directamente
	return res.Text(), nil
}

// @Summary Iniciar tarea asíncrona de Gemini
// @Description Inicia una tarea en segundo plano para procesar un prompt con la API de Gemini.
// @Tags gemini
// @Accept  json
// @Produce  json
// @Param   requestBody body models.PromptRequest true "Prompt a procesar"
// @Success 202 {object} models.TaskIDResponse "Solicitud aceptada y procesando"
// @Failure 400 {object} map[string]string "JSON de solicitud inválido"
// @Router /gemini/process [post]
func ProcessPrompt(c *gin.Context) {
	var requestBody models.PromptRequest
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON de solicitud inválido"})
		return
	}

	taskID := uuid.New().String()

	// Crear registro inicial en DB
	newTask := models.TaskDB{
		ID:     taskID,
		Status: models.StatusPending,
		Prompt: requestBody.Prompt,
	}
	if err := db.DB.Create(&newTask).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear la tarea"})
		return
	}

	// Procesar en segundo plano
	go func(id string, prompt string) {
		service := &GeminiService{}

		// Actualizar estado a en_proceso
		db.DB.Model(&models.TaskDB{}).
			Where("id = ?", id).
			Update("status", models.StatusProcessing)

		result, err := service.GenerateContent(prompt)
		if err != nil {
			db.DB.Model(&models.TaskDB{}).
				Where("id = ?", id).
				Updates(map[string]interface{}{
					"status": models.StatusError,
					"error":  err.Error(),
				})
		} else {
			db.DB.Model(&models.TaskDB{}).
				Where("id = ?", id).
				Updates(map[string]interface{}{
					"status": models.StatusCompleted,
					"result": result,
				})
		}
	}(taskID, requestBody.Prompt)

	c.JSON(http.StatusAccepted, models.TaskIDResponse{TaskID: taskID})
}

// @Summary Obtener estado de la tarea de Gemini
// @Description Consulta el estado y el resultado de una tarea por su ID.
// @Tags gemini
// @Accept  json
// @Produce  json
// @Param   task_id path string true "ID de la tarea"
// @Success 200 {object} models.TaskResponse "Estado de la tarea"
// @Failure 404 {object} map[string]string "ID de tarea no encontrado"
// @Router /gemini/status/{task_id} [get]
func GetTaskStatus(c *gin.Context) {
	taskID := c.Param("task_id")
	var task models.TaskDB

	if err := db.DB.First(&task, "id = ?", taskID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ID de tarea no encontrado"})
		return
	}

	// Convertir TaskDB → TaskResponse para no exponer Prompt ni timestamps
	response := models.TaskResponse{
		ID:     task.ID,
		Status: task.Status,
		Result: task.Result,
		Error:  task.Error,
	}
	c.JSON(http.StatusOK, response)
}
