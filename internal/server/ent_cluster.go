package server

import (
	"game/internal/utils"
	"net/http"
	"strings"
)

func StartENTCluster(svcs *Services, ginHandlers []http.Handler) []*http.Server {
	servers := []*http.Server{
		TaoENTServer(&ENTServices{
			VuKhi:     svcs.VuKhi,
			He:        svcs.He,
			DoHiem:    svcs.DoHiem,
			LoaiVuKhi: svcs.LoaiVuKhi,
		}, ginHandlers[0]),
		TaoENTServer(&ENTServices{
			VuKhi:     svcs.VuKhi,
			He:        svcs.He,
			DoHiem:    svcs.DoHiem,
			LoaiVuKhi: svcs.LoaiVuKhi,
		}, ginHandlers[1]),
		TaoENTServer(&ENTServices{
			VuKhi:     svcs.VuKhi,
			He:        svcs.He,
			DoHiem:    svcs.DoHiem,
			LoaiVuKhi: svcs.LoaiVuKhi,
		}, ginHandlers[2]),
	}

	names := []string{"ENT-1", "ENT-2", "ENT-3"}
	ports := []string{
		utils.GetenvString("HTTP_PORT_ENT1", ":8060"),
		utils.GetenvString("HTTP_PORT_ENT2", ":8070"),
		utils.GetenvString("HTTP_PORT_ENT3", ":8080"),
	}

	for i, srv := range servers {
		port := strings.TrimSpace(ports[i])
		srv.Addr = port
		RunServer(names[i], port, srv)
	}

	return servers
}
