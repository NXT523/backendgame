package handler

import (
	"game/internal/repo"
	"game/internal/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// 1. ENT
type VuKhiHandlerEnt struct{ sv service.VuKhiServiceEnt }

func CreateVuKhiHandlerEnt(s service.VuKhiServiceEnt) *VuKhiHandlerEnt {
	return &VuKhiHandlerEnt{sv: s}
}

type vuKhiReqEnt struct {
	TenVuKhi       string  `json:"ten_vu_khi" binding:"required"`
	SatThuongCoBan int     `json:"sat_thuong_co_ban"`
	TocDoDanh      float64 `json:"toc_do_danh"`
	TamDanh        int     `json:"tam_danh"`
	MoTa           *string `json:"mo_ta"`
	MaLoai         int     `json:"ma_loai" binding:"required"`
	MaDoHiem       int     `json:"ma_do_hiem" binding:"required"`
	MaHe           int     `json:"ma_he" binding:"required"`
}

func (h *VuKhiHandlerEnt) CreateEnt(c *gin.Context) {
	var r vuKhiReqEnt
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON không hợp lệ", "detail": err.Error()})
		return
	}
	in := service.VuKhiInputEnt{
		TenVuKhi:       r.TenVuKhi,
		SatThuongCoBan: r.SatThuongCoBan,
		TocDoDanh:      r.TocDoDanh,
		TamDanh:        r.TamDanh,
		MoTa:           r.MoTa,
		MaLoai:         r.MaLoai,
		MaDoHiem:       r.MaDoHiem,
		MaHe:           r.MaHe,
	}
	v, err := h.sv.CreateEnt(c, in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, v)
}

func (h *VuKhiHandlerEnt) GetAllEnt(c *gin.Context) {
	q := c.Query("q")
	items, err := h.sv.GetAllEnt(c, q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": len(items), "items": items})
}

func (h *VuKhiHandlerEnt) UpdateByNameEnt(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tên không hợp lệ"})
		return
	}
	var r vuKhiReqEnt
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON không hợp lệ", "detail": err.Error()})
		return
	}
	in := service.VuKhiInputEnt{
		TenVuKhi:       r.TenVuKhi,
		SatThuongCoBan: r.SatThuongCoBan,
		TocDoDanh:      r.TocDoDanh,
		TamDanh:        r.TamDanh,
		MoTa:           r.MoTa,
		MaLoai:         r.MaLoai,
		MaDoHiem:       r.MaDoHiem,
		MaHe:           r.MaHe,
	}
	v, err := h.sv.UpdateByNameEnt(c, name, in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, v)
}

func (h *VuKhiHandlerEnt) DeleteByNameEnt(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tên không hợp lệ"})
		return
	}
	if err := h.sv.DeleteByNameEnt(c, name); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy hoặc lỗi khi xoá", "detail": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"deleted_name": name})
}

type VuKhiView struct {
	TenVuKhi       string  `json:"ten_vu_khi"`
	TenLoai        string  `json:"ten_loai"`
	TenHe          string  `json:"ten_he"`
	TenDoHiem      string  `json:"ten_do_hiem"`
	SatThuongCoBan int     `json:"sat_thuong_co_ban"`
	TocDoDanh      float64 `json:"toc_do_danh"`
	TamDanh        int     `json:"tam_danh"`
	SatThuongBonus float64 `json:"sat_thuong_bonus"`
	TocDoDanhBonus float64 `json:"toc_do_danh_bonus"`
	TongSatThuong  float64 `json:"tong_sat_thuong"`
	TongTocDoDanh  float64 `json:"tong_toc_do_danh"`
}

func (h *VuKhiHandlerEnt) GetAllTenVuKhiEnt(c *gin.Context) {
	items, err := h.sv.GetAllTenVuKhiEnt(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": len(items), "items": items})
}

func (h *VuKhiHandlerEnt) GetAllSatThuongCoBanEnt(c *gin.Context) {
	items, err := h.sv.GetAllSatThuongCoBanEnt(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": len(items), "items": items})
}

func (h *VuKhiHandlerEnt) GetAllTocDoDanhEnt(c *gin.Context) {
	items, err := h.sv.GetAllTocDoDanhEnt(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": len(items), "items": items})
}

func (h *VuKhiHandlerEnt) GetAllTamDanhEnt(c *gin.Context) {
	items, err := h.sv.GetAllTamDanhEnt(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": len(items), "items": items})
}

func (h *VuKhiHandlerEnt) SearchEnt(c *gin.Context) {
	var req repo.VuKhiSearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Body JSON không hợp lệ", "detail": err.Error()})
		return
	}

	items, err := h.sv.SearchEnt(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(items) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Không tìm thấy vũ khí"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": len(items), "items": items})
}

// 2. GORM
type VuKhiHandlerGorm struct{ sv service.VuKhiServiceGorm }

func CreateVuKhiHandlerGorm(s service.VuKhiServiceGorm) *VuKhiHandlerGorm {
	return &VuKhiHandlerGorm{sv: s}
}

type vuKhiReqGorm struct {
	TenVuKhi       string  `json:"ten_vu_khi" binding:"required"`
	SatThuongCoBan int     `json:"sat_thuong_co_ban"`
	TocDoDanh      float64 `json:"toc_do_danh"`
	TamDanh        int     `json:"tam_danh"`
	MoTa           *string `json:"mo_ta"`
	MaLoai         int     `json:"ma_loai" binding:"required"`
	MaDoHiem       int     `json:"ma_do_hiem" binding:"required"`
	MaHe           int     `json:"ma_he" binding:"required"`
}

// POST /vukhi/gorm/post
func (h *VuKhiHandlerGorm) CreateGorm(c *gin.Context) {
	var r vuKhiReqGorm
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON không hợp lệ", "detail": err.Error()})
		return
	}
	in := service.VuKhiInputGorm{
		TenVuKhi:       r.TenVuKhi,
		SatThuongCoBan: r.SatThuongCoBan,
		TocDoDanh:      r.TocDoDanh,
		TamDanh:        r.TamDanh,
		MoTa:           r.MoTa,
		MaLoai:         r.MaLoai,
		MaDoHiem:       r.MaDoHiem,
		MaHe:           r.MaHe,
	}
	v, err := h.sv.CreateGorm(c, in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, v)
}

// GET /vukhi/gorm/getall
func (h *VuKhiHandlerGorm) GetAllGorm(c *gin.Context) {
	q := c.Query("q")
	items, err := h.sv.GetAllGorm(c, q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": len(items), "items": items})
}

// PUT /vukhi/gorm/put-by-name/:name
func (h *VuKhiHandlerGorm) UpdateByNameGorm(c *gin.Context) {
	name := strings.TrimSpace(c.Param("name"))
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tên không hợp lệ"})
		return
	}
	var r vuKhiReqGorm
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON không hợp lệ", "detail": err.Error()})
		return
	}
	in := service.VuKhiInputGorm{
		TenVuKhi:       r.TenVuKhi,
		SatThuongCoBan: r.SatThuongCoBan,
		TocDoDanh:      r.TocDoDanh,
		TamDanh:        r.TamDanh,
		MoTa:           r.MoTa,
		MaLoai:         r.MaLoai,
		MaDoHiem:       r.MaDoHiem,
		MaHe:           r.MaHe,
	}
	v, err := h.sv.UpdateByNameGorm(c, name, in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, v)
}

// DELETE /vukhi/gorm/delete-by-name/:name
func (h *VuKhiHandlerGorm) DeleteByNameGorm(c *gin.Context) {
	name := strings.TrimSpace(c.Param("name"))
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tên không hợp lệ"})
		return
	}
	if err := h.sv.DeleteByNameGorm(c, name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"deleted_name": name})
}

// GET /vukhi/gorm/ten
func (h *VuKhiHandlerGorm) GetAllTenVuKhiGorm(c *gin.Context) {
	names, err := h.sv.GetAllTenVuKhiGorm(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": len(names), "items": names})
}

// GET /vukhi/gorm/satthuongcoban
func (h *VuKhiHandlerGorm) GetAllSatThuongCoBanGorm(c *gin.Context) {
	values, err := h.sv.GetAllSatThuongCoBanGorm(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": len(values), "items": values})
}

// GET /vukhi/gorm/tocdo
func (h *VuKhiHandlerGorm) GetAllTocDoDanhGorm(c *gin.Context) {
	values, err := h.sv.GetAllTocDoDanhGorm(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": len(values), "items": values})
}

// GET /vukhi/gorm/tamdanh
func (h *VuKhiHandlerGorm) GetAllTamDanhGorm(c *gin.Context) {
	values, err := h.sv.GetAllTamDanhGorm(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": len(values), "items": values})
}

// POST /vukhi/gorm/search
func (h *VuKhiHandlerGorm) SearchGorm(c *gin.Context) {
	var req repo.VuKhiSearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Body JSON không hợp lệ", "detail": err.Error()})
		return
	}
	items, err := h.sv.SearchGorm(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(items) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Không tìm thấy vũ khí"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": len(items), "items": items})
}

// 3. RAW
type VuKhiHandlerRaw struct{ sv service.VuKhiServiceRaw }

func CreateVuKhiHandlerRaw(s service.VuKhiServiceRaw) *VuKhiHandlerRaw {
	return &VuKhiHandlerRaw{sv: s}
}

type vuKhiReqRaw struct {
	TenVuKhi       string  `json:"ten_vu_khi" binding:"required"`
	SatThuongCoBan int     `json:"sat_thuong_co_ban"`
	TocDoDanh      float64 `json:"toc_do_danh"`
	TamDanh        int     `json:"tam_danh"`
	MoTa           *string `json:"mo_ta"`
	MaLoai         int     `json:"ma_loai" binding:"required"`
	MaDoHiem       int     `json:"ma_do_hiem" binding:"required"`
	MaHe           int     `json:"ma_he" binding:"required"`
}

// POST /vukhi/raw/post
func (h *VuKhiHandlerRaw) CreateRaw(c *gin.Context) {
	var r vuKhiReqRaw
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON không hợp lệ", "detail": err.Error()})
		return
	}
	in := service.VuKhiInputRaw{
		TenVuKhi:       r.TenVuKhi,
		SatThuongCoBan: r.SatThuongCoBan,
		TocDoDanh:      r.TocDoDanh,
		TamDanh:        r.TamDanh,
		MoTa:           r.MoTa,
		MaLoai:         r.MaLoai,
		MaDoHiem:       r.MaDoHiem,
		MaHe:           r.MaHe,
	}
	v, err := h.sv.CreateRaw(c, in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, v)
}

// GET /vukhi/raw/getall
func (h *VuKhiHandlerRaw) GetAllRaw(c *gin.Context) {
	q := c.Query("q")
	items, err := h.sv.GetAllRaw(c, q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": len(items), "items": items})
}

// PUT /vukhi/raw/put-by-name/:name
func (h *VuKhiHandlerRaw) UpdateByNameRaw(c *gin.Context) {
	name := strings.TrimSpace(c.Param("name"))
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tên không hợp lệ"})
		return
	}
	var r vuKhiReqRaw
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON không hợp lệ", "detail": err.Error()})
		return
	}
	in := service.VuKhiInputRaw{
		TenVuKhi:       r.TenVuKhi,
		SatThuongCoBan: r.SatThuongCoBan,
		TocDoDanh:      r.TocDoDanh,
		TamDanh:        r.TamDanh,
		MoTa:           r.MoTa,
		MaLoai:         r.MaLoai,
		MaDoHiem:       r.MaDoHiem,
		MaHe:           r.MaHe,
	}
	v, err := h.sv.UpdateByNameRaw(c, name, in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, v)
}

// DELETE /vukhi/raw/delete-by-name/:name
func (h *VuKhiHandlerRaw) DeleteByNameRaw(c *gin.Context) {
	name := strings.TrimSpace(c.Param("name"))
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tên không hợp lệ"})
		return
	}
	if err := h.sv.DeleteByNameRaw(c, name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"deleted_name": name})
}

// GET /vukhi/raw/ten
func (h *VuKhiHandlerRaw) GetAllTenVuKhiRaw(c *gin.Context) {
	names, err := h.sv.GetAllTenVuKhiRaw(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": len(names), "items": names})
}

// GET /vukhi/raw/satthuongcoban
func (h *VuKhiHandlerRaw) GetAllSatThuongCoBanRaw(c *gin.Context) {
	values, err := h.sv.GetAllSatThuongCoBanRaw(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": len(values), "items": values})
}

// GET /vukhi/raw/tocdo
func (h *VuKhiHandlerRaw) GetAllTocDoDanhRaw(c *gin.Context) {
	values, err := h.sv.GetAllTocDoDanhRaw(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": len(values), "items": values})
}

// GET /vukhi/raw/tamdanh
func (h *VuKhiHandlerRaw) GetAllTamDanhRaw(c *gin.Context) {
	values, err := h.sv.GetAllTamDanhRaw(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": len(values), "items": values})
}

// POST /vukhi/raw/search
func (h *VuKhiHandlerRaw) SearchRaw(c *gin.Context) {
	var req repo.VuKhiSearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Body JSON không hợp lệ", "detail": err.Error()})
		return
	}
	items, err := h.sv.SearchRaw(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(items) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Không tìm thấy vũ khí"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": len(items), "items": items})
}
