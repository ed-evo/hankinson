package seeds

import (
	"context"
	"embed"
	"encoding/json"
	"log"
	"path/filepath"
	"strings"

	"github.com/ed-evo/hankinson/server/ent"
)

//go:embed domande/*
var domandeJsonSeedsFs embed.FS

type rawQuizItem struct {
	IDBlocco         int      `json:"id_blocco"`
	Argomenti        []string `json:"argomenti"`
	Numero           int      `json:"numero"`
	Testo            string   `json:"testo"`
	RispostaCorretta string   `json:"risposta_corretta"`
	Immagine         string   `json:"immagine"`
	PaginaQuiz       int      `json:"pagina_quiz"`
}

func SeedDomande(ctx context.Context, client *ent.Client) {
	seedsDir := "domande"
	files, err := domandeJsonSeedsFs.ReadDir(seedsDir)
	if err != nil {
		log.Fatalf("Failed to read seeds directory: %v", err)
	}

	tx, err := client.Tx(ctx)
	if err != nil {
		log.Fatalf("Failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	argomentiCache := make(map[string]int)

	existingArgomenti, err := tx.Argomento.Query().All(ctx)
	if err != nil {
		log.Fatalf("failed to query existing argomenti: %v", err)
	}
	for _, arg := range existingArgomenti {
		argomentiCache[arg.Nome] = arg.ID
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasPrefix(file.Name(), "quesito_") || !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		filePath := filepath.Join(seedsDir, file.Name())

		data, err := domandeJsonSeedsFs.ReadFile(filePath)
		if err != nil {
			log.Fatalf("❌ Failed to read file %s: %v\n", file.Name(), err)
			continue
		}

		var items []rawQuizItem
		if err := json.Unmarshal(data, &items); err != nil {
			log.Fatalf("❌ Failed to parse JSON in %s: %v\n", file.Name(), err)
			continue
		}
		for _, item := range items {
			log.Printf("Domand %v", item)
			var topicIDs []int

			for _, topicName := range item.Argomenti {
				id, exists := argomentiCache[topicName]
				if !exists {
					topicObj, err := tx.Argomento.Create().
						SetNome(topicName).
						Save(ctx)
					if err == nil {
						id = topicObj.ID
						argomentiCache[topicName] = id
					} else {
						log.Fatalf("❌ Failed to create topic '%s': %v\n", topicName, err)
						continue
					}
				}
				topicIDs = append(topicIDs, id)
			}
			log.Printf("topics %v", topicIDs)
			domandaId, err := tx.Domanda.Create().
				SetNumero(item.Numero).
				SetTesto(item.Testo).
				SetIsTrue(item.RispostaCorretta == "VERO").
				SetNillableImmagine(&item.Immagine).
				SetPaginaQuiz(item.PaginaQuiz).
				SetIDBlocco(item.IDBlocco).
				OnConflict().
				UpdateNewValues().
				ID(ctx)
			if err != nil {
				log.Fatalf("❌ Failed to upser %v", item)
			}

			err = tx.Domanda.
				UpdateOneID(domandaId).
				ClearArgomenti().
				AddArgomentiIDs(topicIDs...).
				Exec(ctx)

			if err != nil {
				log.Fatalf("❌ Error updating argomenti %v.", item)
			}

			log.Printf("Upsert Domand id: %v", domandaId)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Fatalf("Failed to commit transaction: %v", err)
	}
	log.Println("🎉 Seeding finished instantly!")
}
