package handler

import (
	"database/sql"
	"errors"
	"game/internal/service"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 1. ENT
type DoHiemHandlerEnt struct {
	sv service.DoHiemServiceEnt
}

func CreateDoHiemHandlerEnt(s service.DoHiemServiceEnt) *DoHiemHandlerEnt {
	return &DoHiemHandlerEnt{sv: s}
}

type doHiemReqEnt struct {
	TenDoHiem      string   `json:"ten_do_hiem" binding:"required"`
	MauSac         string   `json:"mau_sac" binding:"required"`
	SoLuong        *int     `json:"so_luong"`
	SatThuongBonus *float64 `json:"sat_thuong_bonus"`
	TocDoDanhBonus *float64 `json:"toc_do_danh_bonus"`
}

func (h *DoHiemHandlerEnt) CreateEnt(c *gin.Context) {
	var r doHiemReqEnt
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON không hợp lệ", "detail": err.Error()})
		return
	}

	in := service.DoHiemInputEnt{
		TenDoHiem:      r.TenDoHiem,
		MauSac:         r.MauSac,
		SoLuong:        *r.SoLuong,
		SatThuongBonus: *r.SatThuongBonus,
		TocDoDanhBonus: *r.TocDoDanhBonus,
	}

	v, err := h.sv.CreateEnt(c, in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, v)
}

func (h *DoHiemHandlerEnt) GetAllEnt(c *gin.Context) {
	q := c.Query("q")
	items, err := h.sv.GetAllEnt(c, q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": len(items), "items": items})
}

// GET /dohiem/ent/get-by-name/:name
func (h *DoHiemHandlerEnt) GetByNameEnt(c *gin.Context) {
	raw := c.Param("name")
	if strings.TrimSpace(raw) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tên hệ không hợp lệ"})
		return
	}
	name, err := url.PathUnescape(raw)
	if err != nil || strings.TrimSpace(name) == "" {
		name = raw
	}

	v, err := h.sv.GetByNameEnt(c, name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy"})
		return
	}
	c.JSON(http.StatusOK, v)
}

// PUT /dohiem/ent/put-by-name/:name
func (h *DoHiemHandlerEnt) UpdateByNameEnt(c *gin.Context) {
	rawName := c.Param("name")
	if strings.TrimSpace(rawName) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tên không hợp lệ"})
		return
	}
	name, err := url.PathUnescape(rawName)
	if err != nil || strings.TrimSpace(name) == "" {
		name = rawName
	}

	var r doHiemReqEnt
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON không hợp lệ", "detail": err.Error()})
		return
	}

	in := service.DoHiemInputEnt{
		TenDoHiem:      r.TenDoHiem,
		MauSac:         r.MauSac,
		SoLuong:        *r.SoLuong,
		SatThuongBonus: *r.SatThuongBonus,
		TocDoDanhBonus: *r.TocDoDanhBonus,
	}

	v, err := h.sv.UpdateByNameEnt(c, name, in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, v)
}

// DELETE /dohiem/ent/delete-by-name/:name
func (h *DoHiemHandlerEnt) DeleteByNameEnt(c *gin.Context) {
	rawName := c.Param("name")
	if strings.TrimSpace(rawName) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tên không hợp lệ"})
		return
	}
	name, err := url.PathUnescape(rawName)
	if err != nil || strings.TrimSpace(name) == "" {
		name = rawName
	}

	if err := h.sv.DeleteByNameEnt(c, name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"deleted_name": name})
}

func (h *DoHiemHandlerEnt) GetAllTenEnt(c *gin.Context) {
	items, err := h.sv.GetAllTen(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": len(items), "items": items})
}

func (h *DoHiemHandlerEnt) GetAllSoLuongEnt(c *gin.Context) {
	items, err := h.sv.GetAllSoLuong(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": len(items), "items": items})
}

func (h *DoHiemHandlerEnt) GetAllMauSacEnt(c *gin.Context) {
	items, err := h.sv.GetAllMauSac(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": len(items), "items": items})
}

func (h *DoHiemHandlerEnt) GetAllSatThuongBonusEnt(c *gin.Context) {
	items, err := h.sv.GetAllSatThuongBonus(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": len(items), "items": items})
}

func (h *DoHiemHandlerEnt) GetAllTocDoDanhBonusEnt(c *gin.Context) {
	items, err := h.sv.GetAllTocDoDanhBonus(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": len(items), "items": items})
}

// 2. GORM
type DoHiemHandlerGorm struct {
	sv service.DoHiemServiceGorm
}

func CreateDoHiemHandlerGorm(s service.DoHiemServiceGorm) *DoHiemHandlerGorm {
	return &DoHiemHandlerGorm{sv: s}
}

type doHiemReqGorm struct {
	TenDoHiem      string   `json:"ten_do_hiem" binding:"required"`
	MauSac         string   `json:"mau_sac" binding:"required"`
	SoLuong        *int     `json:"so_luong" binding:"required"`
	SatThuongBonus *float64 `json:"sat_thuong_bonus" binding:"required"`
	TocDoDanhBonus *float64 `json:"toc_do_danh_bonus" binding:"required"`
}

func (h *DoHiemHandlerGorm) CreateGorm(c *gin.Context) {
	var r doHiemReqGorm
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON không hợp lệ", "detail": err.Error()})
		return
	}
	in := service.DoHiemInputGorm{
		TenDoHiem:      r.TenDoHiem,
		MauSac:         r.MauSac,
		SoLuong:        *r.SoLuong,
		SatThuongBonus: *r.SatThuongBonus,
		TocDoDanhBonus: *r.TocDoDanhBonus,
	}
	v, err := h.sv.CreateGorm(c, in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, v)
}

func (h *DoHiemHandlerGorm) GetAllGorm(c *gin.Context) {
	q := c.Query("q")
	items, err := h.sv.GetAllGorm(c, q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": len(items), "items": items})
}

func (h *DoHiemHandlerGorm) GetIdGorm(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}
	v, err := h.sv.GetIdGorm(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy"})
		return
	}
	c.JSON(http.StatusOK, v)
}

func (h *DoHiemHandlerGorm) UpdateGorm(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}
	var r doHiemReqGorm
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON không hợp lệ", "detail": err.Error()})
		return
	}
	in := service.DoHiemInputGorm{
		TenDoHiem:      r.TenDoHiem,
		MauSac:         r.MauSac,
		SoLuong:        *r.SoLuong,
		SatThuongBonus: *r.SatThuongBonus,
		TocDoDanhBonus: *r.TocDoDanhBonus,
	}
	v, err := h.sv.UpdateGorm(c, id, in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, v)
}

func (h *DoHiemHandlerGorm) DeleteGorm(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}
	if err := h.sv.DeleteGorm(c, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"deleted_id": id})
}

func (h *DoHiemHandlerGorm) GetByNameGorm(c *gin.Context) {
	raw := c.Param("name")
	if strings.TrimSpace(raw) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tên hệ không hợp lệ"})
		return
	}
	name, err := url.PathUnescape(raw)
	if err != nil || strings.TrimSpace(name) == "" {
		name = raw
	}
	res, err := h.sv.GetByNameGorm(c, name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *DoHiemHandlerGorm) UpdateByNameGorm(c *gin.Context) {
	raw := c.Param("name")
	if strings.TrimSpace(raw) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tên hệ không hợp lệ"})
		return
	}
	name, err := url.PathUnescape(raw)
	if err != nil || strings.TrimSpace(name) == "" {
		name = raw
	}

	var input doHiemReqGorm
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON không hợp lệ", "detail": err.Error()})
		return
	}

	// Convert sang service.DoHiemInputGorm
	serviceInput := service.DoHiemInputGorm{
		TenDoHiem:      input.TenDoHiem,
		MauSac:         input.MauSac,
		SoLuong:        *input.SoLuong,
		SatThuongBonus: *input.SatThuongBonus,
		TocDoDanhBonus: *input.TocDoDanhBonus,
	}

	res, err := h.sv.UpdateByNameGorm(c, name, serviceInput)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *DoHiemHandlerGorm) DeleteByNameGorm(c *gin.Context) {
	raw := c.Param("name")
	if strings.TrimSpace(raw) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tên hệ không hợp lệ"})
		return
	}
	name, err := url.PathUnescape(raw)
	if err != nil || strings.TrimSpace(name) == "" {
		name = raw
	}

	if err := h.sv.DeleteByNameGorm(c, name); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *DoHiemHandlerGorm) GetAllTenGorm(c *gin.Context) {
	items, err := h.sv.GetAllTenGorm(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": len(items), "items": items})
}

func (h *DoHiemHandlerGorm) GetAllSoLuongGorm(c *gin.Context) {
	items, err := h.sv.GetAllSoLuongGorm(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": len(items), "items": items})
}

func (h *DoHiemHandlerGorm) GetAllMauSacGorm(c *gin.Context) {
	items, err := h.sv.GetAllMauSacGorm(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": len(items), "items": items})
}

func (h *DoHiemHandlerGorm) GetAllSatThuongBonusGorm(c *gin.Context) {
	items, err := h.sv.GetAllSatThuongBonusGorm(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": len(items), "items": items})
}

func (h *DoHiemHandlerGorm) GetAllTocDoDanhBonusGorm(c *gin.Context) {
	items, err := h.sv.GetAllTocDoDanhBonusGorm(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": len(items), "items": items})
}

// 3. RAW
type DoHiemHandlerRaw struct {
	sv service.DoHiemServiceRaw
}

func CreateDoHiemHandlerRaw(s service.DoHiemServiceRaw) *DoHiemHandlerRaw {
	return &DoHiemHandlerRaw{sv: s}
}

type doHiemReqRaw struct {
	TenDoHiem      string   `json:"ten_do_hiem" binding:"required"`
	MauSac         string   `json:"mau_sac" binding:"required"`
	SoLuong        *int     `json:"so_luong" binding:"required"`
	SatThuongBonus *float64 `json:"sat_thuong_bonus" binding:"required"`
	TocDoDanhBonus *float64 `json:"toc_do_danh_bonus" binding:"required"`
}

func (h *DoHiemHandlerRaw) CreateRaw(c *gin.Context) {
	var r doHiemReqRaw
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON không hợp lệ", "detail": err.Error()})
		return
	}
	in := service.DoHiemInputRaw{
		TenDoHiem:      r.TenDoHiem,
		MauSac:         r.MauSac,
		SoLuong:        *r.SoLuong,
		SatThuongBonus: *r.SatThuongBonus,
		TocDoDanhBonus: *r.TocDoDanhBonus,
	}
	v, err := h.sv.CreateRaw(c, in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, v)
}

func (h *DoHiemHandlerRaw) GetAllRaw(c *gin.Context) {
	q := c.Query("q")
	items, err := h.sv.GetAllRaw(c, q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": len(items), "items": items})
}

func (h *DoHiemHandlerRaw) GetIdRaw(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}
	v, err := h.sv.GetIdRaw(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy"})
		return
	}
	c.JSON(http.StatusOK, v)
}

func (h *DoHiemHandlerRaw) UpdateRaw(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}
	var r doHiemReqRaw
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON không hợp lệ", "detail": err.Error()})
		return
	}
	in := service.DoHiemInputRaw{
		TenDoHiem:      r.TenDoHiem,
		MauSac:         r.MauSac,
		SoLuong:        *r.SoLuong,
		SatThuongBonus: *r.SatThuongBonus,
		TocDoDanhBonus: *r.TocDoDanhBonus,
	}
	v, err := h.sv.UpdateRaw(c, id, in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, v)
}

func (h *DoHiemHandlerRaw) DeleteRaw(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}
	if err := h.sv.DeleteRaw(c, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"deleted_id": id})
}

func (h *DoHiemHandlerRaw) GetByNameRaw(c *gin.Context) {
	name := c.Param("name")
	res, err := h.sv.GetByNameRaw(c.Request.Context(), name)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *DoHiemHandlerRaw) UpdateByNameRaw(c *gin.Context) {
	name := c.Param("name")

	var input doHiemReqRaw
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON không hợp lệ", "detail": err.Error()})
		return
	}

	serviceInput := service.DoHiemInputRaw{
		TenDoHiem:      input.TenDoHiem,
		MauSac:         input.MauSac,
		SoLuong:        *input.SoLuong,
		SatThuongBonus: *input.SatThuongBonus,
		TocDoDanhBonus: *input.TocDoDanhBonus,
	}

	res, err := h.sv.UpdateByNameRaw(c.Request.Context(), name, serviceInput)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *DoHiemHandlerRaw) DeleteByNameRaw(c *gin.Context) {
	name := c.Param("name")
	if err := h.sv.DeleteByNameRaw(c.Request.Context(), name); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *DoHiemHandlerRaw) GetAllTenRaw(c *gin.Context) {
	items, err := h.sv.GetAllTenRaw(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": len(items), "items": items})
}

func (h *DoHiemHandlerRaw) GetAllSoLuongRaw(c *gin.Context) {
	items, err := h.sv.GetAllSoLuongRaw(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": len(items), "items": items})
}

func (h *DoHiemHandlerRaw) GetAllMauSacRaw(c *gin.Context) {
	items, err := h.sv.GetAllMauSacRaw(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": len(items), "items": items})
}

func (h *DoHiemHandlerRaw) GetAllSatThuongBonusRaw(c *gin.Context) {
	items, err := h.sv.GetAllSatThuongBonusRaw(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": len(items), "items": items})
}

func (h *DoHiemHandlerRaw) GetAllTocDoDanhBonusRaw(c *gin.Context) {
	items, err := h.sv.GetAllTocDoDanhBonusRaw(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": len(items), "items": items})
}
