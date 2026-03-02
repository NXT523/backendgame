package routes

import (
	handler_Ent "game/internal/handler_gin/ent"

	"github.com/gin-gonic/gin"
)

type VuKhiRoutes struct {
	handlerEnt  *handler_Ent.VuKhiHandlerEnt
}

func NewVuKhiRoutesEnt(handler *handler_Ent.VuKhiHandlerEnt) *VuKhiRoutes {
	return &VuKhiRoutes{
		handlerEnt: handler,
	}
}

func (vk *VuKhiRoutes) RegisterEnt(r *gin.RouterGroup) {
	vukhis := r.Group("/vukhi/api/ent")
	{
		vukhis.POST("/post", vk.handlerEnt.CreateEnt)
		vukhis.GET("/getall", vk.handlerEnt.GetAllEnt)
		vukhis.PUT("/put-by-name/:name", vk.handlerEnt.UpdateByNameEnt)
		vukhis.DELETE("/delete-by-name/:name", vk.handlerEnt.DeleteByNameEnt)

		vukhis.GET("/ten", vk.handlerEnt.GetAllTenVuKhiEnt)
		vukhis.GET("/satthuongcoban", vk.handlerEnt.GetAllSatThuongCoBanEnt)
		vukhis.GET("/tocdo", vk.handlerEnt.GetAllTocDoDanhEnt)
		vukhis.GET("/tamdanh", vk.handlerEnt.GetAllTamDanhEnt)
		vukhis.POST("/search", vk.handlerEnt.SearchEnt)
	}
}

