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
type QuesitoEsame struct {
	ent.Schema
}

func (QuesitoEsame) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "quesiti_esame",
		},
	}
}

func (QuesitoEsame) Fields() []ent.Field {
	return []ent.Field{
		field.Bool("risposta_finale").
			Optional(),
		field.Time("created_at").
			Default(time.Now),
		field.Time("updated_at").
			Default(time.Now),
	}
}

func (QuesitoEsame) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("esame", Esame.Type).
			Ref("quesiti").
			Unique().
			Required(),
		edge.To("domanda_originale", Domanda.Type).
			Unique().
			Required(),
		edge.To("logs", AttivitaQuesitoEsame.Type),
	}
}
