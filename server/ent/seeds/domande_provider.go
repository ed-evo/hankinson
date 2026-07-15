package seeds

import (
	"embed"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io/fs"
	"log"
	"path/filepath"
	"slices"
)

//go:embed capitoli.json
var CapitoliJson []byte

//go:embed domande/*
var DomandeJsonSeedsFs embed.FS

//go:embed seeds_domande.template.sql
var DomandeSqlTemplate string

type capitoloInfo struct {
	ID               int    `json:"capitolo"`
	Nome             string `json:"nome"`
	MinNumeroDomanda int    `json:"min"`
	MaxNumeroDomanda int    `json:"max"`
	TotaleDomande    int    `json:"atteso"`
}

type rawQuizItem struct {
	IDBlocco         int      `json:"id_blocco"`
	Argomenti        []string `json:"argomenti"`
	Numero           int      `json:"numero"`
	Testo            string   `json:"testo"`
	RispostaCorretta string   `json:"risposta_corretta"`
	Immagine         string   `json:"immagine"`
	PaginaQuiz       int      `json:"pagina_quiz"`
}

type Capitolo struct {
	ID               int
	Nome             string
	MinNumeroDomanda int
	MaxNumeroDomanda int
	TotaleDomande    int
}
type Domanda struct {
	IDCapitolo int
	IDBlocco   int
	Argomenti  []int
	Numero     int
	Testo      string
	IsTrue     bool
	Immagine   string
	PaginaQuiz int
}

type Argomento struct {
	ID   int
	Nome string
}

// hashTopic matches your deterministic 32-bit FNV-1a logic
func hashTopic(name string) int {
	h := fnv.New32a()
	_, _ = h.Write([]byte(name))
	return int(h.Sum32())
}
func GetSeeds() ([]Capitolo, []Domanda, []Argomento, error) {

	var capitoliInfo []capitoloInfo
	if err := json.Unmarshal(CapitoliJson, &capitoliInfo); err != nil {
		log.Printf("Error reading file %v", err)
		return nil, nil, nil, err
	}
	// a causa del capitolo 23 che ha domande anche con numero maggiore dei capitoli 24 e 25
	// bisogna riordina per fare in modo che nel ciclo di assegnazione capitolo
	// 24 e 25 abbiamo priorita' rispetto al capitolo 23 per evitare
	// che tutte le domande dei capitoli 23, 24 e 25 finiscano nel 23
	slices.SortFunc(capitoliInfo, func(a, b capitoloInfo) int {
		return a.MaxNumeroDomanda - b.MaxNumeroDomanda
	})

	domandeDir := "domande"

	files, err := fs.ReadDir(DomandeJsonSeedsFs, domandeDir)
	if err != nil {
		log.Printf("Errore lettura file da cartella %v")
		return nil, nil, nil, err
	}

	domande := make(map[int]Domanda)
	argomenti := make(map[int]Argomento)

	for _, f := range files {
		if f.IsDir() || filepath.Ext(f.Name()) != ".json" {
			continue
		}
		path := filepath.Join("domande", f.Name())
		data, err := DomandeJsonSeedsFs.ReadFile(path)
		if err != nil {
			log.Printf("Error reading file %v", err)
			return nil, nil, nil, err
		}

		var raws []rawQuizItem
		if err := json.Unmarshal(data, &raws); err != nil {
			log.Printf("Error reading file %v", err)
			return nil, nil, nil, err
		}

		for _, raw := range raws {

			log.Printf("raw data %v", raw)

			var IDCapitolo int

			for _, c := range capitoliInfo {
				if c.MinNumeroDomanda <= raw.Numero && raw.Numero <= c.MaxNumeroDomanda {
					IDCapitolo = c.ID
					break
				}
			}

			if IDCapitolo == 0 {
				return nil, nil, nil, fmt.Errorf("Capitolo non trovato per domanda %v", raw.Numero)
			}

			var hashedTopicIDs []int
			for _, topicName := range raw.Argomenti {
				topicID := hashTopic(topicName)
				hashedTopicIDs = append(hashedTopicIDs, topicID)

				if _, ok := argomenti[topicID]; !ok {
					argomenti[topicID] = Argomento{ID: topicID, Nome: topicName}
				}
			}

			if d, ok := domande[raw.Numero]; ok {
				return nil, nil, nil, fmt.Errorf("Errore domanda DUPLICATO %v", d)
			} else {
				domande[raw.Numero] = Domanda{
					IDBlocco:   raw.IDBlocco,
					Argomenti:  hashedTopicIDs,
					Numero:     raw.Numero,
					Testo:      raw.Testo,
					IsTrue:     raw.RispostaCorretta == "VERO",
					Immagine:   raw.Immagine,
					IDCapitolo: IDCapitolo,
					PaginaQuiz: raw.PaginaQuiz,
				}
			}
		}

	}

	var listaCapitoli []Capitolo
	for _, c := range capitoliInfo {
		listaCapitoli = append(listaCapitoli, Capitolo{
			ID:               c.ID,
			Nome:             c.Nome,
			MinNumeroDomanda: c.MinNumeroDomanda,
			MaxNumeroDomanda: c.MaxNumeroDomanda,
			TotaleDomande:    c.TotaleDomande,
		})
	}

	var listaDomande []Domanda
	for _, d := range domande {
		listaDomande = append(listaDomande, d)
	}
	var listaArgomenti []Argomento
	for _, a := range argomenti {
		listaArgomenti = append(listaArgomenti, a)
	}

	return listaCapitoli, listaDomande, listaArgomenti, nil
}
