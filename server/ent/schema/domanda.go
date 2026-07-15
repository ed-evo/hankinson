package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Domanda holds the schema definition for the Domanda entity.
type Domanda struct {
	ent.Schema
}

func (Domanda) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "domande",
		},
	}
}

// Fields of the Domanda.
func (Domanda) Fields() []ent.Field {
	autoincrement := false
	return []ent.Field{
		field.Int("id").
			StorageKey("numero").
			Unique().
			Annotations(entsql.Annotation{
				Incremental: &autoincrement,
			}),
		field.Text("testo"),
		field.Bool("is_true"),
		field.String("immagine").
			Optional().
			Nillable(),
		field.Int("id_capitolo"),
		field.Int("pagina_quiz"),
		field.Int("id_blocco"),
	}
}

func (Domanda) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("argomenti", Argomento.Type).
			Ref("domande"),
		edge.From("capitolo", Capitolo.Type).
			Ref("domande").
			Unique(). // Garantisce che sia Many-to-One
			Field("id_capitolo").
			Required(),
	}
}

func (Domanda) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("id_capitolo"),
		index.Fields("is_true"),
		index.Fields("id_blocco"),
		index.Fields("pagina_quiz"),
	}
}
