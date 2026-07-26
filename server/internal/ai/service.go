package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"sync"

	"github.com/ed-evo/hankinson/server/ent"
	"github.com/ed-evo/hankinson/server/internal/config"
	"github.com/ed-evo/hankinson/server/webui"
	"google.golang.org/genai"
)

type AiService interface {
	Spiega(context.Context, *ent.Domanda) (*QuizSpiegazione, error)
}

type QuizSpiegazione struct {
	ID               int    `json:"id"`
	Spiegazione      string `json:"spiegazione"`
	FocusLinguistico string `json:"focus_linguistico"`
	RegolaChiave     string `json:"regola_chiave"`
}

var quizSpiegazioneSchema = &genai.Schema{
	Type:        genai.TypeObject,
	Description: "Spiegazione per la domanda dell'esame patente B",
	Required:    []string{"id", "spiegazione"},
	Properties: map[string]*genai.Schema{
		"id":                {Type: genai.TypeInteger, Description: "Il numero identificativo della domanda"},
		"spiegazione":       {Type: genai.TypeString, Description: "Spiegazione semplice e chiara del perché è VERO/FALSO"},
		"focus_linguistico": {Type: genai.TypeString, Description: "Chiarimento sulle parole trappola (es. 'di norma')"},
		"regola_chiave":     {Type: genai.TypeString, Description: "Pillola mnemonica brevissima"},
	},
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
var systemPromt = genai.Text(`Sei un insegnante di scuola guida. Ti verrà fornito un gruppo di domande legate a una stessa immagine o argomento.
Per OGNI domanda presente nell'elenco, genera una spiegazione semplice (livello B1), evidenzia le parole trappola e fornisci una regola chiave mnemonica. Assicurati di includere il numero identificativo per ciascuna domanda.`,
)[0]

var promptTemplate = template.Must(template.New("quiz").Parse(`Ecco la domanda del blocco da analizzare:

- ID domanda: {{.ID}}
  Testo: {{.Testo}}
  Rsisposta Corretta: {{if .IsTrue}}VERO{{else}}FALSO{{end}}
`))

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

func compileDomandaPart(d *ent.Domanda) (*genai.Part, error) {
	var buf bytes.Buffer
	if err := promptTemplate.Execute(&buf, d); err != nil {
		return nil, fmt.Errorf("Errore creazione prompt: %w", err)
	}

	return &genai.Part{Text: buf.String()}, nil
}

func (ai *Gem) Spiega(ctx context.Context, d *ent.Domanda) (*QuizSpiegazione, error) {
	var buf bytes.Buffer
	if err := promptTemplate.Execute(&buf, d); err != nil {
		return nil, fmt.Errorf("Errore creazione prompt: %w", err)
	}

	aiCfg := &genai.GenerateContentConfig{
		ResponseMIMEType:  "application/json",
		ResponseSchema:    quizSpiegazioneSchema,
		SystemInstruction: systemPromt,
	}

	var parts []*genai.Part

	if i, err := toImagePart(d); err == nil {
		if i != nil {
			parts = append(parts, i)
		}
	} else {
		return nil, fmt.Errorf("Error compiling image %w", err)
	}

	if p, err := compileDomandaPart(d); err == nil {
		parts = append(parts, p)
	} else {
		return nil, fmt.Errorf("Error compiling prompt: %w", err)
	}

	resp, err := ai.client.Models.GenerateContent(ctx, ai.model, []*genai.Content{{Parts: parts}}, aiCfg)

	if err != nil {
		return nil, fmt.Errorf("gemini generation error: %w", err)
	}

	var spiegazione QuizSpiegazione
	if err := json.Unmarshal([]byte(resp.Text()), &spiegazione); err != nil {
		return nil, fmt.Errorf("Errore recupero Spiegazione da ai: %w", err)
	}

	return &spiegazione, nil
}

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
