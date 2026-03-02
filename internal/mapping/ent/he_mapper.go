package mappers

import (
	"game/ent"
	v1 "game/v1/proto"
)

// HeEntToPB map ent.He → proto He
func HeEntToPB(e *ent.He) *v1.He {
	if e == nil {
		return nil
	}
	out := &v1.He{
		MaHe:  int32(e.ID),
		TenHe: e.TenHe,
		MoTa:  e.MoTa,
	}
	return out
}
