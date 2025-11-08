package main

import (
	"log"
	"os"

	"github.com/Efren-Garza-Z/go-api-gemini/db"
	_ "github.com/Efren-Garza-Z/go-api-gemini/docs" // swag docs
	"github.com/Efren-Garza-Z/go-api-gemini/domain/models"
	"github.com/Efren-Garza-Z/go-api-gemini/domain/repositories"
	service "github.com/Efren-Garza-Z/go-api-gemini/services"
	controllers "github.com/Efren-Garza-Z/go-api-gemini/web/controllers"
	"github.com/Efren-Garza-Z/go-api-gemini/web/routes"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// carga .env si usas godotenv en runtime
	// _ = godotenv.Load()

	db.Connect()

	// Migraciones
	if err := db.DB.AutoMigrate(&models.UserDB{}, &models.GeminiProcessingDB{}, &models.GeminiProcessingFileDB{}); err != nil {
		log.Fatalf("Error al migrar modelos: %v", err)
	}

	// Repositorios
	userRepo := repositories.NewUserRepository(db.DB)
	gemRepo := repositories.NewGeminiRepository(db.DB)

	// Services
	userSvc := service.NewUserService(userRepo)
	gemSvc := service.NewGeminiService(gemRepo)

	// Controllers
	userCtrl := controllers.NewUserController(userSvc, db.DB)
	gemCtrl := controllers.NewGeminiController(gemSvc)

	// Gin
	r := gin.Default()

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Routes
	routes.RegisterUserRoutes(r, userCtrl)
	routes.RegisterGeminiRoutes(r, gemCtrl)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Servidor corriendo en http://localhost:%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Error al iniciar servidor: %v", err)
	}
}
