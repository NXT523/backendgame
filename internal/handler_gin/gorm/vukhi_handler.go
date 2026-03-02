package handler_gin

import (
	"game/internal/models"
	service_Gorm "game/internal/service/gin_service/gorm"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type VuKhiHandlerGorm struct {
	service service_Gorm.VuKhiServiceGorm
}

func NewVuKhiHandlerGorm(service service_Gorm.VuKhiServiceGorm) *VuKhiHandlerGorm {
	return &VuKhiHandlerGorm{
		service: service,
	}
}

// Handler: Tạo vũ khí
func (h *VuKhiHandlerGorm) CreateGorm(c *gin.Context) {
	var req models.VuKhis

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.service.CreateGorm(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *VuKhiHandlerGorm) UpdateByNameGorm(c *gin.Context) {
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

	res, err := h.service.UpdateByNameGorm(
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
func (h *VuKhiHandlerGorm) DeleteByNameGorm(c *gin.Context) {
	name := strings.TrimSpace(c.Param("name"))
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tên không hợp lệ"})
		return
	}

	if err := h.service.DeleteByNameGorm(c.Request.Context(), name); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":  "Không tìm thấy hoặc lỗi khi xoá",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Đã xoá thành công vũ khí", "name": name})
}

// Handler: Lấy tất cả các vũ khí
//
// Dùng từ khóa kw để tìm kiếm vũ khí cụ thể
func (h *VuKhiHandlerGorm) GetAllGorm(c *gin.Context) {
	kw := c.Query("kw")

	items, err := h.service.GetAllGorm(c.Request.Context(), kw)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total": len(items),
		"items": items,
	})
}

func (h *VuKhiHandlerGorm) GetAllTenVuKhiGorm(c *gin.Context) {
	items, err := h.service.GetAllTenVuKhiGorm(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": len(items), "items": items})
}

func (h *VuKhiHandlerGorm) GetAllSatThuongCoBanGorm(c *gin.Context) {
	items, err := h.service.GetAllSatThuongCoBanGorm(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": len(items), "items": items})
}

func (h *VuKhiHandlerGorm) GetAllTocDoDanhGorm(c *gin.Context) {
	items, err := h.service.GetAllTocDoDanhGorm(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": len(items), "items": items})
}

func (h *VuKhiHandlerGorm) GetAllTamDanhGorm(c *gin.Context) {
	items, err := h.service.GetAllTamDanhGorm(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": len(items), "items": items})
}

func (h *VuKhiHandlerGorm) SearchGorm(c *gin.Context) {
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

	items, total, err := h.service.SearchGorm(
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
