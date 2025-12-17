package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type He struct{ ent.Schema }

func (He) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			StorageKey("ma_he").
			Immutable().
			Comment("Mã hệ"),
		field.String("ten_he").
			MaxLen(50).NotEmpty().Unique().
			Comment("Tên hệ"),
		field.String("mo_ta").
			MaxLen(255).Optional().Comment("Mô tả"),
	}
}

func (He) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("ten_he").
			Unique().
			StorageKey("uq_he_ten"), // 1
	}
}

func (He) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("vu_khi", VuKhi.Type),
	}
}
