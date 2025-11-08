package models

import "time"

// PromptRequest es la estructura de la solicitud para iniciar una tarea.
type PromptRequest struct {
	Prompt string `json:"prompt" example:"Conoces las becas para poder estudiar en finlandia o noruega?"`
}

// GeminiProcessingIDResponse es la respuesta que contiene el ID de la tarea.
type GeminiProcessingIDResponse struct {
	GeminiProcessingID string `json:"task_id" example:"8b9a1d2e-3c4f-5a6b-7c8d-9e0f1a2b3c4d"`
}

// GeminiProcessingStatus representa el estado de una tarea asíncrona.
type GeminiProcessingStatus string

const (
	// StatusPending indica que la tarea está pendiente de iniciar.
	StatusPending GeminiProcessingStatus = "pendiente"
	// StatusProcessing indica que la tarea está en curso.
	StatusProcessing GeminiProcessingStatus = "en_proceso"
	// StatusCompleted indica que la tarea ha finalizado con éxito.
	StatusCompleted GeminiProcessingStatus = "finalizado"
	// StatusError indica que la tarea ha fallado.
	StatusError GeminiProcessingStatus = "error"
)

// GeminiProcessingResponse representa el estado y el resultado de una tarea.
type GeminiProcessingResponse struct {
	ID     string                 `json:"id" example:"8b9a1d2e-3c4f-5a6b-7c8d-9e0f1a2b3c4d"`
	Status GeminiProcessingStatus `json:"status" example:"finalizado"`
	Result string                 `json:"result,omitempty" example:"Sí, existen varias becas..."`
	Error  string                 `json:"error,omitempty"`
}

// GeminiProcessingDB es el modelo que se guarda en la base de datos.
type GeminiProcessingDB struct {
	ID        string `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Status    GeminiProcessingStatus `gorm:"type:varchar(20);not null"`
	Result    string                 `gorm:"type:text"`
	Error     string                 `gorm:"type:text"`
	Prompt    string                 `gorm:"type:text;not null"`
}

// TableName especifica el nombre de la tabla en la DB.
func (GeminiProcessingDB) TableName() string {
	return "gemini.gemini_processing"
}
