package env

import (
	"fmt"
	"github.com/nogavadu/articles-service/internal/config"
	"os"
)

const (
	pgHostEnv     = "PG_HOST"
	pgPortEnv     = "PG_PORT"
	pgDbNameEnv   = "PG_DB_NAME"
	pgUserEnv     = "PG_USER"
	pgPasswordEnv = "PG_PASSWORD"
)

type pgConfig struct {
	host     string
	port     string
	dbName   string
	user     string
	password string
}

func NewPGConfig() (config.PGConfig, error) {
	const op = "config.NewPGConfig"

	host := os.Getenv(pgHostEnv)
	if host == "" {
		return nil, fmt.Errorf("%s: %s: failed to get env variable", op, pgHostEnv)
	}

	port := os.Getenv(pgPortEnv)
	if port == "" {
		return nil, fmt.Errorf("%s: %s: failed to get env variable", op, pgPortEnv)
	}

	name := os.Getenv(pgDbNameEnv)
	if name == "" {
		return nil, fmt.Errorf("%s: %s: failed to get env variable", op, pgDbNameEnv)
	}

	user := os.Getenv(pgUserEnv)
	if user == "" {
		return nil, fmt.Errorf("%s: %s: failed to get env variable", op, pgUserEnv)
	}

	password := os.Getenv(pgPasswordEnv)
	if password == "" {
		return nil, fmt.Errorf("%s: %s: failed to get env variable", op, pgPasswordEnv)
	}

	return &pgConfig{
		host:     host,
		port:     port,
		dbName:   name,
		user:     user,
		password: password,
	}, nil
}

func (c *pgConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		c.host, c.port, c.dbName, c.user, c.password,
	)
}
