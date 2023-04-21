package config

import "time"

// Cancel return shutdown signal chan
func Cancel() <-chan struct{} { return ctx.Done() }

// Cancelled judge if project shutdown
func Cancelled() bool {
	select {
	case <-Cancel():
		return true
	default:
		return false
	}
}

// ShutdownTimeout return shutdown timeout duration
func ShutdownTimeout() time.Duration {
	return 3 * time.Second
}
