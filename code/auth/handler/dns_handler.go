package handler

import (
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/miekg/dns"

	"auth/config"
)

// RequestContext stores information about the latest DNS request
type RequestContext struct {
	Msg        *dns.Msg
	Writer     dns.ResponseWriter
	RemoteAddr net.Addr
	ReceivedAt time.Time
}

// DNSHandler handles DNS requests and maintains state
type DNSHandler struct {
	config        *config.Config
	counter       uint64           // Atomic counter to track the number of received packets
	latestRequest *RequestContext  // Context of the latest request
	firstReceived time.Time        // Time when the first packet was received
	mu            sync.Mutex       // Mutex to protect shared state
	newRequestCh  chan struct{}    // Channel to notify the UI of new requests
}

// NewDNSHandler creates a new DNS handler
func NewDNSHandler(cfg *config.Config) *DNSHandler {
	return &DNSHandler{
		config:       cfg,
		newRequestCh: make(chan struct{}, 1), // Buffered to 1 to prevent blocking
	}
}

// ServeDNS implements the dns.Handler interface
func (h *DNSHandler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	// Extract source IP
	remoteAddr := w.RemoteAddr()
	ipStr, _, err := net.SplitHostPort(remoteAddr.String())
	if err != nil {
		return // Discard malformed address
	}

	// Check if the request is from the target IP
	if ipStr != h.config.TargetIP {
		return // Discard if not from the target IP
	}

	// Check if it is an NS query
	if len(r.Question) == 0 || r.Question[0].Qtype != dns.TypeNS {
		return // Discard if not an NS query
	}

	// Increment the counter
	atomic.AddUint64(&h.counter, 1)

	h.mu.Lock()

	// Record the time if this is the first packet
	if h.firstReceived.IsZero() {
		h.firstReceived = time.Now()
	}

	// Update the latest request (discard the old one)
	h.latestRequest = &RequestContext{
		Msg:        r.Copy(), // Copy to prevent modification
		Writer:     w,
		RemoteAddr: remoteAddr,
		ReceivedAt: time.Now(),
	}

	h.mu.Unlock()

	// Notify the UI of a new request
	select {
	case h.newRequestCh <- struct{}{}:
		// Notification sent
	default:
		// Channel is full, but that's fine (UI already knows about the new request)
	}
}

// GetCounter returns the current packet counter value
func (h *DNSHandler) GetCounter() uint64 {
	return atomic.LoadUint64(&h.counter)
}

// GetFirstReceivedTime returns the time the first packet was received
func (h *DNSHandler) GetFirstReceivedTime() time.Time {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.firstReceived
}

// HasPendingRequest checks if there is a pending request
func (h *DNSHandler) HasPendingRequest() bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.latestRequest != nil
}

// GetNewRequestChannel returns the channel for new request notifications
func (h *DNSHandler) GetNewRequestChannel() <-chan struct{} {
	return h.newRequestCh
}

// SendResponse sends a response based on the specified mode
func (h *DNSHandler) SendResponse(mode int) bool {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.latestRequest == nil {
		return false // No request to respond to
	}

	req := h.latestRequest
	h.latestRequest = nil // Clear after responding

	var resp *dns.Msg

	switch mode {
	case 1:
		// Mode 1: Return NS records
		resp = h.createNSResponse(req.Msg)
	case 2:
		// Mode 2: Return TC=1 flag
		resp = h.createTCResponse(req.Msg)
	default:
		return false
	}

	// Send the response
	req.Writer.WriteMsg(resp)
	return true
}

// createNSResponse creates a DNS response with NS records
func (h *DNSHandler) createNSResponse(query *dns.Msg) *dns.Msg {
	m := new(dns.Msg)
	m.SetReply(query)
	m.Authoritative = true

	// Add an example NS record
	// In a real implementation, you might use the actual domain from the query
	if len(query.Question) > 0 {
		domain := query.Question[0].Name
		ns, _ := dns.NewRR(domain + " 3600 IN NS ns1." + domain)
		m.Ns = append(m.Ns, ns)

		// Add additional records as needed
		additional, _ := dns.NewRR("ns1." + domain + " 3600 IN A 192.0.2.1")
		m.Extra = append(m.Extra, additional)
	}

	return m
}

// createTCResponse creates a DNS response with the TC=1 flag
func (h *DNSHandler) createTCResponse(query *dns.Msg) *dns.Msg {
	m := new(dns.Msg)
	m.SetReply(query)
	m.Truncated = true // Set the TC=1 flag
	return m
}