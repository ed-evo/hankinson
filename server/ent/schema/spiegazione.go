package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type Spiegazione struct {
	ent.Schema
}

func (Spiegazione) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "spiegazioni",
		},
	}
}

func (Spiegazione) Fields() []ent.Field {
	return []ent.Field{
		field.Int("numero_domanda").
			Immutable(),
		field.String("spiegazione"),
		field.String("focus_linguistico"),
		field.String("regola_chiave"),
	}
}

func (Spiegazione) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("domanda", Domanda.Type).
			Ref("spiegazione").
			Unique().
			Field("numero_domanda").
			Required().
			Immutable(),
	}
}

func (Spiegazione) Indexes() []ent.Index {
	return []ent.Index{
		// Enforces unique constraint via a unique index
		index.Fields("numero_domanda").Unique(),
	}
}
