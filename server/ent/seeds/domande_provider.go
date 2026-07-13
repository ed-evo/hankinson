package seeds

import (
	"embed"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io/fs"
	"log"
	"path/filepath"
)

//go:embed domande/*
var DomandeJsonSeedsFs embed.FS

//go:embed seeds_domande.template.sql
var DomandeSqlTemplate string

type rawQuizItem struct {
	IDBlocco         int      `json:"id_blocco"`
	Argomenti        []string `json:"argomenti"`
	Numero           int      `json:"numero"`
	Testo            string   `json:"testo"`
	RispostaCorretta string   `json:"risposta_corretta"`
	Immagine         string   `json:"immagine"`
	PaginaQuiz       int      `json:"pagina_quiz"`
}

type Domanda struct {
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
func GetSeeds() ([]Domanda, []Argomento, error) {

	domandeDir := "domande"

	files, err := fs.ReadDir(DomandeJsonSeedsFs, domandeDir)
	if err != nil {
		log.Printf("Errore lettura file da cartella %v")
		return nil, nil, err
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
			return nil, nil, err
		}

		var raws []rawQuizItem
		if err := json.Unmarshal(data, &raws); err != nil {
			log.Printf("Error reading file %v", err)
			return nil, nil, err
		}

		for _, raw := range raws {

			log.Printf("raw data %v", raw)

			var hashedTopicIDs []int
			for _, topicName := range raw.Argomenti {
				topicID := hashTopic(topicName)
				hashedTopicIDs = append(hashedTopicIDs, topicID)

				if _, ok := argomenti[topicID]; !ok {
					argomenti[topicID] = Argomento{ID: topicID, Nome: topicName}
				}
			}

			if d, ok := domande[raw.Numero]; ok {
				return nil, nil, fmt.Errorf("Errore domanda DUPLICATO %v", d)

			} else {
				domande[raw.Numero] = Domanda{
					IDBlocco:   raw.IDBlocco,
					Argomenti:  hashedTopicIDs,
					Numero:     raw.Numero,
					Testo:      raw.Testo,
					IsTrue:     raw.RispostaCorretta == "VERO",
					Immagine:   raw.Immagine,
					PaginaQuiz: raw.PaginaQuiz,
				}
			}
		}

	}

	var listaDomande []Domanda
	for _, d := range domande {
		listaDomande = append(listaDomande, d)
	}
	var listaArgomenti []Argomento
	for _, a := range argomenti {
		listaArgomenti = append(listaArgomenti, a)
	}

	return listaDomande, listaArgomenti, nil
}
