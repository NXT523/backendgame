package app

import (
	"game/internal/middleware"
	routes_Gorm "game/internal/routes/gorm"
	"log"

	"github.com/gin-gonic/gin"
)

type GormModule interface {
	RoutesGorm() []routes_Gorm.Route
}
type ApplicationGorm struct {
	router *gin.Engine
	port   string
}

func NewApplicationGorm(port string, modules []GormModule) *ApplicationGorm {

	r := gin.Default()

	go middleware.CleanupClients()

	r.Use(
		middleware.HTTPLoggerMiddlewareGin(),
		middleware.RateLimiterGinMiddleware(),
	)

	routes_Gorm.RegisterRoutesGorm(r, getGormRoutes(modules)...)

	return &ApplicationGorm{
		router: r,
		port:   port,
	}
}

func (a *ApplicationGorm) Run() error {
	log.Printf("🚀 GORM (Gin REST) chạy tại http://localhost%s\n", a.port)
	return a.router.Run(a.port)
}

func getGormRoutes(modules []GormModule) []routes_Gorm.Route {
	var result []routes_Gorm.Route
	for _, m := range modules {
		result = append(result, m.RoutesGorm()...)
	}
	return result
}
