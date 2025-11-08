package services

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/Efren-Garza-Z/go-api-gemini/domain/models"
	"github.com/Efren-Garza-Z/go-api-gemini/domain/repositories"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	genai "google.golang.org/genai"
)

// GeminiService coordina repo + llamada a Gemini
type GeminiService interface {
	ProcessPromptAsync(prompt string) (string, error)
	GetProcessStatus(id string) (*models.GeminiProcessingDB, error)

	ProcessFileAsync(prompt, filename, mimeType string, fileContent []byte) (string, error)
	GetFileProcessStatus(id string) (*models.GeminiProcessingFileDB, error)
}

type geminiService struct {
	repo repositories.GeminiRepository
}

func NewGeminiService(r repositories.GeminiRepository) GeminiService {
	return &geminiService{repo: r}
}

// newClient crea un cliente Gemini usando la variable de entorno GEMINI_API_KEY
func newClient(ctx context.Context) (*genai.Client, context.Context, error) {
	// Cargar variables de entorno
	_ = godotenv.Load()

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, nil, fmt.Errorf("GEMINI_API_KEY no configurada en el entorno")
	}

	// Crear configuración del cliente
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("error creando cliente Gemini: %w", err)
	}

	return client, ctx, nil
}

// GenerateContent llama al modelo Gemini con texto (sin archivos)
func (s *geminiService) GenerateContent(prompt string) (string, error) {
	ctx := context.Background()
	client, _, err := newClient(ctx)
	if err != nil {
		return "", err
	}

	// Modelo base de Gemini
	model := "gemini-2.5-flash"

	chat, err := client.Chats.Create(ctx, model, &genai.GenerateContentConfig{
		Temperature: genai.Ptr[float32](0.5),
	}, nil)
	if err != nil {
		return "", fmt.Errorf("error creando chat: %w", err)
	}

	res, err := chat.SendMessage(ctx, genai.Part{Text: prompt})
	if err != nil {
		return "", fmt.Errorf("error enviando mensaje: %w", err)
	}

	return res.Text(), nil
}

// GenerateWithFile llama al modelo Gemini subiendo un archivo
func (s *geminiService) GenerateWithFile(prompt string, fileReader io.Reader, filename, mimeType string) (string, error) {
	ctx := context.Background()
	client, _, err := newClient(ctx)
	if err != nil {
		return "", err
	}

	// Subir archivo
	f, err := client.Files.Upload(ctx, fileReader, &genai.UploadFileConfig{
		DisplayName: filename,
		MIMEType:    mimeType,
	})
	if err != nil {
		return "", fmt.Errorf("error subiendo archivo: %w", err)
	}

	// Crear chat con el modelo Gemini
	model := "gemini-2.5-flash"
	chat, err := client.Chats.Create(ctx, model, nil, nil)
	if err != nil {
		return "", fmt.Errorf("error creando chat: %w", err)
	}

	parts := []genai.Part{
		{Text: prompt},
		{FileData: &genai.FileData{
			FileURI:  f.URI,
			MIMEType: f.MIMEType,
		}},
	}

	res, err := chat.SendMessage(ctx, parts...)
	if err != nil {
		return "", fmt.Errorf("error enviando mensaje con archivo: %w", err)
	}

	return res.Text(), nil
}

// ProcessPromptAsync crea registro y lanza goroutine para procesamiento de texto
func (s *geminiService) ProcessPromptAsync(prompt string) (string, error) {
	id := genUUID()

	proc := &models.GeminiProcessingDB{
		ID:     id,
		Status: models.StatusPending,
		Prompt: prompt,
	}
	if err := s.repo.CreateProcess(proc); err != nil {
		return "", err
	}

	go func(procID, p string) {
		_ = s.repo.UpdateStatus(procID, models.StatusProcessing, "", "")

		result, err := s.GenerateContent(p)
		if err != nil {
			_ = s.repo.UpdateStatus(procID, models.StatusError, "", err.Error())
			return
		}
		_ = s.repo.UpdateStatus(procID, models.StatusCompleted, result, "")
	}(id, prompt)

	return id, nil
}

// ProcessFileAsync crea registro y lanza goroutine para procesamiento con archivo
func (s *geminiService) ProcessFileAsync(prompt, filename, mimeType string, fileContent []byte) (string, error) {
	id := genUUID()

	proc := &models.GeminiProcessingFileDB{
		ID:       id,
		Status:   models.StatusPending,
		Prompt:   prompt,
		File:     fileContent,
		Filename: filename,
		MimeType: mimeType,
	}
	if err := s.repo.CreateFileProcess(proc); err != nil {
		return "", err
	}

	go func(procID, fname, mtype, p string, content []byte) {
		_ = s.repo.UpdateFileStatus(procID, models.StatusProcessing, "", "")

		result, err := s.GenerateWithFile(p, bytes.NewReader(content), fname, mtype)
		if err != nil {
			_ = s.repo.UpdateFileStatus(procID, models.StatusError, "", err.Error())
			return
		}
		_ = s.repo.UpdateFileStatus(procID, models.StatusCompleted, result, "")
	}(id, filename, mimeType, prompt, fileContent)

	return id, nil
}

func (s *geminiService) GetProcessStatus(id string) (*models.GeminiProcessingDB, error) {
	return s.repo.FindProcessByID(id)
}

func (s *geminiService) GetFileProcessStatus(id string) (*models.GeminiProcessingFileDB, error) {
	return s.repo.FindFileProcessByID(id)
}

// genUUID crea un identificador pseudo-único (usa uuid real en producción)
func genUUID() string {
	return uuid.New().String()
}
