package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Argomento holds the schema definition for the Argomento entity.
type Argomento struct {
	ent.Schema
}

// Fields of the Argomento.
func (Argomento) Fields() []ent.Field {
	return []ent.Field{
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
