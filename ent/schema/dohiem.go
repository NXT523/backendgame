package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type DoHiem struct{ ent.Schema }

func (DoHiem) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			StorageKey("ma_do_hiem").
			Immutable().
			Comment("Mã độ hiếm"),
		field.String("ten_do_hiem").
			MaxLen(20).NotEmpty().Unique().
			Comment("Tên độ hiếm"),
		field.Int("so_luong").
			Default(0).Comment("Số lượng"),
		field.String("mau_sac").
			MaxLen(7).Optional().Comment("Màu sắc HEX"),
		field.Float("sat_thuong_bonus").
			Default(0).Comment("Damage bonus"),
		field.Float("toc_do_danh_bonus").
			Default(0).Comment("Attack speed bonus"),
		field.Int("cap_bac").Comment("Cấp bậc").Unique(),
	}
}

func (DoHiem) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("ten_do_hiem").
			Unique().
			StorageKey("uq_dohiem_ten"), // 1

		index.Fields("so_luong").
			StorageKey("ix_dohiem_so_luong"), // 2

		index.Fields("mau_sac").
			StorageKey("ix_dohiem_mau_sac"), // 3
		index.Fields("sat_thuong_bonus").StorageKey("ix_dohiem_sat_bonus"),
		index.Fields("toc_do_danh_bonus").StorageKey("ix_dohiem_spd_bonus"),
	}
}

func (DoHiem) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("vu_khi", VuKhi.Type),
	}
}
