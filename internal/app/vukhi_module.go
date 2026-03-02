package app

import (
	"game/internal/db"
	handler_Ent "game/internal/handler_gin/ent"
	handler_Gorm "game/internal/handler_gin/gorm"
	repository_Ent "game/internal/repository/gin_repository/ent"
	repository_Gorm "game/internal/repository/gin_repository/gorm"
	routes_Ent "game/internal/routes/ent"
	routes_Gorm "game/internal/routes/gorm"
	service_Ent "game/internal/service/gin_service/ent"
	service_Gorm "game/internal/service/gin_service/gorm"
	"log"
)

type VuKhiModule struct {
	routesEnt  routes_Ent.Route
	routesGorm routes_Gorm.Route
}

func NewVuKhiModule() *VuKhiModule {
	entClient, err := db.OpenEnt()
	if err != nil {
		log.Fatalf("Lỗi kết nối Ent: %v", err)
	}

	gormClient, err := db.OpenGorm()
	if err != nil {
		log.Fatalf("Lỗi kết nối Gorm: %v", err)
	}
	// khởi tạo repository
	vukhiRepoEnt := repository_Ent.NewVuKhiRepositoryEnt(entClient)
	vukhiRepoGorm := repository_Gorm.NewVuKhiRepositoryGorm(gormClient)

	// khởi tạo service
	vukhiServiceEnt := service_Ent.NewVuKhiServiceEnt(vukhiRepoEnt)
	vukhiServiceGorm := service_Gorm.NewVuKhiServiceGorm(vukhiRepoGorm)

	// khởi tạo handler
	vukhiHandlerEnt := handler_Ent.NewVuKhiHandlerEnt(vukhiServiceEnt)
	vukhiHandlerGorm := handler_Gorm.NewVuKhiHandlerGorm(vukhiServiceGorm)

	// khởi tạo routes
	vukhiRoutesEnt := routes_Ent.NewVuKhiRoutesEnt(vukhiHandlerEnt)
	vukhiRoutesGorm := routes_Gorm.NewVuKhiRoutesGorm(vukhiHandlerGorm)

	return &VuKhiModule{
		routesEnt:  vukhiRoutesEnt,
		routesGorm: vukhiRoutesGorm,
	}
}

func (vk *VuKhiModule) RoutesEnt() []routes_Ent.Route {
	return []routes_Ent.Route{vk.routesEnt}
}

func (vk *VuKhiModule) RoutesGorm() []routes_Gorm.Route {
	return []routes_Gorm.Route{vk.routesGorm}
}
