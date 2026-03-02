package mappers

import (
	"game/ent"
	"game/internal/models"
	v1 "game/v1/proto"
)

// VuKhiEntToPB map ent.VuKhi → proto VuKhi
func VuKhiEntToPB(e *ent.VuKhi) *v1.VuKhi {
	if e == nil {
		return nil
	}

	out := &v1.VuKhi{
		MaVuKhi:        int32(e.ID),
		TenVuKhi:       e.TenVuKhi,
		SatThuongCoBan: int32(e.SatThuongCoBan),
		TocDoDanh:      e.TocDoDanh,
		TamDanh:        int32(e.TamDanh),
		MoTa:           e.MoTa,
		MaLoaiVuKhi:    int32(e.MaLoaiVuKhi),
		MaDoHiem:       int32(e.MaDoHiem),
		MaHe:           int32(e.MaHe),
		Version:        int32(e.Version),
	}
	if e.Edges.Loai != nil {
		out.TenLoai = e.Edges.Loai.TenLoaiVuKhi
	}
	if e.Edges.He != nil {
		out.TenHe = e.Edges.He.TenHe
	}
	if e.Edges.DoHiem != nil {
		out.TenDoHiem = e.Edges.DoHiem.TenDoHiem
		out.MauSac = e.Edges.DoHiem.MauSac
		out.SoLuong = int32(e.Edges.DoHiem.SoLuong)
		out.SatThuongBonus = e.Edges.DoHiem.SatThuongBonus
		out.TocDoDanhBonus = e.Edges.DoHiem.TocDoDanhBonus
		out.CapBac = int32(e.Edges.DoHiem.CapBac)
	}
	return out
}

// Trả về 1 kết quả được mapping theo models/vukhi.go
func VuKhiToModel(e *ent.VuKhi) *models.VuKhis {
	if e == nil {
		return nil
	}

	m := &models.VuKhis{
		MaVuKhi:        e.ID,
		TenVuKhi:       e.TenVuKhi,
		SatThuongCoBan: e.SatThuongCoBan,
		TocDoDanh:      e.TocDoDanh,
		TamDanh:        e.TamDanh,
		MoTa:           e.MoTa,
		MaLoaiVuKhi:    e.MaLoaiVuKhi,
		MaDoHiem:       e.MaDoHiem,
		MaHe:           e.MaHe,
		Version:        int(e.Version),
	}

	if e.Edges.Loai != nil {
		m.TenLoaiVuKhi = e.Edges.Loai.TenLoaiVuKhi
	}

	if e.Edges.He != nil {
		m.TenHe = e.Edges.He.TenHe
	}

	if e.Edges.DoHiem != nil {
		m.TenDoHiem = e.Edges.DoHiem.TenDoHiem
		m.MauSac = e.Edges.DoHiem.MauSac
		m.SoLuong = e.Edges.DoHiem.SoLuong
		m.SatThuongBonus = e.Edges.DoHiem.SatThuongBonus
		m.TocDoDanhBonus = e.Edges.DoHiem.TocDoDanhBonus
		m.CapBac = e.Edges.DoHiem.CapBac
	}

	return m
}

// Trả về mảng slice kết quả được mapping theo models/vukhi.go
func ListVuKhiToModels(list []*ent.VuKhi) []*models.VuKhis {
	result := make([]*models.VuKhis, 0, len(list))

	for _, e := range list {
		result = append(result, VuKhiToModel(e))
	}

	return result
}

// Mapping trả về mảng slice models struct VuKhiResponse
//
// Dùng cho search của Ent
func ListVuKhiToResponsesSearchENT(list []*models.VuKhis) []*models.VuKhiResponse {
	res := make([]*models.VuKhiResponse, 0, len(list))

	for _, m := range list {
		if m == nil {
			continue
		}

		res = append(res, &models.VuKhiResponse{
			MaVuKhi:        m.MaVuKhi,
			TenVuKhi:       m.TenVuKhi,
			SatThuongCoBan: m.SatThuongCoBan,
			TocDoDanh:      m.TocDoDanh,
			TamDanh:        m.TamDanh,
			MoTa:           m.MoTa,
			Version:        int(m.Version),

			MaLoaiVuKhi: m.MaLoaiVuKhi,
			MaDoHiem:    m.MaDoHiem,
			MaHe:        m.MaHe,

			TenLoaiVuKhi:   m.TenLoaiVuKhi,
			TenHe:          m.TenHe,
			TenDoHiem:      m.TenDoHiem,
			MauSac:         m.MauSac,
			SoLuong:        m.SoLuong,
			SatThuongBonus: m.SatThuongBonus,
			TocDoDanhBonus: m.TocDoDanhBonus,
			CapBac:         m.CapBac,
		})
	}

	return res
}
