package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type LoaiVuKhi struct{ ent.Schema }

func (LoaiVuKhi) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			StorageKey("ma_loai").
			Immutable().
			Comment("Mã loại vũ khí"),
		field.String("ten_loai").
			MaxLen(50).NotEmpty().Unique().
			Comment("Tên loại"),
		field.String("mo_ta").
			MaxLen(255).Optional().Comment("Mô tả"),
	}
}

func (LoaiVuKhi) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("ten_loai").
			Unique().
			StorageKey("uq_loaivukhi_ten"),// 1
	}
}

func (LoaiVuKhi) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("vu_khi", VuKhi.Type),
	}
}
