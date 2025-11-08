package models

import (
	"time"
)

// GeminiProcessingFileIDResponse es la respuesta que contiene el ID de la tarea.
type GeminiProcessingFileIDResponse struct {
	GeminiProcessingFileID string `json:"task_id" example:"8b9a1d2e-3c4f-5a6b-7c8d-9e0f1a2b3c4d"`
}

// GeminiProcessingFileResponse representa el estado y el resultado de una tarea.
type GeminiProcessingFileResponse struct {
	ID     string                 `json:"id" example:"8b9a1d2e-3c4f-5a6b-7c8d-9e0f1a2b3c4d"`
	Status GeminiProcessingStatus `json:"status" example:"finalizado"`
	Result string                 `json:"result,omitempty" example:"Sí, existen varias becas..."`
	Error  string                 `json:"error,omitempty"`
}

// GeminiProcessingFileDB es el modelo que se guarda en la base de datos.
type GeminiProcessingFileDB struct {
	ID        string `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Status    GeminiProcessingStatus `gorm:"type:varchar(20);not null"`
	Result    string                 `gorm:"type:text"`
	Error     string                 `gorm:"type:text"`
	Prompt    string                 `gorm:"type:text;not null"`
	File      []byte                 `gorm:"type:bytea"`
}

// TableName especifica el nombre de la tabla en la DB.
func (GeminiProcessingFileDB) TableName() string {
	return "gemini.gemini_processing_file"
}
