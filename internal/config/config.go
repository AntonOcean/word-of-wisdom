package config

import "time"

type Config struct {
	Port                string
	MaxConnections      int
	ConnectionTimeout   time.Duration
	ShutdownTimeout     time.Duration
	RateLimitEvery100MS int
}
