package routes

import (
	"github.com/Efren-Garza-Z/go-api-gemini/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.Engine) {
	users := r.Group("/users")
	gemini := r.Group("/gemini")
	{
		users.GET("/:id", controllers.GetUserByID)
		users.POST("", controllers.CreateUser)

		gemini.POST("/process", controllers.ProcessPrompt)
		gemini.GET("/status/:task_id", controllers.GetTaskStatus)
	}
}
