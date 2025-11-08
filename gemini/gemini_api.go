package gemini

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/genai"
)

// Service GeminiService es una abstracción del servicio de la API de Gemini.
type Service struct{}

// newClient inicializa el cliente de Gemini con la configuración estándar.
func newClient() (*genai.Client, context.Context, error) {
	_ = godotenv.Load() // carga .env si existe

	ctx := context.Background()

	// Validación de configuración
	if os.Getenv("GOOGLE_GENAI_USE_VERTEXAI") != "true" && os.Getenv("GEMINI_API_KEY") == "" {
		return nil, nil, fmt.Errorf("falta configurar GOOGLE_API_KEY o habilitar VertexAI (GOOGLE_GENAI_USE_VERTEXAI=true)")
	}

	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("error creando cliente genai: %w", err)
	}

	return client, ctx, nil
}

// GenerateContent genera contenido a partir de un prompt de texto.
func (s *Service) GenerateContent(prompt string) (string, error) {
	client, ctx, err := newClient()
	if err != nil {
		return "", err
	}

	cfg := &genai.GenerateContentConfig{
		Temperature: genai.Ptr[float32](0.5),
	}

	chat, err := client.Chats.Create(ctx, "gemini-2.0-flash", cfg, nil)
	if err != nil {
		return "", fmt.Errorf("error creando chat: %w", err)
	}

	res, err := chat.SendMessage(ctx, genai.Part{Text: prompt})
	if err != nil {
		return "", fmt.Errorf("error enviando mensaje: %w", err)
	}

	return res.Text(), nil
}

// GenerateWithFile genera contenido usando un prompt y un archivo adjunto de cualquier tipo.
func (s *Service) GenerateWithFile(file io.Reader, filename string, mimeType string, prompt string) (string, error) {
	client, ctx, err := newClient()
	if err != nil {
		return "", err
	}

	// Subir el archivo, pasando el MIMEType dinámicamente
	f, err := client.Files.Upload(ctx, file, &genai.UploadFileConfig{
		DisplayName: filename,
		MIMEType:    mimeType, // Ahora el MIMEType es dinámico
	})
	if err != nil {
		return "", fmt.Errorf("upload: %w", err)
	}

	chat, err := client.Chats.Create(ctx, "gemini-2.0-flash", nil, nil)
	if err != nil {
		return "", err
	}

	parts := []genai.Part{
		{Text: prompt},
		{FileData: &genai.FileData{
			FileURI: func() string {
				if f.URI != "" {
					return f.URI
				}
				return f.Name
			}(),
			MIMEType: f.MIMEType,
		}},
	}

	res, err := chat.SendMessage(ctx, parts...)
	if err != nil {
		return "", err
	}
	return res.Text(), nil
}
