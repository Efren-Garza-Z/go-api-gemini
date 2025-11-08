package main

import (
	"log"

	"github.com/Efren-Garza-Z/go-api-gemini/controllers"
	"github.com/Efren-Garza-Z/go-api-gemini/db"
	_ "github.com/Efren-Garza-Z/go-api-gemini/docs" // Importa la documentación Swagger generada
	"github.com/Efren-Garza-Z/go-api-gemini/models"
	"github.com/Efren-Garza-Z/go-api-gemini/routes"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title API GEMINI
// @version 1.0
// @description API RESTful para gestión de usuarios
// @host localhost:8080
// @BasePath /
func main() {

	// Conexión a PostgreSQL
	db.Connect()

	// Migración automática del modelo UserDB
	if err := db.DB.AutoMigrate(&models.UserDB{}, &models.GeminiProcessingDB{}, &models.GeminiProcessingFileDB{}); err != nil {
		log.Fatalf("Error al migrar modelo UserDB: %v", err)
	}

	// Inyectar la conexión DB en los controladores
	controllers.SetDB(db.DB)

	// Crear instancia de Gin
	r := gin.Default()

	// Rutas Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Rutas para usuarios
	routes.RegisterUserRoutes(r)

	// Iniciar servidor
	log.Println("Servidor corriendo en http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Error al iniciar servidor: %v", err)
	}
}
