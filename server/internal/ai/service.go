package ai

import (
	"context"
	"fmt"
	"sync"

	"github.com/ed-evo/hankinson/server/ent"
	"github.com/ed-evo/hankinson/server/internal/config"
	"github.com/ed-evo/hankinson/server/webui"
	"google.golang.org/genai"
)

type AiService interface {
	Spiega(context.Context, *ent.Domanda) (*QuizSpiegazione, error)
	Correggi(context.Context, *ent.Client, int) (*EsameCorrezione, error)
}

type Gem struct {
	client *genai.Client
	model  string
}

var (
	gemini  *Gem
	once    sync.Once
	initErr error
)

func GetGemini(ctx context.Context) (AiService, error) {
	once.Do(func() {
		cfg := config.Get()

		client, err := genai.NewClient(ctx, &genai.ClientConfig{
			APIKey:  cfg.GeminiApiKey,
			Backend: genai.BackendGeminiAPI,
		})
		if err != nil {
			initErr = fmt.Errorf("failed to create Genai Client: %w", err)
			return
		}
		gemini = &Gem{
			client: client,
			model:  cfg.GeminiModel,
		}
	})

	return gemini, initErr
}

func toImagePart(d *ent.Domanda) (*genai.Part, error) {
	if d == nil {
		return nil, fmt.Errorf("Domanda required")
	}
	if d.Immagine == nil || *d.Immagine == "" {
		return nil, nil
	}
	data, err := webui.GetQuizImage(*d.Immagine)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving image: %w", err)
	}
	return &genai.Part{
		InlineData: &genai.Blob{
			MIMEType: "image/png",
			Data:     data,
		},
	}, nil
}
