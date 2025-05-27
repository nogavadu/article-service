package env

import (
	"fmt"
	"github.com/nogavadu/articles-service/internal/config"
	"net"
	"os"
	"strconv"
	"time"
)

const (
	authServiceDomainEnv       = "AUTH_SERVICE_DOMAIN"
	authServicePortEnv         = "AUTH_SERVICE_PORT"
	authServiceTimeoutEnv      = "AUTH_SERVICE_TIMEOUT"
	authServiceRetriesCountEnv = "AUTH_SERVICE_RETRIES_COUNT"
	authServiceInsecureEnv     = "AUTH_SERVICE_INSECURE"
)

type authServiceConfig struct {
	domain   string
	port     string
	timeout  time.Duration
	retries  int
	insecure bool
}

func NewAuthServiceConfig() (config.AuthServiceConfig, error) {
	domain := os.Getenv(authServiceDomainEnv)
	if domain == "" {
		return nil, fmt.Errorf("environment variable %s is not set", authServiceDomainEnv)
	}

	port := os.Getenv(authServicePortEnv)
	if port == "" {
		return nil, fmt.Errorf("environment variable %s is not set", authServicePortEnv)
	}

	timeoutStr := os.Getenv(authServiceTimeoutEnv)
	if timeoutStr == "" {
		return nil, fmt.Errorf("environment variable %s is not set", authServiceTimeoutEnv)
	}
	timeout, err := time.ParseDuration(timeoutStr)

	retriesStr := os.Getenv(authServiceRetriesCountEnv)
	if retriesStr == "" {
		return nil, fmt.Errorf("environment variable %s is not set", authServiceRetriesCountEnv)
	}
	retries, err := strconv.Atoi(retriesStr)
	if err != nil {
		return nil, fmt.Errorf("invalid %s", authServiceRetriesCountEnv)
	}

	insecureStr := os.Getenv(authServiceInsecureEnv)
	if insecureStr == "" {
		return nil, fmt.Errorf("environment variable %s is not set", authServiceInsecureEnv)
	}
	insecure, err := strconv.ParseBool(insecureStr)
	if err != nil {
		return nil, fmt.Errorf("invalid %s", insecureStr)
	}

	return &authServiceConfig{
		domain:   domain,
		port:     port,
		timeout:  timeout,
		retries:  retries,
		insecure: insecure,
	}, nil
}

func (c *authServiceConfig) Address() string {
	return net.JoinHostPort(c.domain, c.port)
}

func (c *authServiceConfig) Timeout() time.Duration {
	return c.timeout
}

func (c *authServiceConfig) RetriesCount() int {
	return c.retries
}

func (c *authServiceConfig) Insecure() bool {
	return c.insecure
}
