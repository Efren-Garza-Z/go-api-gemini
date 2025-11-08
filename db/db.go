package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "edgz"
	}
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "1234"
	}
	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "edgz"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		host, user, password, dbname, port)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}

	if err := DB.Exec("CREATE SCHEMA IF NOT EXISTS service;").Error; err != nil {
		log.Fatalf("Error al crear schema: %v", err)
	}

	log.Println("Conexión a la base de datos exitosa")
}
