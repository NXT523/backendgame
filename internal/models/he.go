package models

type Hes struct {
	MaVuKhi        int     `json:"ma_vu_khi"`
	TenVuKhi       string  `json:"ten_vu_khi"`
	SatThuongCoBan int     `json:"sat_thuong_co_ban"`
	TocDoDanh      float64 `json:"toc_do_danh"`
	TamDanh        int     `json:"tam_danh"`
	MoTa           string  `json:"mo_ta"`
	Version        int     `json:"version"`

	MaLoaiVuKhi int `json:"ma_loai_vu_khi"`
	MaDoHiem    int `json:"ma_do_hiem"`
	MaHe        int `json:"ma_he"`

	TenLoaiVuKhi   string  `json:"ten_loai_vu_khi,omitempty"`
	TenDoHiem      string  `json:"ten_do_hiem,omitempty"`
	TenHe          string  `json:"ten_he,omitempty"`
	SatThuongBonus float64 `json:"sat_thuong_bonus,omitempty"`
	TocDoDanhBonus float64 `json:"toc_do_danh_bonus,omitempty"`
	SoLuong        int     `json:"so_luong,omitempty"`
	MauSac         string  `json:"mau_sac,omitempty"`
	CapBac         int     `json:"cap_bac,omitempty"`
}
