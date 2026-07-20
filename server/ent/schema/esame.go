package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Domanda holds the schema definition for the Domanda entity.
type Esame struct {
	ent.Schema
}

func (Esame) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "esami",
		},
	}
}

func (Esame) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("tipo").
			Values("ministeriale", "parziale", "aperto"),
		field.Int("numero_quesiti"),
		field.Int("max_errori").
			Default(3),
		field.Int("minuti_disponibili").
			Default(30),
		field.Time("created_at").
			Default(time.Now),
		field.Time("updated_at").
			Default(time.Now),
	}
}

func (Esame) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("utente", Utente.Type).
			Unique().
			Required(),
		edge.To("quesiti", QuesitoEsame.Type),
	}
}

func (Esame) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("tipo").
			Edges("utente").
			Unique().
			Annotations(
				entsql.IndexWhere("tipo = 'aperto'"),
			),
	}
}
