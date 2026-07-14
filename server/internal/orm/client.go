package orm

import (
	"context"
	"database/sql"
	"log"
	"sync"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/ed-evo/hankinson/server/ent"
	"github.com/ed-evo/hankinson/server/ent/seeds"

	_ "modernc.org/sqlite"
)

var (
	client *ent.Client
	once   sync.Once
)

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

		err = seeds.SeedDomande(ctx, db)
	})

	return client
}
