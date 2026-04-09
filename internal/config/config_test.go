package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setRequiredEnv(t *testing.T) {
	t.Helper()
	t.Setenv("ROCKETCHAT_URL", "https://chat.example.com")
	t.Setenv("ROCKETCHAT_USER_ID", "uid123")
	t.Setenv("ROCKETCHAT_AUTH_TOKEN", "tok456")
}

func TestLoadFromEnv_ValidConfig(t *testing.T) {
	setRequiredEnv(t)

	cfg, err := LoadFromEnv()
	require.NoError(t, err)
	assert.Equal(t, "https://chat.example.com", cfg.URL)
	assert.Equal(t, "uid123", cfg.UserID)
	assert.Equal(t, "tok456", cfg.AuthToken)
	assert.False(t, cfg.ReadOnly)
	assert.False(t, cfg.AllowHTTP)
	assert.Equal(t, "info", cfg.LogLevel)
	assert.Equal(t, ":8080", cfg.HTTPAddr)
	assert.Equal(t, int64(10*1024*1024), cfg.MaxBodySize)
}

func TestLoadFromEnv_MissingURL(t *testing.T) {
	t.Setenv("ROCKETCHAT_USER_ID", "uid123")
	t.Setenv("ROCKETCHAT_AUTH_TOKEN", "tok456")

	_, err := LoadFromEnv()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "ROCKETCHAT_URL is required")
}

func TestLoadFromEnv_MissingUserID(t *testing.T) {
	t.Setenv("ROCKETCHAT_URL", "https://chat.example.com")
	t.Setenv("ROCKETCHAT_AUTH_TOKEN", "tok456")

	_, err := LoadFromEnv()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "ROCKETCHAT_USER_ID is required")
}

func TestLoadFromEnv_MissingAuthToken(t *testing.T) {
	t.Setenv("ROCKETCHAT_URL", "https://chat.example.com")
	t.Setenv("ROCKETCHAT_USER_ID", "uid123")

	_, err := LoadFromEnv()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "ROCKETCHAT_AUTH_TOKEN is required")
}

func TestLoadFromEnv_HTTPNotAllowed(t *testing.T) {
	t.Setenv("ROCKETCHAT_URL", "http://chat.example.com")
	t.Setenv("ROCKETCHAT_USER_ID", "uid123")
	t.Setenv("ROCKETCHAT_AUTH_TOKEN", "tok456")

	_, err := LoadFromEnv()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "must use HTTPS")
}

func TestLoadFromEnv_HTTPAllowedWithFlag(t *testing.T) {
	t.Setenv("ROCKETCHAT_URL", "http://chat.example.com")
	t.Setenv("ROCKETCHAT_USER_ID", "uid123")
	t.Setenv("ROCKETCHAT_AUTH_TOKEN", "tok456")
	t.Setenv("ROCKETCHAT_ALLOW_HTTP", "true")

	cfg, err := LoadFromEnv()
	require.NoError(t, err)
	assert.True(t, cfg.AllowHTTP)
}

func TestLoadFromEnv_InvalidScheme(t *testing.T) {
	t.Setenv("ROCKETCHAT_URL", "ftp://chat.example.com")
	t.Setenv("ROCKETCHAT_USER_ID", "uid123")
	t.Setenv("ROCKETCHAT_AUTH_TOKEN", "tok456")
	t.Setenv("ROCKETCHAT_ALLOW_HTTP", "true")

	_, err := LoadFromEnv()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "scheme must be http or https")
}

func TestLoadFromEnv_ReadOnly(t *testing.T) {
	setRequiredEnv(t)
	t.Setenv("ROCKETCHAT_READ_ONLY", "true")

	cfg, err := LoadFromEnv()
	require.NoError(t, err)
	assert.True(t, cfg.ReadOnly)
}

func TestLoadFromEnv_CustomValues(t *testing.T) {
	setRequiredEnv(t)
	t.Setenv("ROCKETCHAT_LOG_LEVEL", "debug")
	t.Setenv("ROCKETCHAT_MCP_ADDR", ":9090")
	t.Setenv("ROCKETCHAT_MAX_BODY_SIZE", "5242880")

	cfg, err := LoadFromEnv()
	require.NoError(t, err)
	assert.Equal(t, "debug", cfg.LogLevel)
	assert.Equal(t, ":9090", cfg.HTTPAddr)
	assert.Equal(t, int64(5242880), cfg.MaxBodySize)
}

func TestLoadFromEnv_InvalidLogLevel(t *testing.T) {
	setRequiredEnv(t)
	t.Setenv("ROCKETCHAT_LOG_LEVEL", "verbose")

	_, err := LoadFromEnv()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "LOG_LEVEL must be one of")
}

func TestLoadFromEnv_InvalidMaxBodySize(t *testing.T) {
	setRequiredEnv(t)
	t.Setenv("ROCKETCHAT_MAX_BODY_SIZE", "not-a-number")

	cfg, err := LoadFromEnv()
	require.NoError(t, err)
	// Falls back to default
	assert.Equal(t, int64(10*1024*1024), cfg.MaxBodySize)
}

func TestValidate_NoHost(t *testing.T) {
	cfg := Config{
		URL:       "https://",
		UserID:    "uid",
		AuthToken: "tok",
		LogLevel:  "info",
		MaxBodySize: defaultMaxBodySize,
	}
	err := cfg.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "must have a host")
}
