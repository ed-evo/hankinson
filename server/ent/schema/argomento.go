package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Argomento holds the schema definition for the Argomento entity.
type Argomento struct {
	ent.Schema
}

func (Argomento) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "argomenti",
		},
	}
}

// Fields of the Argomento.
func (Argomento) Fields() []ent.Field {
	autoincremeent := false
	return []ent.Field{
		field.Int("id").
			Unique().
			Annotations(entsql.Annotation{
				Incremental: &autoincremeent,
			}),
		field.String("nome").NotEmpty().Unique(),
	}
}

// Edges of the Argomento.
func (Argomento) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("domande", Domanda.Type).
			StorageKey(
				edge.Table("argomenti_domande"),
				edge.Columns("argomento_id", "domanda_id"),
			),
	}
}

func (Argomento) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("nome"),
	}
}
