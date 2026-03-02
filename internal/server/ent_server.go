package server

import (
	"context"
	grpcsrv "game/internal/grpc"
	"game/internal/middleware"
	"game/internal/service/grpc_service"
	v1 "game/v1/proto"
	"net/http"
	"time"

	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

type ENTServices struct {
	VuKhi     *grpc_service.VuKhiService
	He        *grpc_service.HeService
	DoHiem    *grpc_service.DoHiemService
	LoaiVuKhi *grpc_service.LoaiVuKhiService
}

func TaoENTServer(svcs *ENTServices, ginHandler http.Handler) *http.Server {

	//Nếu client ping < 10s/lần → server sẽ đóng connection. Tránh client spam ping gây tốn tài nguyên.
	kaep := keepalive.EnforcementPolicy{
		MinTime:             10 * time.Second, // Chỉ yêu cầu client ping mỗi 10s (đỡ gắt hơn 5s)
		PermitWithoutStream: true,
	}

	kasp := keepalive.ServerParameters{
		MaxConnectionIdle: 5 * time.Minute,  // Idle 5 phút
		MaxConnectionAge:  1 * time.Hour,    // Sống 1 tiếng mới reset
		Time:              20 * time.Second, // Server ping client mỗi 20s
		Timeout:           5 * time.Second,
	}

	grpcSrv := grpc.NewServer(
		grpc.MaxConcurrentStreams(1000),
		grpc.KeepaliveEnforcementPolicy(kaep),
		grpc.KeepaliveParams(kasp),

		// Middleware
		grpc.ChainUnaryInterceptor(
			middleware.GRPCUnaryLoggerMiddleware(),
			grpc_recovery.UnaryServerInterceptor(),
		),
		grpc.ChainStreamInterceptor(
			middleware.GRPCStreamLoggerMiddleware(),
			grpc_recovery.StreamServerInterceptor(),
		),
	)

	vukhiImpl := grpcsrv.NewVuKhiGRPCServer(svcs.VuKhi)
	heImpl := grpcsrv.NewHeGRPCServer(svcs.He)
	dohiemImpl := grpcsrv.NewDoHiemGRPCServer(svcs.DoHiem)
	loaivukhiImpl := grpcsrv.NewLoaiVuKhiGRPCServer(svcs.LoaiVuKhi)

	v1.RegisterVuKhiServiceServer(grpcSrv, vukhiImpl)
	v1.RegisterHeServiceServer(grpcSrv, heImpl)
	v1.RegisterDoHiemServiceServer(grpcSrv, dohiemImpl)
	v1.RegisterLoaiVuKhiServiceServer(grpcSrv, loaivukhiImpl)

	reflection.Register(grpcSrv)

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
	)

	ctx := context.Background()
	// Register Handler Server (In-Process)
	if err := v1.RegisterVuKhiServiceHandlerServer(ctx, gw, vukhiImpl); err != nil {
		panic(err)
	}
	if err := v1.RegisterHeServiceHandlerServer(ctx, gw, heImpl); err != nil {
		panic(err)
	}
	if err := v1.RegisterDoHiemServiceHandlerServer(ctx, gw, dohiemImpl); err != nil {
		panic(err)
	}
	if err := v1.RegisterLoaiVuKhiServiceHandlerServer(ctx, gw, loaivukhiImpl); err != nil {
		panic(err)
	}

	return GatewayBuildH2CServer(grpcSrv, gw, ginHandler)
}
