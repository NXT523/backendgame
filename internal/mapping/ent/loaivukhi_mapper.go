package mappers

import (
	"game/ent"
	v1 "game/v1/proto"
)

func LoaiVuKhiEntToPB(m *ent.LoaiVuKhi) *v1.LoaiVuKhi {
	if m == nil {
		return nil
	}
	return &v1.LoaiVuKhi{
		MaLoaiVuKhi:  int32(m.ID),
		TenLoaiVuKhi: m.TenLoaiVuKhi,
		MoTa:         m.MoTa,
	}
}
