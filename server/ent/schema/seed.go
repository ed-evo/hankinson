package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type Seed struct {
	ent.Schema
}

func (Seed) Fields() []ent.Field {
	return []ent.Field{
		field.String("hash").
			Unique().
			Immutable().NotEmpty(),
		field.Time("created_at").
			Default(time.Now),
	}
}
