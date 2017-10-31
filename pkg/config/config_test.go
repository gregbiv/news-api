package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	setGlobalConfigEnv()

	cfg, err := LoadEnv()
	assert.NoError(t, err)

	assertConfig(t, cfg)
}

func assertConfig(t *testing.T, cfg *Specification) {
	assert.Equal(t, 8090, cfg.Port)
	assert.Equal(t, "info", cfg.LogLevel)
}

func setGlobalConfigEnv() {
	os.Setenv("PORT", "8090")
	os.Setenv("LOG_LEVEL", "info")
}
