package router

import (
	"game/internal/handler"
	"game/internal/service"

	"github.com/gin-gonic/gin"
)

// 1. ENT
func MountLoaiVuKhiEnt(r gin.IRouter, svc service.LoaiVuKhiServiceEnt) {
	h := handler.CreateLoaiVuKhiHandlerEnt(svc)
	lvk := r.Group("/loaivukhi/ent")
	{
		lvk.POST("/post", h.CreateEnt)
		lvk.GET("/getall", h.GetAllEnt)
		lvk.GET("/get/:id", h.GetIdEnt)
		lvk.PUT("/put/:id", h.UpdateEnt)
		lvk.DELETE("/delete/:id", h.DeleteEnt)

		lvk.GET("/nameloaivukhi", h.GetAllLoaiVuKhiEnt) // 1: chỉ lấy tên loại vũ khí
	}
}

// 2. GORM
func MountLoaiVuKhiGorm(r gin.IRouter, svc service.LoaiVuKhiServiceGorm) {
	h := handler.CreateLoaiVuKhiHandlerGorm(svc)
	g := r.Group("/loaivukhi/gorm")
	{
		g.POST("/post", h.CreateGorm)
		g.GET("/getall", h.GetAllGorm)
		g.GET("/get/:id", h.GetIdGorm)
		g.PUT("/put/:id", h.UpdateGorm)
		g.DELETE("/delete/:id", h.DeleteGorm)

		g.GET("/nameloaivukhi", h.GetAllLoaiVuKhiGorm)
	}
}

// 3. RAW
func MountLoaiVuKhiRaw(r gin.IRouter, svc service.LoaiVuKhiServiceRaw) {
	h := handler.CreateLoaiVuKhiHandlerRaw(svc)
	g := r.Group("/loaivukhi/raw")
	{
		g.POST("/post", h.CreateRaw)
		g.GET("/getall", h.GetAllRaw)
		g.GET("/get/:id", h.GetIdRaw)
		g.PUT("/put/:id", h.UpdateRaw)
		g.DELETE("/delete/:id", h.DeleteRaw)

		g.GET("/nameloaivukhi", h.GetAllLoaiVuKhiRaw)
	}
}
