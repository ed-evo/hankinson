package orm

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/ed-evo/hankinson/server/ent"
	"github.com/ed-evo/hankinson/server/ent/domanda"
	"github.com/ed-evo/hankinson/server/ent/seeds"

	_ "modernc.org/sqlite"
)

var (
	client *ent.Client
	once   sync.Once
)

var DomandeByCapitolo map[int][]int = make(map[int][]int)

func GetClient(ctx context.Context) *ent.Client {
	once.Do(func() {
		// 1. Open the database using standard Go code (accepts "sqlite")
		db, err := sql.Open("sqlite", "file:hankinson.db?cache=shared&_pragma=journal_mode(WAL)&_pragma=foreign_keys(1)")
		db.SetMaxOpenConns(1)
		if err != nil {
			log.Fatalf("failed opening connection to sqlite: %v", err)
		}

		// 2. Wrap the connection. This tells Ent: "Use this active DB pool,
		//    but format all the SQL queries using the standard SQLite dialect syntax."
		drv := entsql.OpenDB(dialect.SQLite, db)

		// 3. Create your Ent client
		client = ent.NewClient(ent.Driver(drv))

		if err := client.Schema.Create(ctx); err != nil {
			log.Fatalf("failed creating schema resources: %v", err)
		}

		err = seeds.SeedDomande(ctx, client)

		if err != nil {
			log.Fatalf("Errore Seeding Domande: %v", err)
		}

		err = seeds.VerificaCapitoli(ctx, client)

		if err != nil {
			log.Fatalf("Errore Capitoli: %v", err)
		}

		records, err := client.Domanda.Query().
			Select(domanda.FieldID, domanda.FieldIDCapitolo).
			All(ctx)

		if err != nil {
			log.Fatalf("Error reading Domande %v", err)
		}

		for _, d := range records {
			DomandeByCapitolo[d.IDCapitolo] = append(DomandeByCapitolo[d.IDCapitolo], d.ID)
		}

	})

	return client
}

func WithTx(ctx context.Context, client *ent.Client, fn func(tx *ent.Tx) error) error {
	tx, err := client.Tx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if v := recover(); v != nil {
			tx.Rollback()
			panic(v)
		}
	}()
	if err := fn(tx); err != nil {
		log.Printf("Errore esecuzione query %v", err)
		if rerr := tx.Rollback(); rerr != nil {
			err = fmt.Errorf("%w: rolling back transaction: %v", err, rerr)
		}
		return err
	}
	if err := tx.Commit(); err != nil {
		log.Printf("Errore commit transation %v", err)
		return fmt.Errorf("committing transaction: %w", err)
	}
	return nil
}
