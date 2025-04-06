package ui

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"auth/handler"
)

// UI 代表用户界面
type UI struct {
	handler *handler.DNSHandler
	done    chan struct{}
}

// NewUI 创建新的用户界面
func NewUI(h *handler.DNSHandler) *UI {
	return &UI{
		handler: h,
		done:    make(chan struct{}),
	}
}

// Run 启动UI
func (u *UI) Run() {
	// 启动一个goroutine定期更新状态
	go u.statusUpdater()

	// 启动一个goroutine监听新请求
	go u.requestListener()

	// 主输入循环
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("DNS服务器已启动。按Enter退出。")
	fmt.Println("当收到请求时，输入1返回NS响应，输入2返回TC=1响应。")

	for scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())

		if input == "" {
			continue
		}

		mode, err := strconv.Atoi(input)
		if err != nil || (mode != 1 && mode != 2) {
			fmt.Println("请输入1返回NS响应或2返回TC=1响应。")
			continue
		}

		if !u.handler.HasPendingRequest() {
			fmt.Println("没有待处理的请求。")
			continue
		}

		success := u.handler.SendResponse(mode)
		if success {
			fmt.Printf("已使用模式%d发送响应。\n", mode)
		} else {
			fmt.Println("发送响应失败。")
		}
	}

	close(u.done)
}

// Stop 停止UI
func (u *UI) Stop() {
	close(u.done)
}

// statusUpdater 定期更新状态显示
func (u *UI) statusUpdater() {
	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			u.printStatus()
		case <-u.done:
			return
		}
	}
}

// requestListener 监听新请求
func (u *UI) requestListener() {
	newRequestCh := u.handler.GetNewRequestChannel()

	for {
		select {
		case <-newRequestCh:
			fmt.Println("\n收到新请求！输入1返回NS响应，输入2返回TC=1响应。")
		case <-u.done:
			return
		}
	}
}

// printStatus 打印当前状态
func (u *UI) printStatus() {
	counter := u.handler.GetCounter()
	firstTime := u.handler.GetFirstReceivedTime()

	if firstTime.IsZero() {
		fmt.Printf("\r正在等待请求... 计数器: %d", counter)
		return
	}

	elapsed := time.Since(firstTime)
	fmt.Printf("\r请求数: %d | 自第一个请求以来的时间: %s", counter, formatDuration(elapsed))
}

// formatDuration 将持续时间格式化为可读字符串
func formatDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}