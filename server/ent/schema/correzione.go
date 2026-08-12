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

type Correzione struct {
	ent.Schema
}

func (Correzione) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "correzioni",
		},
	}
}

func (Correzione) Fields() []ent.Field {
	return []ent.Field{
		field.Int("esame_id").Immutable(),
		field.Enum("type").
			Values("human", "ai"),
		field.Text("esaminatore"),
		field.Bool("is_promosso").
            StructTag(`json:"is_promosso"`),
		field.Text("testo"),
		field.Text("meta").Optional(),
		field.Time("created_at").
			Default(time.Now),
	}
}

func (Correzione) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("esame", Esame.Type).
			Ref("correzioni").
			Unique().
			Field("esame_id").
			Required().
			Immutable(),
	}
}

func (Correzione) Indexes() []ent.Index {
	return []ent.Index{
		// Enforces unique constraint via a unique index
		index.Fields("esaminatore", "esame_id").Unique(),
	}
}
