package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// Domanda holds the schema definition for the Domanda entity.
type Utente struct {
	ent.Schema
}

func (Utente) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "utenti",
		},
	}
}

func (Utente) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			StorageKey("email").
			NotEmpty().
			Unique(),
		field.Time("created_at").
			Default(time.Now),
	}
}
