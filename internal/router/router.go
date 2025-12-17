package router

import (
	"database/sql"
	"game/ent"
	"game/internal/repo"
	"game/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 1. ENT
func CreateRouterEnt(client *ent.Client) *gin.Engine {
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true, "engine": "ent"}) })

	vkRepo := repo.CreateKhoVuKhiEnt(client)
	vkService := service.CreateVuKhiServiceEnt(vkRepo)

	lvkRepo := repo.CreateKhoLoaiVuKhiEnt(client)
	lvkService := service.CreateLoaiVuKhiServiceEnt(lvkRepo)

	heRepo := repo.CreateKhoHeEnt(client)
	heService := service.CreateHeServiceEnt(heRepo)

	dhRepo := repo.CreateKhoDoHiemEnt(client)
	dhService := service.CreateDoHiemServiceEnt(dhRepo)

	MountVuKhiEnt(r, vkService)
	MountLoaiVuKhiEnt(r, lvkService)
	MountHeEnt(r, heService)
	MountDoHiemEnt(r, dhService)

	return r
}

// 2. GORM
func CreateRouterGorm(gdb *gorm.DB) *gin.Engine {
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true, "engine": "gorm"}) })

	// DoHiem
	dhRepo := repo.CreateKhoDoHiemGorm(gdb)
	dhService := service.CreateDoHiemServiceGorm(dhRepo)

	// He
	heRepo := repo.CreateKhoHeGorm(gdb)
	heService := service.CreateHeServiceGorm(heRepo)

	// LoaiVuKhi
	lvkRepo := repo.CreateKhoLoaiVuKhiGorm(gdb)
	lvkService := service.CreateLoaiVuKhiServiceGorm(lvkRepo)

	// VuKhi
	vkRepo := repo.CreateKhoVuKhiGorm(gdb)
	vkService := service.CreateVuKhiServiceGorm(vkRepo)

	MountHeGorm(r, heService)
	MountLoaiVuKhiGorm(r, lvkService)
	MountVuKhiGorm(r, vkService)
	MountDoHiemGorm(r, dhService)
	return r
}

// 3. RAW
func CreateRouterRaw(sqlDB *sql.DB) *gin.Engine {
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true, "engine": "raw"}) })

	// DoHiem
	dhRepo := repo.CreateKhoDoHiemRaw(sqlDB)
	dhService := service.CreateDoHiemServiceRaw(dhRepo)

	// He
	heRepo := repo.CreateKhoHeRaw(sqlDB)
	heService := service.CreateHeServiceRaw(heRepo)

	// VuKhi
	vkRepo := repo.CreateKhoVuKhiRaw(sqlDB)
	vkService := service.CreateVuKhiServiceRaw(vkRepo)

	lvkRepo := repo.CreateKhoLoaiVuKhiRaw(sqlDB)
	lvkService := service.CreateLoaiVuKhiServiceRaw(lvkRepo)

	MountLoaiVuKhiRaw(r, lvkService)
	MountDoHiemRaw(r, dhService)
	MountHeRaw(r, heService)
	MountVuKhiRaw(r, vkService)
	return r
}
