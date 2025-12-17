package router

import (
	"game/internal/handler"
	"game/internal/service"

	"github.com/gin-gonic/gin"
)

// 1. ENT
func MountVuKhiEnt(r gin.IRouter, svc service.VuKhiServiceEnt) {
	h := handler.CreateVuKhiHandlerEnt(svc)
	g := r.Group("/vukhi/ent")
	{
		g.POST("/post", h.CreateEnt)                         // thêm vu khi
		g.PUT("/put-by-name/:name", h.UpdateByNameEnt)       // (vukhi 1) update theo ten vũ khí
		g.DELETE("/delete-by-name/:name", h.DeleteByNameEnt) // (vukhi 1) xóa theo tên vũ khí

		g.GET("/getall", h.GetAllEnt)                       // get tất cả các vũ khí
		g.GET("/ten", h.GetAllTenVuKhiEnt)                  // chỉ lấy tên vũ khí
		g.GET("/satthuongcoban", h.GetAllSatThuongCoBanEnt) // chỉ lấy sát thương cơ bản
		g.GET("/tocdo", h.GetAllTocDoDanhEnt)               // chỉ lấy tốc độ đánh
		g.GET("/tamdanh", h.GetAllTamDanhEnt)               // chỉ lấy tầm đánh

		g.POST("/search", h.SearchEnt)
	}
}

// 2. GORM
func MountVuKhiGorm(r gin.IRouter, svc service.VuKhiServiceGorm) {
	h := handler.CreateVuKhiHandlerGorm(svc)
	g := r.Group("/vukhi/gorm")
	{

		g.POST("/post", h.CreateGorm)                         // thêm vũ khí
		g.PUT("/put-by-name/:name", h.UpdateByNameGorm)       // update theo tên
		g.DELETE("/delete-by-name/:name", h.DeleteByNameGorm) // xóa theo tên

		g.GET("/getall", h.GetAllGorm)                       // lấy tất cả
		g.GET("/ten", h.GetAllTenVuKhiGorm)                  // chỉ tên
		g.GET("/satthuongcoban", h.GetAllSatThuongCoBanGorm) // chỉ sát thương cơ bản
		g.GET("/tocdo", h.GetAllTocDoDanhGorm)               // chỉ tốc độ đánh
		g.GET("/tamdanh", h.GetAllTamDanhGorm)               // chỉ tầm đánh

		g.POST("/search", h.SearchGorm)
	}
}

// 3. RAW
func MountVuKhiRaw(r gin.IRouter, svc service.VuKhiServiceRaw) {
	h := handler.CreateVuKhiHandlerRaw(svc)
	g := r.Group("/vukhi/raw")
	{
		g.POST("/post", h.CreateRaw)                         // thêm vũ khí
		g.PUT("/put-by-name/:name", h.UpdateByNameRaw)       // update theo tên
		g.DELETE("/delete-by-name/:name", h.DeleteByNameRaw) // xóa theo tên

		g.GET("/getall", h.GetAllRaw)                       // lấy tất cả
		g.GET("/ten", h.GetAllTenVuKhiRaw)                  // chỉ tên
		g.GET("/satthuongcoban", h.GetAllSatThuongCoBanRaw) // chỉ sát thương cơ bản
		g.GET("/tocdo", h.GetAllTocDoDanhRaw)               // chỉ tốc độ đánh
		g.GET("/tamdanh", h.GetAllTamDanhRaw)               // chỉ tầm đánh

		g.POST("/search", h.SearchRaw)
	}
}
