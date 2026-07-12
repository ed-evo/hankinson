package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Domanda holds the schema definition for the Domanda entity.
type Domanda struct {
	ent.Schema
}

// Fields of the Domanda.
func (Domanda) Fields() []ent.Field {
	return []ent.Field{
		field.Int("numero").Unique(),
		field.Text("testo"),
		field.Bool("is_true"),
		field.String("immagine").
			Optional().
			Nillable(),
		field.Int("pagina_quiz"),
		field.Int("id_blocco"),
	}
}

func (Domanda) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("argomenti", Argomento.Type).
			Ref("domande"),
	}
}

func (Domanda) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("is_true"),
		index.Fields("id_blocco"),
		index.Fields("pagina_quiz"),
	}
}
