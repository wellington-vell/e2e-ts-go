package lib

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
)

type NodeEnv string

const (
	NodeEnvDevelopment NodeEnv = "development"
	NodeEnvProduction  NodeEnv = "production"
	NodeEnvTest        NodeEnv = "test"
)

type EnvConfig struct {
	NodeEnv     NodeEnv
	ServerPort  int
	CorsOrigin  string
	DatabaseURL string
	ServerURL   string
	AuthSecret  string
}

var Env *EnvConfig

func LoadEnv() {
	cfg, err := parseEnv()
	if err != nil {
		panic(fmt.Sprintf("failed to load env: %v", err))
	}
	Env = cfg
}

func parseEnv() (*EnvConfig, error) {
	serverPortRaw, ok := os.LookupEnv("SERVER_PORT")
	if !ok {
		return nil, fmt.Errorf("SERVER_PORT is required")
	}
	serverPort, err := strconv.Atoi(serverPortRaw)
	if err != nil {
		return nil, fmt.Errorf("SERVER_PORT must be a valid integer: %w", err)
	}

	corsOrigin, ok := os.LookupEnv("CORS_ORIGIN")
	if !ok {
		return nil, fmt.Errorf("CORS_ORIGIN is required")
	}
	if err := validateOrigin("CORS_ORIGIN", corsOrigin); err != nil {
		return nil, err
	}

	databaseURL, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}
	if err := validateURL("DATABASE_URL", databaseURL); err != nil {
		return nil, err
	}

	serverURL, ok := os.LookupEnv("SERVER_URL")
	if !ok {
		return nil, fmt.Errorf("SERVER_URL is required")
	}
	if err := validateURL("SERVER_URL", serverURL); err != nil {
		return nil, err
	}

	authSecret, ok := os.LookupEnv("AUTH_SECRET")
	if !ok {
		return nil, fmt.Errorf("AUTH_SECRET is required")
	}
	if err := validateHex("AUTH_SECRET", authSecret); err != nil {
		return nil, err
	}

	nodeEnvRaw, ok := os.LookupEnv("NODE_ENV")
	if !ok {
		nodeEnvRaw = "development"
	}
	var nodeEnv NodeEnv
	switch nodeEnvRaw {
	case "development":
		nodeEnv = NodeEnvDevelopment
	case "production":
		nodeEnv = NodeEnvProduction
	case "test":
		nodeEnv = NodeEnvTest
	default:
		return nil, fmt.Errorf("NODE_ENV must be development, production, or test")
	}

	return &EnvConfig{
		NodeEnv:     nodeEnv,
		ServerPort:  serverPort,
		CorsOrigin:  corsOrigin,
		DatabaseURL: databaseURL,
		ServerURL:   serverURL,
		AuthSecret:  authSecret,
	}, nil
}

func validateURL(name, raw string) error {
	u, err := url.Parse(raw)
	if err != nil {
		return fmt.Errorf("%s must be a valid URL: %w", name, err)
	}
	if u.Scheme == "" {
		return fmt.Errorf("%s must have a scheme (e.g., http:// or postgres://)", name)
	}
	if u.Host == "" {
		return fmt.Errorf("%s must have a host", name)
	}
	return nil
}

func validateHex(name, raw string) error {
	if len(raw) != 64 {
		return fmt.Errorf("%s must be exactly 64 hex characters", name)
	}
	for i := 0; i < len(raw); i++ {
		c := raw[i]
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return fmt.Errorf("%s must be a valid hex string", name)
		}
	}
	return nil
}

func validateOrigin(name, raw string) error {
	if raw == "*" {
		return fmt.Errorf("%s must not be a wildcard (*)", name)
	}
	u, err := url.Parse(raw)
	if err != nil {
		return fmt.Errorf("%s must be a valid URL: %w", name, err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("%s must use http or https scheme", name)
	}
	if u.Host == "" {
		return fmt.Errorf("%s must have a host", name)
	}
	return nil
}
