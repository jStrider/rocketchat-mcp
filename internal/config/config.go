package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const (
	defaultHTTPAddr    = ":8080"
	defaultLogLevel    = "info"
	defaultMaxBodySize = 10 * 1024 * 1024 // 10 MB
)

// Config holds all configuration for the MCP server.
// Immutable after creation via LoadFromEnv.
type Config struct {
	URL         string
	UserID      string
	AuthToken   string
	ReadOnly    bool
	AllowHTTP   bool
	LogLevel    string
	HTTPAddr    string
	MaxBodySize int64
}

// LoadFromEnv reads configuration from environment variables and validates it.
func LoadFromEnv() (Config, error) {
	cfg := Config{
		URL:         os.Getenv("ROCKETCHAT_URL"),
		UserID:      os.Getenv("ROCKETCHAT_USER_ID"),
		AuthToken:   os.Getenv("ROCKETCHAT_AUTH_TOKEN"),
		ReadOnly:    envBool("ROCKETCHAT_READ_ONLY"),
		AllowHTTP:   envBool("ROCKETCHAT_ALLOW_HTTP"),
		LogLevel:    envOr("ROCKETCHAT_LOG_LEVEL", defaultLogLevel),
		HTTPAddr:    envOr("ROCKETCHAT_MCP_ADDR", defaultHTTPAddr),
		MaxBodySize: envInt64("ROCKETCHAT_MAX_BODY_SIZE", defaultMaxBodySize),
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

// Validate checks that all required fields are set and valid.
func (c Config) Validate() error {
	if c.URL == "" {
		return fmt.Errorf("ROCKETCHAT_URL is required")
	}
	if c.UserID == "" {
		return fmt.Errorf("ROCKETCHAT_USER_ID is required")
	}
	if c.AuthToken == "" {
		return fmt.Errorf("ROCKETCHAT_AUTH_TOKEN is required")
	}

	u, err := url.Parse(c.URL)
	if err != nil {
		return fmt.Errorf("ROCKETCHAT_URL is not a valid URL: %w", err)
	}
	if u.Scheme != "https" && !c.AllowHTTP {
		return fmt.Errorf("ROCKETCHAT_URL must use HTTPS (set ROCKETCHAT_ALLOW_HTTP=true for dev)")
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("ROCKETCHAT_URL scheme must be http or https, got %q", u.Scheme)
	}
	if u.Host == "" {
		return fmt.Errorf("ROCKETCHAT_URL must have a host")
	}

	validLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
	if !validLevels[c.LogLevel] {
		return fmt.Errorf("ROCKETCHAT_LOG_LEVEL must be one of debug, info, warn, error; got %q", c.LogLevel)
	}

	if c.MaxBodySize <= 0 {
		return fmt.Errorf("ROCKETCHAT_MAX_BODY_SIZE must be positive")
	}

	return nil
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func envBool(key string) bool {
	return strings.EqualFold(os.Getenv(key), "true")
}

func envInt64(key string, def int64) int64 {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	n, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return def
	}
	return n
}
