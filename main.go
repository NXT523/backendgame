package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"game/ent"
	"game/internal/config"
	"game/internal/db"
	grpcsrv "game/internal/grpc"
	"game/internal/router"
	"game/pkg/elastic"
	"game/pkg/interceptor"
	kconsumer "game/pkg/kafka"
	"game/pkg/logx"
	"game/pkg/memcached"
	"game/pkg/middleware"
	"game/pkg/proxy"
	redisx "game/pkg/redis"
	v1 "game/v1"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

type statusWriter struct {
	http.ResponseWriter
	status   int
	bytesOut int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}
func (w *statusWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.bytesOut += n
	return n, err
}

func taoServerENTChungCong(
	cfg config.CauHinh,
	entClient *ent.Client,
	rdb *redis.Client,
	mc *memcache.Client,
	kw *kconsumer.KafkaWriters, // có thể nil
	elog *logx.LoggerElastic, // nhận logger từ main
) *http.Server {
	rl := middleware.TaoRateLimiter(3, time.Minute)
	grpcSrv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logx.GRPCGhiLog(elog),
			middleware.RateLimitGrpc(rl),
			interceptor.MasterInterceptor(),
		),
	)

	heImpl := grpcsrv.NewHeGRPCServerMemcached(entClient, mc, 10*time.Minute)
	dohiemImpl := grpcsrv.NewDoHiemGRPCServerMemcached(entClient, mc, 10*time.Minute)
	loaivukhiImpl := grpcsrv.NewLoaiVuKhiGRPCServer(entClient)

	var vukhiImpl *grpcsrv.VuKhiGRPCServer
	if kw != nil {
		vukhiImpl = grpcsrv.NewVuKhiGRPCServerKafka(entClient, rdb, kw.Created, kw.Updated)
	} else {
		vukhiImpl = grpcsrv.NewVuKhiGRPCServer(entClient, rdb)
	}

	v1.RegisterHeServiceServer(grpcSrv, heImpl)
	v1.RegisterVuKhiServiceServer(grpcSrv, vukhiImpl)
	v1.RegisterDoHiemServiceServer(grpcSrv, dohiemImpl)
	v1.RegisterLoaiVuKhiServiceServer(grpcSrv, loaivukhiImpl)
	reflection.Register(grpcSrv)

	restGin := router.CreateRouterEnt(entClient)
	restGin.Use(logx.GinGhiLogTruyCap(elog))

	gw := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames:   true,
				EmitUnpopulated: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: false,
			},
		}),
		runtime.WithErrorHandler(logx.HttpErrorHandler(elog)),
	)

	ctx := context.Background()
	if err := v1.RegisterHeServiceHandlerServer(ctx, gw, heImpl); err != nil {
		log.Fatalf("[GW] Đăng ký He handler lỗi: %v", err)
	}
	if err := v1.RegisterVuKhiServiceHandlerServer(ctx, gw, vukhiImpl); err != nil {
		log.Fatalf("[GW] Đăng ký VuKhi handler lỗi: %v", err)
	}
	if err := v1.RegisterDoHiemServiceHandlerServer(ctx, gw, dohiemImpl); err != nil {
		log.Fatalf("[GW] Đăng ký DoHiem handler lỗi: %v", err)
	}
	if err := v1.RegisterLoaiVuKhiServiceHandlerServer(ctx, gw, loaivukhiImpl); err != nil {
		log.Fatalf("[GW] Đăng ký LoaiVuKhi handler lỗi: %v", err)
	}

	gwWithLog := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sw := &statusWriter{ResponseWriter: w, status: 200}
		start := time.Now()

		r2 := r.WithContext(context.WithValue(r.Context(), logx.CtxKeyHTTPStart, start))
		gw.ServeHTTP(sw, r2)

		if sw.status < 400 {
			ip := r.Header.Get("X-Forwarded-For")
			if ip == "" {
				host, _, err := net.SplitHostPort(r.RemoteAddr)
				if err == nil && host != "" {
					ip = host
				} else {
					ip = r.RemoteAddr
				}
			}

			kv := map[string]any{
				"service.component":         "gateway",
				"event.dataset":             "http",
				"http.request.method":       r.Method,
				"http.request.url":          r.URL.Path,
				"http.response.status_code": sw.status,
				"client.ip":                 ip,
				"duration_ms":               time.Since(start).Milliseconds(),
			}
			if tid := r.Header.Get("X-Trace-Id"); tid != "" {
				kv["trace_id"] = tid
			}
			if rid := r.Header.Get("X-Request-Id"); rid != "" {
				kv["request_id"] = rid
			}

			elog.Ghi(r.Context(), "info", "HTTP request handled", kv)
		}
	})

	mux := http.NewServeMux()
	mux.Handle("/he/ent/", gwWithLog)
	mux.Handle("/vukhi/ent/", gwWithLog)
	mux.Handle("/dohiem/ent/", gwWithLog)
	mux.Handle("/loaivukhi/ent/", gwWithLog)
	mux.Handle("/", restGin)

	h2cHandler := h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.HasPrefix(r.Header.Get("Content-Type"), "application/grpc") {
			grpcSrv.ServeHTTP(w, r)
			return
		}
		mux.ServeHTTP(w, r)
	}), &http2.Server{})

	// rateLimiter := middleware.TaoRateLimiter(100, time.Minute) // 100 req/phút từ mỗi IP

	// Trả server với Handler; Addr sẽ được gán từ main
	return &http.Server{
		// Handler: rateLimiter.RateLimitHTTP(h2cHandler),
		Handler: h2cHandler,
	}
}

func main() {

	if err := grpcsrv.SetupAccessLogger(); err != nil {
		log.Fatalf("Setup access logger lỗi: %v", err)
	}
	cfg := config.DocCauHinh()
	ctx := context.Background()

	// ===== ELASTIC LOGGER =====
	elog, err := elastic.TaoElasticLogger(ctx, cfg)
	if err != nil {
		log.Fatalf("Không tạo được logger ES: %v", err)
	}

	// ===== REDIS =====
	rdb, err := redisx.NewRedisClient(ctx, cfg)
	if err != nil {
		log.Fatalf("Không thể kết nối Redis: %v", err)
	}
	log.Println("Kết nối Redis thành công")

	// ===== MEMCACHED =====
	mc, err := memcached.NewMemcacheClient(ctx, cfg)
	if err != nil {
		log.Fatalf("Không thể kết nối Memcached: %v", err)
	}
	log.Println("Kết nối Memcached thành công")

	// ===== ENT =====
	if err := db.CreateDBEnt(cfg.DSNChuaDB, cfg.TenDBEnt); err != nil {
		log.Fatalf("[ENT] Tạo DB lỗi: %v", err)
	}
	entClient, err := db.OpenEnt(cfg.DSNChuaDB, cfg.TenDBEnt)
	if err != nil {
		log.Fatalf("[ENT] Kết nối MySQL lỗi: %v", err)
	}
	if err := db.RunMigrationEnt(context.Background(), entClient); err != nil {
		log.Fatalf("[ENT] Migration lỗi: %v", err)
	}

	// ===== KAFKA WRITER =====
	kafkaWriter, err := kconsumer.TaoKafkaWriters(ctx, cfg)
	if err != nil {
		log.Fatalf("Khởi tạo Kafka writers lỗi: %v", err)
	}
	defer kafkaWriter.Created.Close()
	defer kafkaWriter.Updated.Close()

	// ===== KAFKA CONSUMERS =====
	go kconsumer.ChayConsumerAuditVuKhi(ctx, elog)
	go kconsumer.ChayConsumerAuditVuKhi(ctx, elog)
	go kconsumer.ChayConsumerAuditVuKhi(ctx, elog)

	// ===== OUTBOX WORKERS =====
	workerCount := 3 // số worker cho EventVuKhiCreated (match partitions)
	workerCtx, workerCancel := context.WithCancel(ctx)
	defer workerCancel()
	var workerWg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		workerWg.Add(1)
		go func() {
			defer workerWg.Done()
			kconsumer.SenderWorker(
				workerCtx,
				rdb,
				kconsumer.EventVuKhiCreated,
				"outbox:vukhicreate",
				kafkaWriter.Created,
				elog,
			)
		}()
	}

	// ===== ENT SERVER (có Kafka & logger) =====
	entSrv1 := taoServerENTChungCong(cfg, entClient, rdb, mc, kafkaWriter, elog)
	entSrv2 := taoServerENTChungCong(cfg, entClient, rdb, mc, kafkaWriter, elog)
	entSrv3 := taoServerENTChungCong(cfg, entClient, rdb, mc, kafkaWriter, elog)

	// Gán Addr từ config (HTTPPortEnt1/2/3)
	// Nếu config không set đầy đủ, fallback sẽ reuse port cuối cùng
	ports := []string{cfg.HTTPPortEnt1, cfg.HTTPPortEnt2, cfg.HTTPPortEnt3}
	for i := range ports {
		ports[i] = strings.TrimSpace(ports[i])
		if ports[i] == "" {
			ports[i] = ":0" // nếu rỗng, allow net.Listen lắng nghe random port (hiếm khi dùng)
		}
	}
	entSrv1.Addr = ports[0]
	entSrv2.Addr = ports[1]
	entSrv3.Addr = ports[2]

	// ===== GORM =====
	if err := db.CreateDBGorm(cfg.DSNChuaDB, cfg.TenDBGorm); err != nil {
		log.Fatalf("[GORM] Tạo DB lỗi: %v", err)
	}
	gdb, err := db.OpenGorm(cfg.DSNChuaDB, cfg.TenDBGorm)
	if err != nil {
		log.Fatalf("[GORM] Kết nối MySQL lỗi: %v", err)
	}
	if err := db.RunMigrationGorm(gdb); err != nil {
		log.Fatalf("[GORM] Migration lỗi: %v", err)
	}
	gormSrv := &http.Server{Addr: cfg.HTTPPortGorm, Handler: router.CreateRouterGorm(gdb)}

	// ===== RAW =====
	if err := db.CreateDBRaw(cfg.DSNChuaDB, cfg.TenDBRaw); err != nil {
		log.Fatalf("[RAW] Tạo DB lỗi: %v", err)
	}
	rawDB, err := db.OpenRaw(cfg.DSNChuaDB, cfg.TenDBRaw)
	if err != nil {
		log.Fatalf("[RAW] Kết nối MySQL lỗi: %v", err)
	}
	if err := db.RunMigrationRaw(rawDB); err != nil {
		log.Fatalf("[RAW] Migration lỗi: %v", err)
	}
	rawSrv := &http.Server{Addr: cfg.HTTPPortRaw, Handler: router.CreateRouterRaw(rawDB)}

	// ===== REVERSE PROXY - ROUND ROBIN =====
	rrRegistry := proxy.NewUpstreamRegistry()

	backendURLsRoundRobin := []string{
		"http://localhost" + cfg.HTTPPortEnt1,
		"http://localhost" + cfg.HTTPPortEnt2,
		"http://localhost" + cfg.HTTPPortEnt3,
	}

	for _, u := range backendURLsRoundRobin {
		up, err := proxy.NewUpstream(u, 100)
		if err != nil {
			log.Fatalf("Không parse upstream %s: %v", u, err)
		}
		rrRegistry.Register(up)
		log.Printf("[RR] Added upstream: %s", u)
	}

	rrLB := proxy.NewRoundRobin(rrRegistry)
	// rateLimiter := middleware.TaoRateLimiter(100, time.Minute)

	// proxyHandler :=
	// 	middleware.InjectionGuard(
	// 		rateLimiter.RateLimitHTTP(lb),
	// 	)
	proxySrvRoundRobin := &http.Server{
		Addr:    cfg.ProxyPortRoundRobin,
		Handler: rrLB,
	}

	// ===== REVERSE PROXY - LEAST CONNECTIONS =====
	lcPool := proxy.NewBackendPool()

	backendURLs := []string{
		"http://localhost" + cfg.HTTPPortEnt1,
		"http://localhost" + cfg.HTTPPortEnt2,
		"http://localhost" + cfg.HTTPPortEnt3,
	}

	for _, u := range backendURLs {
		backend, err := proxy.NewBackendNode(u, 100) // max 100 concurrent
		if err != nil {
			log.Fatalf("Không parse backend %s: %v", u, err)
		}
		lcPool.AddBackend(backend)
		log.Printf("[LC] Added backend: %s", u)
	}

	lcLB := proxy.NewSmartLoadBalancer(lcPool)

	proxySrvLeastConnections := &http.Server{
		Addr:    cfg.ProxyPortLeastConnections,
		Handler: lcLB,
	}
	// Start servers song song (ENT, GORM, RAW, Proxy)
	go func() {
		log.Printf("🚀 Least-Connections Proxy chạy tại http://localhost%s", cfg.ProxyPortLeastConnections)
		if err := proxySrvLeastConnections.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[LC PROXY] lỗi: %v", err)
		}
	}()
	go func() {
		log.Printf("🚀 Round Robin Proxy chạy tại http://localhost%s", cfg.ProxyPortRoundRobin)
		if err := proxySrvRoundRobin.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[PROXY] lỗi: %v", err)
		}
	}()

	go func() {
		log.Println("🚀 ENT1  chạy tại http://localhost" + cfg.HTTPPortEnt1)
		if err := entSrv1.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[ENT] ListenAndServe lỗi: %v", err)
		}
	}()
	go func() {
		log.Println("🚀 ENT2  chạy tại http://localhost" + cfg.HTTPPortEnt2)
		if err := entSrv2.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[ENT] ListenAndServe lỗi: %v", err)
		}
	}()
	go func() {
		log.Println("🚀 ENT3  chạy tại http://localhost" + cfg.HTTPPortEnt3)
		if err := entSrv3.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[ENT] ListenAndServe lỗi: %v", err)
		}
	}()
	go func() {
		log.Println("🚀 GORM chạy tại http://localhost" + cfg.HTTPPortGorm)
		if err := gormSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[GORM] ListenAndServe lỗi: %v", err)
		}
	}()
	go func() {
		log.Println("🚀 RAW  chạy tại http://localhost" + cfg.HTTPPortRaw)
		if err := rawSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[RAW] ListenAndServe lỗi: %v", err)
		}
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	log.Println("⏳ Đang tắt server...")

	// 1) Cancel outbox workers first
	workerCancel()
	log.Println("⏳ Đã gửi cancel tới outbox workers")

	// 1.1) Force-unblock BLPOP bằng cách đóng Redis client
	if err := rdb.Close(); err != nil {
		log.Printf("Đóng redis client lỗi: %v", err)
	} else {
		log.Println("Đã đóng Redis client để unblock BLPOP")
	}

	// 1.2) Wait for workers with timeout
	done := make(chan struct{})
	go func() {
		workerWg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Println("✅ Outbox workers đã dừng")
	case <-time.After(10 * time.Second):
		log.Println("⚠️ Timeout chờ outbox workers; tiếp tục shutdown (workers có thể vẫn đang chạy)")
	}

	// 2) Shutdown HTTP servers (5s timeout)
	shutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	

	_ = entSrv1.Shutdown(shutCtx)
	_ = entSrv2.Shutdown(shutCtx)
	_ = entSrv3.Shutdown(shutCtx)
	_ = gormSrv.Shutdown(shutCtx)
	_ = rawSrv.Shutdown(shutCtx)
	_ = proxySrvLeastConnections.Shutdown(shutCtx)
	_ = proxySrvRoundRobin.Shutdown(shutCtx)

	// 4) Close DBs, flush logs
	entClient.Close()
	if sqlDB, err := gdb.DB(); err == nil {
		_ = sqlDB.Close()
	}
	_ = rawDB.Close()

	_ = elog.DongBo(context.Background())

	log.Println("✅ Đã tắt xong.")
}
