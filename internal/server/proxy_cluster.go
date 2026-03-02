package server

import (
	"game/internal/utils"
	"game/pkg/proxy"
	"net/http"
)

// StartProxyServers khởi động các proxy server (LC, RR)
func StartProxyCluster(backends []string) []*http.Server {
	// Least Connections
	lcProxy := proxy.NewLeastConnServer("Proxy-LeastConn", backends, 100)

	// Round Robin
	rrProxy := proxy.NewRoundRobinServer("Proxy-RoundRobin", backends, 100)

	lcPort := utils.GetenvString("PROXY_PORT_LEAST_CONNECTIONS", "9000")
	rrPort := utils.GetenvString("PROXY_PORT_ROUND_ROBIN", "9010")

	RunServer("Proxy-LeastConn", lcPort, lcProxy)
	RunServer("Proxy-RoundRobin", rrPort, rrProxy)

	return []*http.Server{
		lcProxy,
		rrProxy,
	}
}

// return các server ent
func BuildBackendENTs() []string {
	return []string{
		"http://localhost" + utils.GetenvString("HTTP_PORT_ENT1", ":8060"),
		"http://localhost" + utils.GetenvString("HTTP_PORT_ENT2", ":8070"),
		"http://localhost" + utils.GetenvString("HTTP_PORT_ENT3", ":8080"),
	}
}
