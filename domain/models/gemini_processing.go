package models

import "time"

// GeminiProcessingStatus tipo para los estados de la tarea
type GeminiProcessingStatus string

const (
	StatusPending    GeminiProcessingStatus = "pendiente"
	StatusProcessing GeminiProcessingStatus = "en_proceso"
	StatusCompleted  GeminiProcessingStatus = "finalizado"
	StatusError      GeminiProcessingStatus = "error"
)

// GeminiProcessingDB es el modelo que se guarda en la DB (tabla service.gemini_processing)
type GeminiProcessingDB struct {
	ID        string                 `gorm:"primaryKey" json:"id" example:"8b9a1d2e-3c4f-5a6b-7c8d-9e0f1a2b3c4d"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
	Status    GeminiProcessingStatus `gorm:"type:varchar(20);not null" json:"status" example:"pendiente"`
	Result    string                 `gorm:"type:text" json:"result,omitempty" example:"Resultado del modelo"`
	Error     string                 `gorm:"type:text" json:"error,omitempty"`
	Prompt    string                 `gorm:"type:text;not null" json:"prompt" example:"Qué es Go?"`
}

func (GeminiProcessingDB) TableName() string {
	return "service.gemini_processing"
}

// DTOs para Swagger / responses
type PromptRequest struct {
	Prompt string `json:"prompt" example:"Conoces las becas para Finlandia?" binding:"required"`
}

type GeminiProcessingIDResponse struct {
	GeminiProcessingID string `json:"task_id" example:"8b9a1d2e-3c4f-5a6b-7c8d-9e0f1a2b3c4d"`
}

type GeminiProcessingResponse struct {
	ID     string                 `json:"id" example:"8b9a1d2e-3c4f-5a6b-7c8d-9e0f1a2b3c4d"`
	Status GeminiProcessingStatus `json:"status" example:"finalizado"`
	Result string                 `json:"result,omitempty" example:"Sí, existen varias becas..."`
	Error  string                 `json:"error,omitempty"`
}
