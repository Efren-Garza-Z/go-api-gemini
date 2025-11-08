package models

import "time"

// GeminiProcessingFileDB guarda el archivo (bytea) junto con su metadata.
// Lo dejamos en la DB tal como lo tenías (campo File []byte).
type GeminiProcessingFileDB struct {
	ID        string                 `gorm:"primaryKey" json:"id" example:"8b9a1d2e-3c4f-5a6b-7c8d-9e0f1a2b3c4d"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
	Status    GeminiProcessingStatus `gorm:"type:varchar(20);not null" json:"status" example:"pendiente"`
	Result    string                 `gorm:"type:text" json:"result,omitempty"`
	Error     string                 `gorm:"type:text" json:"error,omitempty"`
	Prompt    string                 `gorm:"type:text;not null" json:"prompt"`
	File      []byte                 `gorm:"type:bytea" json:"-"`
	Filename  string                 `gorm:"type:varchar(255)" json:"filename,omitempty"`
	MimeType  string                 `gorm:"type:varchar(100)" json:"mime_type,omitempty"`
}

func (GeminiProcessingFileDB) TableName() string {
	return "service.gemini_processing_file"
}

// Response DTO
type GeminiProcessingFileIDResponse struct {
	GeminiProcessingFileID string `json:"task_id" example:"8b9a1d2e-3c4f-5a6b-7c8d-9e0f1a2b3c4d"`
}

type GeminiProcessingFileResponse struct {
	ID     string                 `json:"id"`
	Status GeminiProcessingStatus `json:"status"`
	Result string                 `json:"result,omitempty"`
	Error  string                 `json:"error,omitempty"`
}
