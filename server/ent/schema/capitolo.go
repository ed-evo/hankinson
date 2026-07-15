package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Capitolo holds the schema definition for the Capitolo entity.
type Capitolo struct {
	ent.Schema
}

func (Capitolo) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "capitoli",
		},
	}
}

// Fields of the Capitolo.
func (Capitolo) Fields() []ent.Field {
	autoincrement := false
	return []ent.Field{
		field.Int("id").
			Unique().
			Annotations(entsql.Annotation{
				Incremental: &autoincrement,
			}),
		field.String("nome").
			NotEmpty(),
		field.Int("min_numero_domanda"),
		field.Int("max_numero_domanda"),
		field.Int("totale_domande"),
	}
}

// Edges of the Capitolo.
func (Capitolo) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("domande", Domanda.Type),
	}
}
