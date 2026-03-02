package routes

// import (
// 	"game/internal/handler"
// 	"game/internal/service"

// 	"github.com/gin-gonic/gin"
// )

// // 1. ENT
// func MountHeEnt(r gin.IRouter, svc service.HeServiceEnt) {
// 	h := handler.CreateHeHandlerEnt(svc)
// 	heGroup := r.Group("/he/ent")
// 	{
// 		heGroup.POST("/post", h.CreateEnt)
// 		heGroup.GET("/getall", h.GetAllEnt)
// 		heGroup.GET("/get/:id", h.GetIdEnt)
// 		heGroup.PUT("/put/:id", h.UpdateEnt)
// 		heGroup.DELETE("/delete/:id", h.DeleteEnt)

// 		heGroup.GET("/get-by-name/:name", h.GetByNameEnt)
// 		heGroup.PUT("/put-by-name/:name", h.UpdateByNameEnt)
// 		heGroup.DELETE("/delete-by-name/:name", h.DeleteByNameEnt)

// 		heGroup.GET("/names", h.GetNamesEnt)
// 	}
// }

// // 2. GORM
// func MountHeGorm(r gin.IRouter, svc service.HeServiceGorm) {
// 	h := handler.CreateHeHandlerGorm(svc)
// 	g := r.Group("/he/gorm")
// 	{
// 		g.POST("/post", h.CreateGorm)
// 		g.GET("/getall", h.GetAllGorm)
// 		g.GET("/get/:id", h.GetIdGorm)
// 		g.PUT("/put/:id", h.UpdateGorm)
// 		g.DELETE("/delete/:id", h.DeleteGorm)

// 		g.GET("/get-by-name/:name", h.GetByNameGorm)
// 		g.PUT("/put-by-name/:name", h.UpdateByNameGorm)
// 		g.DELETE("/delete-by-name/:name", h.DeleteByNameGorm)

// 		g.GET("/names", h.GetNamesGorm)
// 	}
// }

// // 3. RAW
// func MountHeRaw(r gin.IRouter, svc service.HeServiceRaw) {
// 	h := handler.CreateHeHandlerRaw(svc)
// 	g := r.Group("/he/raw")
// 	{
// 		g.POST("/post", h.CreateRaw)
// 		g.GET("/getall", h.GetAllRaw)
// 		g.GET("/get/:id", h.GetIdRaw)
// 		g.PUT("/put/:id", h.UpdateRaw)
// 		g.DELETE("/delete/:id", h.DeleteRaw)

// 		g.GET("/get-by-name/:name", h.GetByNameRaw)
// 		g.PUT("/put-by-name/:name", h.UpdateByNameRaw)
// 		g.DELETE("/delete-by-name/:name", h.DeleteByNameRaw)

// 		g.GET("/names", h.GetNamesRaw)
// 	}
// }
