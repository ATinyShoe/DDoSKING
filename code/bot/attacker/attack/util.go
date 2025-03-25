package attack

import (
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/time/rate"
)

var bandwidthLimiter *rate.Limiter

// Set random seed
func init() {
	rand.Seed(time.Now().UnixNano())
}

// Generate random IP address
func RandIPv4() string {
	rand.Seed(time.Now().UnixNano()) // Set random seed
	return fmt.Sprintf("%d.%d.%d.%d",
		rand.Intn(256), // Each segment ranges from 0-255
		rand.Intn(256),
		rand.Intn(256),
		rand.Intn(256),
	)
}

// Generate random port
func RandPort() int {
	return rand.Intn(65535-1) + 1 // Generate a random port between 1 and 65535
}

// ResetStopChannel resets the stop channel
func ResetStopChannel() {
	// Reset the stop channel
	select {
	case <-STOP:
		// Channel is closed, needs to be recreated
		STOP = make(chan struct{})
	default:
		// Channel is still open, no action needed
	}
}

// Initialize bandwidth limiter
func InitBandwidthLimiter() {
	if BandwidthLimit <= 0 {
		bandwidthLimiter = nil // No limit
		return
	}

	// Convert Mbps to bytes/second (1 kbps = 125 bytes/second), 1 token represents 1 byte
	bytesPerSecond := BandwidthLimit * 125
	bandwidthLimiter = rate.NewLimiter(rate.Limit(bytesPerSecond), bytesPerSecond*PacketBurst)
}
