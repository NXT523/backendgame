package mappers

import (
	"game/internal/db"
	"game/internal/models"
)

// Mapping trả về mảng slice theo models struct VuKhiResponse
//
// Dùng cho search của Gorm
func ListVuKhiToResponsesSearchGORM(list []db.GVuKhi) []*models.VuKhiResponse {
	res := make([]*models.VuKhiResponse, 0, len(list))

	for _, vk := range list {

		item := &models.VuKhiResponse{
			MaVuKhi:        vk.MaVuKhi,
			TenVuKhi:       vk.TenVuKhi,
			SatThuongCoBan: vk.SatThuongCoBan,
			TocDoDanh:      vk.TocDoDanh,
			TamDanh:        vk.TamDanh,
			Version:        vk.Version,

			MaLoaiVuKhi: vk.MaLoaiVuKhi,
			MaDoHiem:    vk.MaDoHiem,
			MaHe:        vk.MaHe,
		}

		if vk.MoTa != nil {
			item.MoTa = *vk.MoTa
		}

		// Join data
		item.TenLoaiVuKhi = vk.LoaiVuKhi.TenLoaiVuKhi
		item.TenHe = vk.He.TenHe
		item.TenDoHiem = vk.DoHiem.TenDoHiem

		if vk.DoHiem.MauSac != nil {
			item.MauSac = *vk.DoHiem.MauSac
		}

		item.SoLuong = vk.DoHiem.SoLuong
		item.SatThuongBonus = vk.DoHiem.SatThuongBonus
		item.TocDoDanhBonus = vk.DoHiem.TocDoDanhBonus
		item.CapBac = vk.DoHiem.CapBac

		res = append(res, item)
	}

	return res
}

// GormVuKhiToModel chuyển đổi một thực thể GORM (db.GVuKhi) tới mô hình miền (models.VuKhis).
//
// Hàm này dùng để ánh xạ dữ liệu lấy từ cơ sở dữ liệu
// vào mô hình lớp nghiệp vụ để kho lưu trữ không
// hiển thị các thực thể cơ sở dữ liệu bên ngoài lớp của nó.
//
// Nó cũng hủy đăng ký một cách an toàn các trường có thể rỗng như MoTa.
func GormVuKhiToModel(vk *db.GVuKhi) *models.VuKhis {
	if vk == nil {
		return nil
	}

	m := &models.VuKhis{
		MaVuKhi:        vk.MaVuKhi,
		TenVuKhi:       vk.TenVuKhi,
		SatThuongCoBan: vk.SatThuongCoBan,
		TocDoDanh:      vk.TocDoDanh,
		TamDanh:        vk.TamDanh,
		MaLoaiVuKhi:    vk.MaLoaiVuKhi,
		MaDoHiem:       vk.MaDoHiem,
		MaHe:           vk.MaHe,
		Version:        vk.Version,
	}

	if vk.MoTa != nil {
		m.MoTa = *vk.MoTa
	}

	return m
}

// ModelToGormVuKhi chuyển đổi mô hình miền (models.VuKhis) thành một thực thể GORM (db.GVuKhi).
//
// Hàm này được sử dụng trước khi lưu dữ liệu vào cơ sở dữ liệu.
// Nó đảm bảo kho lưu trữ chỉ hoạt động với các thực thể cơ sở dữ liệu,
// trong khi lớp dịch vụ hoạt động với các mô hình miền.
//
// Các trường chuỗi như MoTa được chuyển đổi thành kiểu con trỏ
// vì các cột cơ sở dữ liệu có thể rỗng được biểu diễn dưới dạng con trỏ trong GORM.
func ModelToGormVuKhi(m *models.VuKhis) *db.GVuKhi {
	if m == nil {
		return nil
	}

	return &db.GVuKhi{
		MaVuKhi:        m.MaVuKhi,
		TenVuKhi:       m.TenVuKhi,
		SatThuongCoBan: m.SatThuongCoBan,
		TocDoDanh:      m.TocDoDanh,
		TamDanh:        m.TamDanh,
		MoTa:           &m.MoTa,
		MaLoaiVuKhi:    m.MaLoaiVuKhi,
		MaDoHiem:       m.MaDoHiem,
		MaHe:           m.MaHe,
		Version:        m.Version,
	}
}

// Slice mapping
func ListGormVuKhiToModels(list []db.GVuKhi) []*models.VuKhis {
	res := make([]*models.VuKhis, 0, len(list))
	for i := range list {
		res = append(res, GormVuKhiToModel(&list[i]))
	}
	return res
}
