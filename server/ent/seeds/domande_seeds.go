package seeds

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"log"
)

//go:generate go run generators/domande_generate.go
//go:embed seeds_domande.sql
var domandeSqlSeeds string

type CapitoloVerifica struct {
	ID               int    `json:"capitolo_id"`
	Nome             string `json:"capitolo_nome"`
	MinNumeroDomanda int    `json:"min_numero_domanda"`
	MaxNumeroDomanda int    `json:"max_numero_domanda"`
	TotaleDomande    int    `json:"totale_domande"`
	ConteggioReale   int    `json:"conteggio_reale"`
	TotaleCorretto   bool   `json:"totale_corretto"`
}

func SeedDomande(ctx context.Context, db *sql.DB) error {

	_, err := db.ExecContext(ctx, domandeSqlSeeds)

	if err != nil {
		log.Printf("Errore esecuzione seeds %v", err)
		return err
	}

	log.Print(">>>>> Seed Domande Completato <<<<<<<<")

	return nil
}

func VerificaCapitoli(ctx context.Context, db *sql.DB) error {
	query := `
		SELECT 
			c.id,
			c.totale_domande,
			COUNT(d.numero) AS conteggio_reale
		FROM capitoli c
		LEFT JOIN domande d ON c.id = d.id_capitolo
		GROUP BY c.id
		ORDER BY c.id;
	`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return fmt.Errorf("errore nell'esecuzione della query di verifica: %w", err)
	}
	defer rows.Close()

	discrepanze := 0

	for rows.Next() {
		var id, totaleDomande, count int
		if err := rows.Scan(&id, &totaleDomande, &count); err != nil {
			return fmt.Errorf("errore durante lo scan dei dati: %w", err)
		}

		if totaleDomande != count {
			discrepanze++
			log.Printf("%v: Trovate %v domande per il Capitolo %v che deve averne %v", discrepanze, count, id, totaleDomande)
		}
	}

	if err = rows.Err(); err != nil {
		return fmt.Errorf("errore durante l'iterazione delle righe: %w", err)
	}

	if discrepanze > 0 {
		return fmt.Errorf("Ci sono %v Capitolo con numero di Domande non corrispondenti.", discrepanze)
	}

	return nil
}
