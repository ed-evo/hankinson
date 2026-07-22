package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Domanda holds the schema definition for the Domanda entity.
type AttivitaQuesitoEsame struct {
	ent.Schema
}

func (AttivitaQuesitoEsame) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "attivita_quesito_esame",
		},
	}
}

func (AttivitaQuesitoEsame) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("tipo").
			Values("salta", "risposta", "pausa", "prossimo"),
		field.Bool("risposta_data").
			Optional().
			Nillable().
			Immutable(),
		field.Time("inizio").Immutable(),
		field.Int("durata_ms"),
		field.Time("timestamp").
			Immutable().
			Default(time.Now),
	}
}

func (AttivitaQuesitoEsame) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("quesito_esame", QuesitoEsame.Type).
			Ref("attivita").
			Unique().
			Required(),
	}
}
