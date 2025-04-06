package server

import (
	"fmt"
	"strconv"

	"github.com/miekg/dns"

	"auth/config"
)

// UDPServer represents a DNS server over UDP
type UDPServer struct {
	server *dns.Server
}

// NewUDPServer creates a new UDP DNS server
func NewUDPServer(cfg *config.Config, handler dns.Handler) *UDPServer {
	return &UDPServer{
		server: &dns.Server{
			Addr:    ":" + strconv.Itoa(cfg.ListenPort),
			Net:     "udp",
			Handler: handler,
		},
	}
}

// Start starts the UDP server
func (s *UDPServer) Start() error {
	fmt.Printf("Starting UDP DNS server, listening on %s\n", s.server.Addr)
	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *UDPServer) Shutdown() error {
	return s.server.Shutdown()
}