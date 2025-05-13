package config

type PGConfig interface {
	DSN() string
}

type HTTPServerConfig interface {
	Port() int
	Address() string
}
