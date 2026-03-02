package handler_gin

// import (
// 	"database/sql"
// 	"errors"
// 	"net/http"
// 	"net/url"
// 	"strconv"
// 	"strings"

// 	"game/internal/service"

// 	"github.com/gin-gonic/gin"
// 	"gorm.io/gorm"
// )

// // 1. ENT
// type HeHandlerEnt struct {
// 	sv service.HeServiceEnt
// }

// func CreateHeHandlerEnt(s service.HeServiceEnt) *HeHandlerEnt {
// 	return &HeHandlerEnt{sv: s}
// }

// type heReqEnt struct {
// 	TenHe string  `json:"ten_he" binding:"required"`
// 	MoTa  *string `json:"mo_ta"`
// }

// func (h *HeHandlerEnt) CreateEnt(c *gin.Context) {
// 	var r heReqEnt
// 	if err := c.ShouldBindJSON(&r); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON không hợp lệ", "detail": err.Error()})
// 		return
// 	}
// 	in := service.HeInputEnt{
// 		TenHe: r.TenHe,
// 		MoTa:  r.MoTa,
// 	}
// 	v, err := h.sv.CreateEnt(c, in)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusCreated, v)
// }

// func (h *HeHandlerEnt) GetAllEnt(c *gin.Context) {
// 	q := c.Query("q")
// 	items, err := h.sv.GetAllEnt(c, q)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"total": len(items), "items": items})
// }

// func (h *HeHandlerEnt) GetIdEnt(c *gin.Context) {
// 	id, err := strconv.Atoi(c.Param("id"))
// 	if err != nil || id <= 0 {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
// 		return
// 	}
// 	v, err := h.sv.GetIdEnt(c, id)
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy"})
// 		return
// 	}
// 	c.JSON(http.StatusOK, v)
// }

// func (h *HeHandlerEnt) UpdateEnt(c *gin.Context) {
// 	id, err := strconv.Atoi(c.Param("id"))
// 	if err != nil || id <= 0 {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
// 		return
// 	}
// 	var r heReqEnt
// 	if err := c.ShouldBindJSON(&r); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON không hợp lệ", "detail": err.Error()})
// 		return
// 	}
// 	in := service.HeInputEnt{
// 		TenHe: r.TenHe,
// 		MoTa:  r.MoTa,
// 	}
// 	v, err := h.sv.UpdateEnt(c, id, in)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, v)
// }

// func (h *HeHandlerEnt) DeleteEnt(c *gin.Context) {
// 	id, err := strconv.Atoi(c.Param("id"))
// 	if err != nil || id <= 0 {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
// 		return
// 	}
// 	if err := h.sv.DeleteEnt(c, id); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"deleted_id": id})
// }

// // GET /he/ent/get-by-name/:name
// func (h *HeHandlerEnt) GetByNameEnt(c *gin.Context) {
// 	raw := c.Param("name")
// 	if strings.TrimSpace(raw) == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Tên hệ không hợp lệ"})
// 		return
// 	}
// 	name, err := url.PathUnescape(raw)
// 	if err != nil || strings.TrimSpace(name) == "" {
// 		name = raw
// 	}

// 	v, err := h.sv.GetByNameEnt(c, name)
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy"})
// 		return
// 	}
// 	c.JSON(http.StatusOK, v)
// }

// // PUT /he/ent/put-by-name/:name
// func (h *HeHandlerEnt) UpdateByNameEnt(c *gin.Context) {
// 	raw := c.Param("name")
// 	if strings.TrimSpace(raw) == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Tên hệ không hợp lệ"})
// 		return
// 	}
// 	name, err := url.PathUnescape(raw)
// 	if err != nil || strings.TrimSpace(name) == "" {
// 		name = raw
// 	}

// 	var r heReqEnt
// 	if err := c.ShouldBindJSON(&r); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON không hợp lệ", "detail": err.Error()})
// 		return
// 	}
// 	in := service.HeInputEnt{TenHe: r.TenHe, MoTa: r.MoTa}

// 	v, err := h.sv.UpdateByNameEnt(c, name, in)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, v)
// }

// // DELETE /he/ent/delete-by-name/:name
// func (h *HeHandlerEnt) DeleteByNameEnt(c *gin.Context) {
// 	raw := c.Param("name")
// 	if strings.TrimSpace(raw) == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Tên hệ không hợp lệ"})
// 		return
// 	}
// 	name, err := url.PathUnescape(raw)
// 	if err != nil || strings.TrimSpace(name) == "" {
// 		name = raw
// 	}

// 	if err := h.sv.DeleteByNameEnt(c, name); err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy hoặc lỗi khi xoá", "detail": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"deleted_name": name})
// }

// func (h *HeHandlerEnt) GetNamesEnt(c *gin.Context) {
// 	q := c.Query("q")
// 	names, err := h.sv.GetAllNamesEnt(c, q)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{
// 		"total": len(names),
// 		"items": names,
// 	})
// }

// // 2. GORM
// type HeHandlerGorm struct {
// 	sv service.HeServiceGorm
// }

// func CreateHeHandlerGorm(s service.HeServiceGorm) *HeHandlerGorm {
// 	return &HeHandlerGorm{sv: s}
// }

// type heReqGorm struct {
// 	TenHe string  `json:"ten_he" binding:"required"`
// 	MoTa  *string `json:"mo_ta"`
// }

// func (h *HeHandlerGorm) CreateGorm(c *gin.Context) {
// 	var r heReqGorm
// 	if err := c.ShouldBindJSON(&r); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON không hợp lệ", "detail": err.Error()})
// 		return
// 	}
// 	in := service.HeInputGorm{
// 		TenHe: r.TenHe,
// 		MoTa:  r.MoTa,
// 	}
// 	v, err := h.sv.CreateGorm(c, in)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusCreated, v)
// }

// func (h *HeHandlerGorm) GetAllGorm(c *gin.Context) {
// 	q := c.Query("q")
// 	items, err := h.sv.GetAllGorm(c, q)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"total": len(items), "items": items})
// }

// func (h *HeHandlerGorm) GetIdGorm(c *gin.Context) {
// 	id, err := strconv.Atoi(c.Param("id"))
// 	if err != nil || id <= 0 {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
// 		return
// 	}
// 	v, err := h.sv.GetIdGorm(c, id)
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy"})
// 		return
// 	}
// 	c.JSON(http.StatusOK, v)
// }

// func (h *HeHandlerGorm) UpdateGorm(c *gin.Context) {
// 	id, err := strconv.Atoi(c.Param("id"))
// 	if err != nil || id <= 0 {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
// 		return
// 	}
// 	var r heReqGorm
// 	if err := c.ShouldBindJSON(&r); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON không hợp lệ", "detail": err.Error()})
// 		return
// 	}
// 	in := service.HeInputGorm{
// 		TenHe: r.TenHe,
// 		MoTa:  r.MoTa,
// 	}
// 	v, err := h.sv.UpdateGorm(c, id, in)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, v)
// }

// func (h *HeHandlerGorm) DeleteGorm(c *gin.Context) {
// 	id, err := strconv.Atoi(c.Param("id"))
// 	if err != nil || id <= 0 {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
// 		return
// 	}
// 	if err := h.sv.DeleteGorm(c, id); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"deleted_id": id})
// }

// // GET /he/gorm/get-by-name/:name
// func (h *HeHandlerGorm) GetByNameGorm(c *gin.Context) {
// 	raw := c.Param("name")
// 	if strings.TrimSpace(raw) == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Tên hệ không hợp lệ"})
// 		return
// 	}
// 	name, err := url.PathUnescape(raw)
// 	if err != nil || strings.TrimSpace(name) == "" {
// 		name = raw
// 	}

// 	v, err := h.sv.GetByNameGorm(c, name)
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy"})
// 			return
// 		}
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, v)
// }

// // PUT /he/gorm/put-by-name/:name
// func (h *HeHandlerGorm) UpdateByNameGorm(c *gin.Context) {
// 	raw := c.Param("name")
// 	if strings.TrimSpace(raw) == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Tên hệ không hợp lệ"})
// 		return
// 	}
// 	name, err := url.PathUnescape(raw)
// 	if err != nil || strings.TrimSpace(name) == "" {
// 		name = raw
// 	}

// 	var r heReqGorm
// 	if err := c.ShouldBindJSON(&r); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON không hợp lệ", "detail": err.Error()})
// 		return
// 	}
// 	in := service.HeInputGorm{TenHe: r.TenHe, MoTa: r.MoTa}

// 	v, err := h.sv.UpdateByNameGorm(c, name, in)
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy"})
// 			return
// 		}
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, v)
// }

// // DELETE /he/gorm/delete-by-name/:name
// func (h *HeHandlerGorm) DeleteByNameGorm(c *gin.Context) {
// 	raw := c.Param("name")
// 	if strings.TrimSpace(raw) == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Tên hệ không hợp lệ"})
// 		return
// 	}
// 	name, err := url.PathUnescape(raw)
// 	if err != nil || strings.TrimSpace(name) == "" {
// 		name = raw
// 	}

// 	if err := h.sv.DeleteByNameGorm(c, name); err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy"})
// 			return
// 		}
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"deleted_name": name})
// }

// func (h *HeHandlerGorm) GetNamesGorm(c *gin.Context) {
// 	q := c.Query("q")
// 	names, err := h.sv.GetAllNamesGorm(c, q)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"total": len(names), "items": names})
// }

// // 3. RAW
// type HeHandlerRaw struct{ sv service.HeServiceRaw }

// func CreateHeHandlerRaw(s service.HeServiceRaw) *HeHandlerRaw {
// 	return &HeHandlerRaw{sv: s}
// }

// type heReqRaw struct {
// 	TenHe string  `json:"ten_he" binding:"required"`
// 	MoTa  *string `json:"mo_ta"`
// }

// func (h *HeHandlerRaw) CreateRaw(c *gin.Context) {
// 	var r heReqRaw
// 	if err := c.ShouldBindJSON(&r); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON không hợp lệ", "detail": err.Error()})
// 		return
// 	}
// 	in := service.HeInputRaw{
// 		TenHe: r.TenHe,
// 		MoTa:  r.MoTa,
// 	}
// 	v, err := h.sv.CreateRaw(c, in)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusCreated, v)
// }

// func (h *HeHandlerRaw) GetAllRaw(c *gin.Context) {
// 	q := c.Query("q")
// 	items, err := h.sv.GetAllRaw(c, q)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"total": len(items), "items": items})
// }

// func (h *HeHandlerRaw) GetIdRaw(c *gin.Context) {
// 	id, err := strconv.Atoi(c.Param("id"))
// 	if err != nil || id <= 0 {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
// 		return
// 	}
// 	v, err := h.sv.GetIdRaw(c, id)
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy"})
// 		return
// 	}
// 	c.JSON(http.StatusOK, v)
// }

// func (h *HeHandlerRaw) UpdateRaw(c *gin.Context) {
// 	id, err := strconv.Atoi(c.Param("id"))
// 	if err != nil || id <= 0 {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
// 		return
// 	}
// 	var r heReqRaw
// 	if err := c.ShouldBindJSON(&r); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON không hợp lệ", "detail": err.Error()})
// 		return
// 	}
// 	in := service.HeInputRaw{
// 		TenHe: r.TenHe,
// 		MoTa:  r.MoTa,
// 	}
// 	v, err := h.sv.UpdateRaw(c, id, in)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, v)
// }

// func (h *HeHandlerRaw) DeleteRaw(c *gin.Context) {
// 	id, err := strconv.Atoi(c.Param("id"))
// 	if err != nil || id <= 0 {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
// 		return
// 	}
// 	if err := h.sv.DeleteRaw(c, id); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"deleted_id": id})
// }

// // GET /he/raw/get-by-name/:name
// func (h *HeHandlerRaw) GetByNameRaw(c *gin.Context) {
// 	raw := c.Param("name")
// 	if strings.TrimSpace(raw) == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Tên hệ không hợp lệ"})
// 		return
// 	}
// 	name, err := url.PathUnescape(raw)
// 	if err != nil || strings.TrimSpace(name) == "" {
// 		name = raw
// 	}

// 	v, err := h.sv.GetByNameRaw(c, name)
// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy"})
// 			return
// 		}
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, v)
// }

// // PUT /he/raw/put-by-name/:name
// func (h *HeHandlerRaw) UpdateByNameRaw(c *gin.Context) {
// 	raw := c.Param("name")
// 	if strings.TrimSpace(raw) == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Tên hệ không hợp lệ"})
// 		return
// 	}
// 	name, err := url.PathUnescape(raw)
// 	if err != nil || strings.TrimSpace(name) == "" {
// 		name = raw
// 	}

// 	var r heReqRaw
// 	if err := c.ShouldBindJSON(&r); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON không hợp lệ", "detail": err.Error()})
// 		return
// 	}
// 	in := service.HeInputRaw{TenHe: r.TenHe, MoTa: r.MoTa}

// 	v, err := h.sv.UpdateByNameRaw(c, name, in)
// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy"})
// 			return
// 		}
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, v)
// }

// // DELETE /he/raw/delete-by-name/:name
// func (h *HeHandlerRaw) DeleteByNameRaw(c *gin.Context) {
// 	raw := c.Param("name")
// 	if strings.TrimSpace(raw) == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Tên hệ không hợp lệ"})
// 		return
// 	}
// 	name, err := url.PathUnescape(raw)
// 	if err != nil || strings.TrimSpace(name) == "" {
// 		name = raw
// 	}

// 	if err := h.sv.DeleteByNameRaw(c, name); err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy"})
// 			return
// 		}
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"deleted_name": name})
// }

// func (h *HeHandlerRaw) GetNamesRaw(c *gin.Context) {
// 	q := c.Query("q")
// 	names, err := h.sv.GetAllNamesRaw(c, q)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"total": len(names), "items": names})
// }
