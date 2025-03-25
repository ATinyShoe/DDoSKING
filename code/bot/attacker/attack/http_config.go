// http_config.go
package attack

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"sync/atomic"
	"time"
)

// addHeaders Adds HTTP request headers - Optimized to ensure headers are correctly added
func (h *HTTP) addHeaders(req *http.Request) {
	// Use custom headers if available
	if h.Header != nil && len(h.Header) > 0 {
		for k, v := range h.Header {
			// Set headers
			req.Header.Set(k, v)
		}
		// Add a random ID to prevent detection as attack traffic
		req.Header.Set("X-Request-ID", generateRandomID())
		return
	}

	// Otherwise, use default headers
	for k, v := range GetDefaultHeaders() {
		req.Header.Set(k, v)
	}
}

// generateRandomID Generates a random ID to increase request randomness
func generateRandomID() string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	id := make([]byte, 16)
	for i := range id {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		id[i] = chars[n.Int64()]
	}
	return string(id)
}

// GetDefaultHeaders Retrieves default HTTP request headers
func GetDefaultHeaders() map[string]string {
	// Add more commonly used headers to make traffic appear more realistic
	return map[string]string{
		"User-Agent":                GetRandomUserAgent(),
		"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8",
		"Accept-Language":           GetRandomLanguage(),
		"Accept-Encoding":           "gzip, deflate, br",
		"Connection":                GetRandomConnection(),
		"Cache-Control":             "no-cache",
		"Pragma":                    "no-cache",
		"Sec-Fetch-Dest":            "document",
		"Sec-Fetch-Mode":            "navigate",
		"Sec-Fetch-Site":            "none",
		"Sec-Fetch-User":            "?1",
		"Upgrade-Insecure-Requests": "1",
	}
}

// GetRandomConnection Randomly returns a connection type
func GetRandomConnection() string {
	connections := []string{
		"keep-alive",
		"close",
	}
	idx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(connections))))
	return connections[idx.Int64()]
}

// GetRandomLanguage Retrieves a random language setting
func GetRandomLanguage() string {
	languages := []string{
		"en-US,en;q=0.9",
		"en-GB,en;q=0.8,en-US;q=0.7",
		"zh-CN,zh;q=0.9,en;q=0.8",
		"zh-TW,zh;q=0.9,en-US;q=0.8,en;q=0.7",
		"ja-JP,ja;q=0.9,en-US;q=0.8,en;q=0.7",
		"ko-KR,ko;q=0.9,en-US;q=0.8,en;q=0.7",
		"de-DE,de;q=0.9,en;q=0.8",
		"fr-FR,fr;q=0.9,en;q=0.8",
		"ru-RU,ru;q=0.9,en;q=0.8",
	}
	idx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(languages))))
	return languages[idx.Int64()]
}

// GetRandomUserAgent Retrieves a random user agent - Adds more modern browser UAs
func GetRandomUserAgent() string {
	userAgents := []string{
		// Chrome
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36",
		// Firefox
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/117.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:109.0) Gecko/20100101 Firefox/117.0",
		"Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/117.0",
		// Safari
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.5 Safari/605.1.15",
		// Edge
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36 Edg/115.0.1901.200",
		// Mobile
		"Mozilla/5.0 (iPhone; CPU iPhone OS 16_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (Linux; Android 13; SM-S918B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Mobile Safari/537.36",
	}

	// Randomly select a user agent
	idx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(userAgents))))
	return userAgents[idx.Int64()]
}

// newClient Creates an HTTP client - Adds connection pooling and more configuration options
func (h *HTTP) newClient() *http.Client {
	// Create custom transport configuration
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,             // Ignore SSL certificate errors
			MinVersion:         tls.VersionTLS12, // Use secure TLS versions
		},
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			dialer := &net.Dialer{
				Timeout:   5 * time.Second,  // Connection timeout
				KeepAlive: 30 * time.Second, // TCP keepalive
			}
			conn, err := dialer.DialContext(ctx, network, addr)
			return conn, err
		},
		MaxIdleConns:          1000,             // Maximum idle connections
		MaxConnsPerHost:       100,              // Maximum connections per host
		MaxIdleConnsPerHost:   100,              // Maximum idle connections per host
		IdleConnTimeout:       90 * time.Second, // Idle connection timeout
		TLSHandshakeTimeout:   10 * time.Second, // TLS handshake timeout
		ExpectContinueTimeout: 1 * time.Second,  // 100-continue timeout
		DisableCompression:    true,             // Disable compression for faster requests
		DisableKeepAlives:     false,            // Enable KeepAlive
		ForceAttemptHTTP2:     true,             // Attempt HTTP/2
	}

	return &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second, // Overall request timeout
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // Do not follow redirects
		},
	}
}

// printSummary Prints attack statistics
func (h *HTTP) printSummary() {
	endTime := time.Now()
	duration := endTime.Sub(h.startTime)

	// Calculate requests per second
	requestsPerSecond := float64(0)
	if duration.Seconds() > 0 {
		requestsPerSecond = float64(atomic.LoadInt64(&h.sentRequests)) / duration.Seconds()
	}

	// Calculate success rate
	successRate := float64(0)
	sentReqs := atomic.LoadInt64(&h.sentRequests)
	if sentReqs > 0 {
		successRate = float64(atomic.LoadInt64(&h.recvResponses)) / float64(sentReqs) * 100
	}

	fmt.Printf("\n=== HTTP Attack Statistics ===\n")
	fmt.Printf("Start Time: %s\n", h.startTime.Format("2006/01/02 15:04:05"))
	fmt.Printf("End Time: %s\n", endTime.Format("2006/01/02 15:04:05"))
	fmt.Printf("Target: %s\n", h.Target)
	fmt.Printf("Request Method: %s\n", h.Method)
	fmt.Printf("Threads: %d\n", h.Threads)
	fmt.Printf("Requests Sent: %d (%.2f requests/sec)\n",
		sentReqs,
		requestsPerSecond)
	fmt.Printf("Responses Received: %d (Success Rate: %.2f%%)\n",
		atomic.LoadInt64(&h.recvResponses),
		successRate)
	fmt.Printf("Data Sent: %.2f MB, Data Received: %.2f MB\n",
		float64(atomic.LoadInt64(&h.sentBytes))/(1024*1024),
		float64(atomic.LoadInt64(&h.recvBytes))/(1024*1024))
	fmt.Printf("Connection Errors: %d\n",
		atomic.LoadInt64(&h.connErrors))
	fmt.Printf("====================\n")
}

// estimateRequestSize Estimates the size of an HTTP request
func estimateRequestSize(req *http.Request) int {
	size := len(req.Method) + len(req.URL.Path) + len("HTTP/1.1") + 4 // Request line

	// Headers
	for name, values := range req.Header {
		for _, value := range values {
			size += len(name) + len(value) + 4 // Includes colon, space, and CRLF
		}
	}

	size += 2 // Blank line between headers and body

	return size
}
