package routes

// import (
// 	"game/internal/handler"
// 	"game/internal/service"

// 	"github.com/gin-gonic/gin"
// )

// // 1. ENT
// func MountDoHiemEnt(r gin.IRouter, svc service.DoHiemServiceEnt) {
// 	h := handler.CreateDoHiemHandlerEnt(svc)
// 	g := r.Group("/dohiem/ent")
// 	{
// 		g.POST("/post", h.CreateEnt)

// 		g.GET("/get-by-name/:name", h.GetByNameEnt)
// 		g.PUT("/put-by-name/:name", h.UpdateByNameEnt)
// 		g.DELETE("/delete-by-name/:name", h.DeleteByNameEnt)

// 		g.GET("/getall", h.GetAllEnt)
// 		g.GET("/ten", h.GetAllTenEnt)
// 		g.GET("/soluong", h.GetAllSoLuongEnt)
// 		g.GET("/mausac", h.GetAllMauSacEnt)
// 		g.GET("/satthuong", h.GetAllSatThuongBonusEnt)
// 		g.GET("/tocdo", h.GetAllTocDoDanhBonusEnt)
// 	}
// }

// // 2. GORM
// func MountDoHiemGorm(r gin.IRouter, svc service.DoHiemServiceGorm) {
// 	h := handler.CreateDoHiemHandlerGorm(svc)
// 	g := r.Group("/dohiem/gorm")
// 	{
// 		g.POST("/post", h.CreateGorm)
// 		g.GET("/getall", h.GetAllGorm)
// 		g.GET("/get/:id", h.GetIdGorm)
// 		g.PUT("/put/:id", h.UpdateGorm)
// 		g.DELETE("/delete/:id", h.DeleteGorm)

// 		g.GET("/get-by-name/:name", h.GetByNameGorm)
// 		g.PUT("/put-by-name/:name", h.UpdateByNameGorm)
// 		g.DELETE("/delete-by-name/:name", h.DeleteByNameGorm)

// 		g.GET("/ten", h.GetAllTenGorm)
// 		g.GET("/soluong", h.GetAllSoLuongGorm)
// 		g.GET("/mausac", h.GetAllMauSacGorm)
// 		g.GET("/satthuong", h.GetAllSatThuongBonusGorm)
// 		g.GET("/tocdo", h.GetAllTocDoDanhBonusGorm)
// 	}
// }

// // 3. RAW
// func MountDoHiemRaw(r gin.IRouter, svc service.DoHiemServiceRaw) {
// 	h := handler.CreateDoHiemHandlerRaw(svc)
// 	g := r.Group("/dohiem/raw")
// 	{
// 		g.POST("/post", h.CreateRaw)
// 		g.GET("/getall", h.GetAllRaw)
// 		g.GET("/get/:id", h.GetIdRaw)
// 		g.PUT("/put/:id", h.UpdateRaw)
// 		g.DELETE("/delete/:id", h.DeleteRaw)

// 		g.GET("/get-by-name/:name", h.GetByNameRaw)
// 		g.PUT("/put-by-name/:name", h.UpdateByNameRaw)
// 		g.DELETE("/delete-by-name/:name", h.DeleteByNameRaw)

// 		g.GET("/ten", h.GetAllTenRaw)
// 		g.GET("/soluong", h.GetAllSoLuongRaw)
// 		g.GET("/mausac", h.GetAllMauSacRaw)
// 		g.GET("/satthuong", h.GetAllSatThuongBonusRaw)
// 		g.GET("/tocdo", h.GetAllTocDoDanhBonusRaw)
// 	}
// }
