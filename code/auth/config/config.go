package config

// Config holds the application configuration
type Config struct {
	TargetIP   string // The target IP allowed to communicate with this server
	ListenPort int    // TCP and UDP listening port
}

// NewConfig creates the default configuration
func NewConfig() *Config {
	return &Config{
		TargetIP:   "10.152.0.71", // Default target IP, should be modified as needed
		ListenPort: 53,            // Standard DNS port
	}
}