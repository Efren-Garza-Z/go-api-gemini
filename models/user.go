package models

import (
	"time"
)

// Modelo para la base de datos (GORM)
type UserDB struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	FullName string `gorm:"not null"`
	Password string `gorm:"not null"`
	Email    string `gorm:"uniqueIndex;not null"`
}

func (UserDB) TableName() string {
	return "gemini.users"
}

// Modelo para respuesta JSON y Swagger (no expone password)
type User struct {
	ID       uint   `json:"id" example:"1"`
	FullName string `json:"full_name" example:"Efren David"`
	Email    string `json:"email" example:"efren@example.com"`
}

// Modelo para recibir creación de usuario (input con password)
type CreateUserInput struct {
	FullName string `json:"full_name" binding:"required" example:"Efren David"`
	Email    string `json:"email" binding:"required,email" example:"efren@example.com"`
	Password string `json:"password" binding:"required" example:"miPasswordSeguro123"`
}

// Convierte UserDB a User para la respuesta
func (u *UserDB) ToSwagger() User {
	return User{
		ID:       u.ID,
		FullName: u.FullName,
		Email:    u.Email,
	}
}
