package repositories

import (
	"github.com/Efren-Garza-Z/go-api-gemini/domain/models"
	"gorm.io/gorm"
)

type GeminiRepository interface {
	CreateProcess(p *models.GeminiProcessingDB) error
	FindProcessByID(id string) (*models.GeminiProcessingDB, error)
	UpdateStatus(id string, status models.GeminiProcessingStatus, result string, processError string) error

	CreateFileProcess(f *models.GeminiProcessingFileDB) error
	FindFileProcessByID(id string) (*models.GeminiProcessingFileDB, error)
	UpdateFileStatus(id string, status models.GeminiProcessingStatus, result string, processError string) error
}

type geminiRepository struct {
	db *gorm.DB
}

func NewGeminiRepository(db *gorm.DB) GeminiRepository {
	return &geminiRepository{db: db}
}

func (r *geminiRepository) CreateProcess(p *models.GeminiProcessingDB) error {
	return r.db.Create(p).Error
}

func (r *geminiRepository) FindProcessByID(id string) (*models.GeminiProcessingDB, error) {
	var p models.GeminiProcessingDB
	if err := r.db.First(&p, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *geminiRepository) UpdateStatus(id string, status models.GeminiProcessingStatus, result string, processError string) error {
	updates := map[string]interface{}{"status": status}
	if result != "" {
		updates["result"] = result
	}
	if processError != "" {
		updates["error"] = processError
	}
	return r.db.Model(&models.GeminiProcessingDB{}).Where("id = ?", id).Updates(updates).Error
}

func (r *geminiRepository) CreateFileProcess(f *models.GeminiProcessingFileDB) error {
	return r.db.Create(f).Error
}

func (r *geminiRepository) FindFileProcessByID(id string) (*models.GeminiProcessingFileDB, error) {
	var f models.GeminiProcessingFileDB
	if err := r.db.First(&f, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *geminiRepository) UpdateFileStatus(id string, status models.GeminiProcessingStatus, result string, processError string) error {
	updates := map[string]interface{}{"status": status}
	if result != "" {
		updates["result"] = result
	}
	if processError != "" {
		updates["error"] = processError
	}
	return r.db.Model(&models.GeminiProcessingFileDB{}).Where("id = ?", id).Updates(updates).Error
}
