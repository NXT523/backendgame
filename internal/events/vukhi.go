package events

// SuKienVuKhiCreated là payload gửi lên Kafka khi tạo vũ khí mới
type SuKienVuKhiCreated struct {
	Type           string  `json:"type"`       // "VU_KHI_CREATED"
	ID             int     `json:"id"`         // ID vũ khí vừa tạo
	TenVuKhi       string  `json:"ten_vu_khi"` // tên vũ khí
	SatThuongCoBan int     `json:"sat_thuong_co_ban"`
	TocDoDanh      float64 `json:"toc_do_danh"`
	TamDanh        int     `json:"tam_danh"`
	MoTa           string  `json:"mo_ta"`
	MaLoai         int     `json:"ma_loai"`
	MaHe           int     `json:"ma_he"`
	MaDoHiem       int     `json:"ma_do_hiem"`
	CreatedAt      int64   `json:"created_at"` // time.Now().Unix()
	CreatedBy      string  `json:"created_by"` // sau này nếu có auth thì nhét user vào
}

type SuKienVuKhiUpdated struct {
	Type           string  `json:"type"`
	MaVuKhi        int     `json:"ma_vu_khi"`
	TenVuKhi       string  `json:"ten_vu_khi"`
	SatThuongCoBan int     `json:"sat_thuong_co_ban"`
	TocDoDanh      float64 `json:"toc_do_danh"`
	TamDanh        int     `json:"tam_danh"`
	MoTa           string  `json:"mo_ta"`
	MaLoai         int     `json:"ma_loai"`
	MaDoHiem       int     `json:"ma_do_hiem"`
	MaHe           int     `json:"ma_he"`
	Ts             int64   `json:"ts"`
	UpdatedAt      int64   `json:"updated_at"`
}

type VuKhiUpdateEvent struct {
    RequestID string `json:"request_id"`
    Name      string `json:"name"`
    Payload   []byte `json:"payload"`
    CreatedAt int64  `json:"created_at"`
}