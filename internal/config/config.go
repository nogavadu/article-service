package config

import "time"

type PGConfig interface {
	DSN() string
}

type HTTPServerConfig interface {
	Port() int
	Address() string
}

type AuthServiceConfig interface {
	Address() string
	Timeout() time.Duration
	RetriesCount() int
	Insecure() bool
}
