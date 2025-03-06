package config

import "time"

type Config struct {
	Port            string
	MaxConnections  int
	ShutdownTimeout time.Duration
}
