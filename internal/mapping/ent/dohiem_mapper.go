package mappers

import (
	"game/ent"
	v1 "game/v1/proto"
)

// DoHiemEntToPB map ent.DoHiem → proto DoHiem
func DoHiemEntToPB(e *ent.DoHiem) *v1.DoHiem {
	if e == nil {
		return nil
	}
	out := &v1.DoHiem{
		MaDoHiem:       int32(e.ID),
		TenDoHiem:      e.TenDoHiem,
		SoLuong:        int32(e.SoLuong),
		MauSac:         e.MauSac,
		SatThuongBonus: e.SatThuongBonus,
		TocDoDanhBonus: e.TocDoDanhBonus,
		CapBac:         int32(e.CapBac),
	}
	return out
}
