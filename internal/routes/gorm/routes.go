package routes

import (
	"game/internal/middleware"

	"github.com/gin-gonic/gin"
)

type Route interface {
	RegisterGorm(r *gin.RouterGroup)
}

func RegisterRoutesGorm(r *gin.Engine, routes ...Route) {
	r.Use(middleware.AuthMiddleware())
	api := r.Group("")

	for _, route := range routes {
		route.RegisterGorm(api)
	}
}
