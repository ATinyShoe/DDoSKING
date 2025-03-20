package attack

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync/atomic"
)

// HTTP GET方法
// 需要优化发包策略，现在发包是需要等待的
func (h *HTTP) httpGet() error {
	client := h.newClient()

	// 构建请求目标
	reqTarget := h.Target

	// 使用提供的Path (如果有)，否则就访问根目录
	if h.Path != "" {
		parsedURL, err := url.Parse(h.Target)
		if err == nil {
			parsedURL.Path = h.Path
			reqTarget = parsedURL.String()
		}
	}

	req, err := http.NewRequest("GET", reqTarget, nil)
	if err != nil {
		return err
	}

	// 添加请求头
	for k, v := range h.generateHeaders() {
		req.Header.Set(k, v)
	}

	// 记录发送请求
	atomic.AddInt64(&h.sentRequests, 1)
	atomic.AddInt64(&h.sentBytes, int64(estimateRequestSize(req)))

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	// 记录接收响应
	atomic.AddInt64(&h.recvResponses, 1)

	if resp.Body != nil {
		// 读取并丢弃响应体，但统计大小
		n, _ := io.Copy(io.Discard, resp.Body)
		atomic.AddInt64(&h.recvBytes, n)
		resp.Body.Close()
	}

	return nil
}

// HTTP POST方法
func (h *HTTP) httpPost() error {
	client := h.newClient()

	// 构建请求目标
	reqTarget := h.Target

	// 使用提供的Path (如果有)，否则就访问根目录
	if h.Path != "" {
		parsedURL, err := url.Parse(h.Target)
		if err == nil {
			parsedURL.Path = h.Path
			reqTarget = parsedURL.String()
		}
	}

	// 获取请求体内容
	payload := h.Payload

	// 如果没有指定负载，根据方法生成
	if payload == "" {
		payload = GetPayload(h.Method)
	}

	// 创建请求体
	body := strings.NewReader(payload)

	// 创建请求
	req, err := http.NewRequest("POST", reqTarget, body)
	if err != nil {
		return err
	}

	// 添加请求头
	for k, v := range h.generateHeaders() {
		req.Header.Set(k, v)
	}

	// 设置Content-Type为application/json（除非头部已经设置）
	if _, exists := req.Header["Content-Type"]; !exists {
		req.Header.Set("Content-Type", "application/json")
	}

	// 记录发送请求
	atomic.AddInt64(&h.sentRequests, 1)
	atomic.AddInt64(&h.sentBytes, int64(len(payload)+estimateRequestSize(req)))

	go func() {
		resp, err := client.Do(req)
		if err != nil {
			atomic.AddInt64(&h.connErrors, 1)
			return
		}

		// 记录接收响应
		atomic.AddInt64(&h.recvResponses, 1)

		if resp.Body != nil {
			// 读取并丢弃响应体，但统计大小
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
		}
	}()

	return nil
}

// HTTP Cookie方法 - 发送大量Cookie的GET请求
func (h *HTTP) httpCookie() error {
	client := h.newClient()

	// 构建请求目标
	reqTarget := h.Target

	// 使用提供的Path (如果有)，否则就访问根目录
	if h.Path != "" {
		parsedURL, err := url.Parse(h.Target)
		if err == nil {
			parsedURL.Path = h.Path
			reqTarget = parsedURL.String()
		}
	}

	req, err := http.NewRequest("GET", reqTarget, nil)
	if err != nil {
		return err
	}

	// 添加请求头
	for k, v := range h.generateHeaders() {
		req.Header.Set(k, v)
	}

	// 添加大量Cookie
	cookieStr := generateCookieString()
	req.Header.Set("Cookie", cookieStr)

	// 记录发送请求
	atomic.AddInt64(&h.sentRequests, 1)
	atomic.AddInt64(&h.sentBytes, int64(estimateRequestSize(req)))

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	// 记录接收响应
	atomic.AddInt64(&h.recvResponses, 1)

	if resp.Body != nil {
		// 读取并丢弃响应体，但统计大小
		n, _ := io.Copy(io.Discard, resp.Body)
		atomic.AddInt64(&h.recvBytes, n)
		resp.Body.Close()
	}

	return nil
}

// 估算HTTP请求大小
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
