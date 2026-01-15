// vu_khis.go
package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type VuKhi struct{ ent.Schema }

func (VuKhi) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			StorageKey("ma_vu_khi").
			Immutable().
			Comment("Mã vũ khí"),
		field.String("ten_vu_khi").
			MaxLen(100).
			NotEmpty(),
		field.Int("sat_thuong_co_ban").
			Default(0),
		field.Float("toc_do_danh").
			Default(1.0),
		field.Int("tam_danh").
			Default(1),
		field.String("mo_ta").
			MaxLen(255).
			Optional(),
		field.Int64("version").
			Default(1),

		// FK: NOT NULL vì edge .Required()
		field.Int("ma_loai").
			Comment("FK -> loai_vu_khis.ma_loai"),
		field.Int("ma_do_hiem").
			Comment("FK -> do_hiems.ma_do_hiem"),
		field.Int("ma_he").
			Comment("FK -> hes.ma_he"),
	}
}

func (VuKhi) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("ten_vu_khi").
			Unique().
			StorageKey("uq_vukhi_ten"), // 1

		index.Fields("ma_loai", "ma_do_hiem").
			StorageKey("ix_vukhi_loai_dohiem"), // 2

		index.Fields("ma_loai", "ma_he").
			StorageKey("ix_vukhi_loai_he"), // 3

		index.Fields("ma_loai", "ma_he", "ma_do_hiem").
			StorageKey("ix_vukhi_loai_he_dohiem"), // 4

		index.Fields("ma_loai", "sat_thuong_co_ban").
			StorageKey("ix_vukhi_loai_satthuong"), // 5

		index.Fields("ma_do_hiem").
			StorageKey("ix_vukhi_ma_dohiem"), // 6
	}
}

func (VuKhi) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("loai", LoaiVuKhi.Type).
			Ref("vu_khi").
			Field("ma_loai").
			Unique().
			Required(),

		edge.From("do_hiem", DoHiem.Type).
			Ref("vu_khi").
			Field("ma_do_hiem").
			Unique().
			Required(),

		edge.From("he", He.Type).
			Ref("vu_khi").
			Field("ma_he").
			Unique().
			Required(),
	}
}
