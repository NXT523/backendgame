package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"game/internal/app"
	server "game/internal/server"
	"game/internal/utils"
	kconsumer "game/pkg/kafka"

	"github.com/joho/godotenv"
)

func main() {
	log.Println("Server start............................")
	// code của thư viện gin
	vkModule := app.NewVuKhiModule()

	entApp := app.NewApplicationEnt(utils.GetenvString("HTTP_PORT_ENT3", ":8080"), []app.EntModule{vkModule})
	gormApp := app.NewApplicationGorm(utils.GetenvString("HTTP_PORT_GORM", ":8081"), []app.GormModule{vkModule})

	go func() {
		if err := gormApp.Run(); err != nil {
			log.Fatal("GORM gin server error:", err)
		}
	}()
	// code cửa thư viện gin

	ctx, cancelAll := context.WithCancel(context.Background())
	defer cancelAll()

	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ không tìm thấy .env")
	}

	infra, err := server.InitInfra(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer infra.Close()

	// Services & servers
	svcs := server.InitServices(infra)
	handlers := []http.Handler{
		entApp.GetHandler(),
		entApp.GetHandler(),
		entApp.GetHandler(),
	}
	entServers := server.StartENTCluster(svcs, handlers)
	proxyServers := server.StartProxyCluster(server.BuildBackendENTs())

	// Kafka consumers
	for i := 0; i < 3; i++ {
		go kconsumer.ChayConsumerAuditVuKhi(ctx, infra.Elastic)
	}

	// Wait signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("⏳ Đang tắt server...")
	cancelAll()

	// Shutdown HTTP servers
	shutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.ShutdownAll(shutCtx, append(entServers, proxyServers...)...)

	// Flush elastic
	_ = infra.Elastic.DongBo(context.Background())

	log.Println("✅ Tắt server thành công. Hẹn gặp lại 👋")
}
