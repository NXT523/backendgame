package routes

import (
	handler_Gorm "game/internal/handler_gin/gorm"

	"github.com/gin-gonic/gin"
)

type VuKhiRoutes struct {
	handlerGorm  *handler_Gorm.VuKhiHandlerGorm
}

func NewVuKhiRoutesGorm(handler *handler_Gorm.VuKhiHandlerGorm) *VuKhiRoutes {
	return &VuKhiRoutes{
		handlerGorm: handler,
	}
}

func (vk *VuKhiRoutes) RegisterGorm(r *gin.RouterGroup) {
	vukhis := r.Group("/vukhi/api/gorm")
	{

		vukhis.POST("/post", vk.handlerGorm.CreateGorm)
		vukhis.GET("/getall", vk.handlerGorm.GetAllGorm)
		vukhis.PUT("/put-by-name/:name", vk.handlerGorm.UpdateByNameGorm)
		vukhis.DELETE("/delete-by-name/:name", vk.handlerGorm.UpdateByNameGorm)

		vukhis.GET("/ten", vk.handlerGorm.GetAllTenVuKhiGorm)
		vukhis.GET("/satthuongcoban", vk.handlerGorm.GetAllSatThuongCoBanGorm)
		vukhis.GET("/tocdo", vk.handlerGorm.GetAllTocDoDanhGorm)
		vukhis.GET("/tamdanh", vk.handlerGorm.GetAllTamDanhGorm)
		vukhis.POST("/search", vk.handlerGorm.SearchGorm)
	}
}

