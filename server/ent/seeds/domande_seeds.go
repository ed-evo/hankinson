package seeds

import (
	"context"
	"database/sql"
	_ "embed"
	"log"
)

//go:generate go run generators/domande_generate.go
//go:embed seeds_domande.sql
var domandeSqlSeeds string

func SeedDomande(ctx context.Context, db *sql.DB) error {

	_, err := db.ExecContext(ctx, domandeSqlSeeds)

	if err != nil {
		log.Printf("Errore esecuzione seeds %v", err)
		return err
	}

	log.Print(">>>>> Seed Domande Completato <<<<<<<<")

	return nil
}
