package controllers

import (
	"github.com/Efren-Garza-Z/go-api-gemini/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var DB *gorm.DB

func SetDB(db *gorm.DB) {
	DB = db
}

// GET /users/:id
// @Summary Obtener un usuario por ID
// @Description Devuelve la información de un usuario por su ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "ID del usuario"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/{id} [get]
func GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var userDB models.UserDB
	if err := DB.First(&userDB, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error en la base de datos"})
		}
		return
	}

	c.JSON(http.StatusOK, userDB.ToSwagger())
}

// POST /users
// @Summary Crear un usuario
// @Description Crea un usuario nuevo con contraseña
// @Tags users
// @Accept json
// @Produce json
// @Param input body models.CreateUserInput true "Datos para crear usuario"
// @Success 201 {object} models.User
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users [post]
func CreateUser(c *gin.Context) {
	var input models.CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userDB := models.UserDB{
		FullName: input.FullName,
		Email:    input.Email,
		Password: input.Password, // Aquí la contraseña se guarda en DB (idealmente con hash)
	}

	if err := DB.Create(&userDB).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear el usuario"})
		return
	}

	c.JSON(http.StatusCreated, userDB.ToSwagger())
}
