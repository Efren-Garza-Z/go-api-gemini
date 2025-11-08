package routes

import (
	"github.com/Efren-Garza-Z/go-api-gemini/web/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterGeminiRoutes(r *gin.Engine, gc *controllers.GeminiController) {
	g := r.Group("/gemini")
	{
		g.POST("/process", gc.ProcessPrompt)
		g.GET("/status/:gemini_processing_id", gc.GetTaskStatus)

		g.POST("/process-file", gc.ProcessFile)
		g.GET("/status-file/:gemini_processing_id", gc.GetFileStatus)
	}
}
