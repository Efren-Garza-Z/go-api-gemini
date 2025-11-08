package routes

import (
	"github.com/Efren-Garza-Z/go-api-gemini/web/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.Engine, uc *controllers.UserController) {
	users := r.Group("/users")
	{
		users.POST("", uc.CreateUser)
		users.GET("", uc.GetAll)
		users.GET("/:id", uc.GetByID)
		users.PUT("/:id", uc.Update)
		users.DELETE("/:id", uc.Delete)
	}
}
