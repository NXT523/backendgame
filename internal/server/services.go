package server

import (
	"game/internal/service/grpc_service"
	"game/internal/utils"
)

type Services struct {
	VuKhi     *grpc_service.VuKhiService
	He        *grpc_service.HeService
	DoHiem    *grpc_service.DoHiemService
	LoaiVuKhi *grpc_service.LoaiVuKhiService
}

func InitServices(infra *Infra) *Services {
	lock := utils.NewRedisLock(infra.Redis)

	return &Services{
		VuKhi:     grpc_service.NewVuKhiService(infra.Ent, infra.Redis, lock),
		He:        grpc_service.NewHeService(infra.Ent, infra.Redis, lock),
		DoHiem:    grpc_service.NewDoHiemService(infra.Ent, infra.Redis, lock),
		LoaiVuKhi: grpc_service.NewLoaiVuKhiService(infra.Ent, infra.Redis, lock),
	}
}
