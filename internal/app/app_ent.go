package app

import (
	"game/internal/middleware"
	routes_Ent "game/internal/routes/ent"

	"github.com/gin-gonic/gin"
)

type EntModule interface {
	RoutesEnt() []routes_Ent.Route
}
type ApplicationEnt struct {
	router *gin.Engine
	port   string
}

func NewApplicationEnt(port string, modules []EntModule) *ApplicationEnt {

	r := gin.Default()

	go middleware.CleanupClients()

	r.Use(
		middleware.HTTPLoggerMiddlewareGin(),
		middleware.RateLimiterGinMiddleware(),
	)

	routes_Ent.RegisterRoutes(r, getEntRoutes(modules)...)

	return &ApplicationEnt{
		router: r,
		port:   port,
	}
}

func getEntRoutes(modules []EntModule) []routes_Ent.Route {
	var result []routes_Ent.Route
	for _, m := range modules {
		result = append(result, m.RoutesEnt()...)
	}
	return result
}
func (a *ApplicationEnt) GetHandler() *gin.Engine {
	return a.router
}
