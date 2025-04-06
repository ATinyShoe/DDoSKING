package server

import (
	"fmt"
	"strconv"

	"github.com/miekg/dns"

	"auth/config"
)

// TCPServer represents a DNS server over TCP
type TCPServer struct {
	server *dns.Server
}

// NewTCPServer creates a new TCP DNS server
func NewTCPServer(cfg *config.Config, handler dns.Handler) *TCPServer {
	return &TCPServer{
		server: &dns.Server{
			Addr:    ":" + strconv.Itoa(cfg.ListenPort),
			Net:     "tcp",
			Handler: handler,
		},
	}
}

// Start starts the TCP server
func (s *TCPServer) Start() error {
	fmt.Printf("Starting TCP DNS server, listening on %s\n", s.server.Addr)
	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *TCPServer) Shutdown() error {
	return s.server.Shutdown()
}