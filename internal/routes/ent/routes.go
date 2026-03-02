package routes

import (
	"game/internal/middleware"

	"github.com/gin-gonic/gin"
)

type Route interface {
	RegisterEnt(r *gin.RouterGroup)
}

func RegisterRoutes(r *gin.Engine, routes ...Route) {
	r.Use(middleware.AuthMiddleware())
	api := r.Group("")

	for _, route := range routes {
		route.RegisterEnt(api)
	}
}
