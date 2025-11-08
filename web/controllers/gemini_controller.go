package controllers

import (
	"io"
	"net/http"

	"github.com/Efren-Garza-Z/go-api-gemini/domain/models"
	"github.com/Efren-Garza-Z/go-api-gemini/services"
	"github.com/gin-gonic/gin"
)

type GeminiController struct {
	service services.GeminiService
}

func NewGeminiController(s services.GeminiService) *GeminiController {
	return &GeminiController{service: s}
}

// @Summary Iniciar procesamiento de prompt
// @Tags gemini
// @Accept json
// @Produce json
// @Param requestBody body models.PromptRequest true "Prompt a procesar"
// @Success 202 {object} models.GeminiProcessingIDResponse
// @Failure 400 {object} map[string]string
// @Router /gemini/process [post]
func (gc *GeminiController) ProcessPrompt(c *gin.Context) {
	var req models.PromptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}
	id, err := gc.service.ProcessPromptAsync(req.Prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo iniciar proceso"})
		return
	}
	c.JSON(http.StatusAccepted, models.GeminiProcessingIDResponse{GeminiProcessingID: id})
}

// @Summary Obtener estado de procesamiento
// @Tags gemini
// @Produce json
// @Param gemini_processing_id path string true "ID del proceso"
// @Success 200 {object} models.GeminiProcessingResponse
// @Router /gemini/status/{gemini_processing_id} [get]
func (gc *GeminiController) GetTaskStatus(c *gin.Context) {
	id := c.Param("gemini_processing_id")
	p, err := gc.service.GetProcessStatus(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Proceso no encontrado"})
		return
	}
	resp := models.GeminiProcessingResponse{
		ID:     p.ID,
		Status: p.Status,
		Result: p.Result,
		Error:  p.Error,
	}
	c.JSON(http.StatusOK, resp)
}

// @Summary Iniciar procesamiento con archivo
// @Tags gemini
// @Accept multipart/form-data
// @Produce json
// @Param prompt formData string true "Prompt"
// @Param file formData file true "Archivo (pdf/png/jpg)"
// @Success 202 {object} models.GeminiProcessingFileIDResponse
// @Failure 400 {object} map[string]string
// @Router /gemini/process-file [post]
func (gc *GeminiController) ProcessFile(c *gin.Context) {
	prompt := c.PostForm("prompt")
	if prompt == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Prompt requerido"})
		return
	}
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Archivo requerido"})
		return
	}
	f, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo abrir archivo"})
		return
	}
	defer f.Close()
	content, err := io.ReadAll(f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo leer archivo"})
		return
	}
	id, err := gc.service.ProcessFileAsync(prompt, fileHeader.Filename, fileHeader.Header.Get("Content-Type"), content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo iniciar procesamiento de archivo"})
		return
	}
	c.JSON(http.StatusAccepted, models.GeminiProcessingFileIDResponse{GeminiProcessingFileID: id})
}

// @Summary Obtener estado de procesamiento de archivo
// @Tags gemini
// @Produce json
// @Param gemini_processing_id path string true "ID del proceso"
// @Success 200 {object} models.GeminiProcessingFileResponse
// @Router /gemini/status-file/{gemini_processing_id} [get]
func (gc *GeminiController) GetFileStatus(c *gin.Context) {
	id := c.Param("gemini_processing_id")
	f, err := gc.service.GetFileProcessStatus(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Proceso no encontrado"})
		return
	}
	resp := models.GeminiProcessingFileResponse{
		ID:     f.ID,
		Status: f.Status,
		Result: f.Result,
		Error:  f.Error,
	}
	c.JSON(http.StatusOK, resp)
}
