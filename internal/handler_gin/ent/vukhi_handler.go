package handler_gin

import (
	"game/internal/models"
	serviceEnt "game/internal/service/gin_service/ent"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// 1. ENT
type VuKhiHandlerEnt struct {
	service serviceEnt.VuKhiServiceEnt
}

func NewVuKhiHandlerEnt(service serviceEnt.VuKhiServiceEnt) *VuKhiHandlerEnt {
	return &VuKhiHandlerEnt{
		service: service,
	}
}

// Handler: Tạo vũ khí
func (h *VuKhiHandlerEnt) CreateEnt(c *gin.Context) {
	var req models.VuKhis

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.service.CreateEnt(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *VuKhiHandlerEnt) UpdateByNameEnt(c *gin.Context) {
	name := strings.TrimSpace(c.Param("name"))
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name required"})
		return
	}

	var req models.UpdateVuKhiReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Version <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "version required"})
		return
	}

	in := models.VuKhiUpdateInput{
		TenVuKhi:       req.TenVuKhi,
		SatThuongCoBan: req.SatThuongCoBan,
		TocDoDanh:      req.TocDoDanh,
		TamDanh:        req.TamDanh,
		MoTa:           req.MoTa,
		MaLoaiVuKhi:    req.MaLoaiVuKhi,
		MaDoHiem:       req.MaDoHiem,
		MaHe:           req.MaHe,
		Version:        req.Version,
	}

	res, err := h.service.UpdateByNameEnt(
		c.Request.Context(),
		name,
		in,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// Handler: Xóa vũ khí theo tên vũ khí truyền vào
func (h *VuKhiHandlerEnt) DeleteByNameEnt(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tên không hợp lệ"})
		return
	}
	if err := h.service.DeleteByNameEnt(c, name); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy hoặc lỗi khi xoá", "detail": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Đá xóa thành công vũ khí ": name})
}

// Handler: Lấy tất cả các vũ khí
//
// Dùng từ khóa kw để tìm kiếm vũ khí cụ thể
func (h *VuKhiHandlerEnt) GetAllEnt(c *gin.Context) {
	kw := c.Query("kw")
	items, err := h.service.GetAllEnt(c, kw)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": len(items), "items": items})
}

func (h *VuKhiHandlerEnt) GetAllTenVuKhiEnt(c *gin.Context) {
	items, err := h.service.GetAllTenVuKhiEnt(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": len(items), "items": items})
}

func (h *VuKhiHandlerEnt) GetAllSatThuongCoBanEnt(c *gin.Context) {
	items, err := h.service.GetAllSatThuongCoBanEnt(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": len(items), "items": items})
}

func (h *VuKhiHandlerEnt) GetAllTocDoDanhEnt(c *gin.Context) {
	items, err := h.service.GetAllTocDoDanhEnt(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": len(items), "items": items})
}

func (h *VuKhiHandlerEnt) GetAllTamDanhEnt(c *gin.Context) {
	items, err := h.service.GetAllTamDanhEnt(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": len(items), "items": items})
}

func (h *VuKhiHandlerEnt) SearchEnt(c *gin.Context) {
	var req models.VuKhiSearchRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Body JSON không hợp lệ",
			"detail": err.Error(),
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	arrangeAsc := c.Query("arrange_asc")
	arrangeDesc := c.Query("arrange_desc")

	items, total, err := h.service.SearchEnt(
		c.Request.Context(),
		req,
		page,
		pageSize,
		arrangeAsc,
		arrangeDesc,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if total == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Không tìm thấy vũ khí",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"items":     items,
	})
}
