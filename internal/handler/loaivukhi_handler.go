package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"game/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 1. ENT
type LoaiVuKhiHandlerEnt struct {
	sv service.LoaiVuKhiServiceEnt
}

func CreateLoaiVuKhiHandlerEnt(s service.LoaiVuKhiServiceEnt) *LoaiVuKhiHandlerEnt {
	return &LoaiVuKhiHandlerEnt{sv: s}
}

type loaiVKReqEnt struct {
	TenLoai string  `json:"ten_loai" binding:"required"`
	MoTa    *string `json:"mo_ta"`
}

func (h *LoaiVuKhiHandlerEnt) CreateEnt(c *gin.Context) {
	var r loaiVKReqEnt
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON không hợp lệ", "detail": err.Error()})
		return
	}
	in := service.LoaiVuKhiInputEnt{
		TenLoai: r.TenLoai,
		MoTa:    r.MoTa,
	}
	v, err := h.sv.CreateEnt(c, in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, v)
}

func (h *LoaiVuKhiHandlerEnt) GetAllEnt(c *gin.Context) {
	q := c.Query("q")
	items, err := h.sv.GetAllEnt(c, q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": len(items), "items": items})
}

func (h *LoaiVuKhiHandlerEnt) GetIdEnt(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}
	v, err := h.sv.GetIdEnt(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy"})
		return
	}
	c.JSON(http.StatusOK, v)
}

func (h *LoaiVuKhiHandlerEnt) UpdateEnt(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}
	var r loaiVKReqEnt
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON không hợp lệ", "detail": err.Error()})
		return
	}
	in := service.LoaiVuKhiInputEnt{
		TenLoai: r.TenLoai,
		MoTa:    r.MoTa,
	}
	v, err := h.sv.CreateEnt(c, in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, v)
}

func (h *LoaiVuKhiHandlerEnt) DeleteEnt(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}
	if err := h.sv.DeleteEnt(c, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"deleted_id": id})
}

func (h *LoaiVuKhiHandlerEnt) GetAllLoaiVuKhiEnt(c *gin.Context) {
	q := c.Query("q")
	names, err := h.sv.GetTenLoaiEnt(c, q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"total": len(names),
		"items": names,
	})
}

// 2. GORM
type LoaiVuKhiHandlerGorm struct {
	sv service.LoaiVuKhiServiceGorm
}

func CreateLoaiVuKhiHandlerGorm(s service.LoaiVuKhiServiceGorm) *LoaiVuKhiHandlerGorm {
	return &LoaiVuKhiHandlerGorm{sv: s}
}

type loaiVKReqGorm struct {
	TenLoai string  `json:"ten_loai" binding:"required"`
	MoTa    *string `json:"mo_ta"`
}

func (h *LoaiVuKhiHandlerGorm) CreateGorm(c *gin.Context) {
	var r loaiVKReqGorm
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON không hợp lệ", "detail": err.Error()})
		return
	}
	in := service.LoaiVuKhiInputGorm{
		TenLoai: r.TenLoai,
		MoTa:    r.MoTa,
	}
	v, err := h.sv.CreateGorm(c, in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, v)
}

func (h *LoaiVuKhiHandlerGorm) GetAllGorm(c *gin.Context) {
	q := c.Query("q")
	items, err := h.sv.GetAllGorm(c, q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": len(items), "items": items})
}

func (h *LoaiVuKhiHandlerGorm) GetIdGorm(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}
	v, err := h.sv.GetIdGorm(c, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, v)
}

func (h *LoaiVuKhiHandlerGorm) UpdateGorm(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}
	var r loaiVKReqGorm
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON không hợp lệ", "detail": err.Error()})
		return
	}
	in := service.LoaiVuKhiInputGorm{
		TenLoai: r.TenLoai,
		MoTa:    r.MoTa,
	}
	v, err := h.sv.UpdateGorm(c, id, in)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, v)
}

func (h *LoaiVuKhiHandlerGorm) DeleteGorm(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}
	if err := h.sv.DeleteGorm(c, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"deleted_id": id})
}

func (h *LoaiVuKhiHandlerGorm) GetAllLoaiVuKhiGorm(c *gin.Context) {
	q := c.Query("q")
	names, err := h.sv.GetTenLoaiGorm(c, q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": len(names), "items": names})
}

// 3. RAW SQL
type LoaiVuKhiHandlerRaw struct {
	sv service.LoaiVuKhiServiceRaw
}

func CreateLoaiVuKhiHandlerRaw(s service.LoaiVuKhiServiceRaw) *LoaiVuKhiHandlerRaw {
	return &LoaiVuKhiHandlerRaw{sv: s}
}

type loaiVKReqRaw struct {
	TenLoai string  `json:"ten_loai" binding:"required"`
	MoTa    *string `json:"mo_ta"`
}

func (h *LoaiVuKhiHandlerRaw) CreateRaw(c *gin.Context) {
	var r loaiVKReqRaw
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON không hợp lệ", "detail": err.Error()})
		return
	}
	in := service.LoaiVuKhiInputRaw{
		TenLoai: r.TenLoai,
		MoTa:    r.MoTa,
	}
	v, err := h.sv.CreateRaw(c, in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, v)
}

func (h *LoaiVuKhiHandlerRaw) GetAllRaw(c *gin.Context) {
	q := c.Query("q")
	items, err := h.sv.GetAllRaw(c, q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": len(items), "items": items})
}

func (h *LoaiVuKhiHandlerRaw) GetIdRaw(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}
	v, err := h.sv.GetIdRaw(c, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, v)
}

func (h *LoaiVuKhiHandlerRaw) UpdateRaw(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}
	var r loaiVKReqRaw
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON không hợp lệ", "detail": err.Error()})
		return
	}
	in := service.LoaiVuKhiInputRaw{
		TenLoai: r.TenLoai,
		MoTa:    r.MoTa,
	}
	v, err := h.sv.UpdateRaw(c, id, in)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, v)
}

func (h *LoaiVuKhiHandlerRaw) DeleteRaw(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}
	if err := h.sv.DeleteRaw(c, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"deleted_id": id})
}

func (h *LoaiVuKhiHandlerRaw) GetAllLoaiVuKhiRaw(c *gin.Context) {
	q := c.Query("q")
	names, err := h.sv.GetTenLoaiRaw(c, q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": len(names), "items": names})
}
