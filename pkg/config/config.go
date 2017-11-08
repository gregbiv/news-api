package config

import (
	"github.com/kelseyhightower/envconfig"
)

// Specification represents structured configuration variables
type Specification struct {
	Debug     bool   `envconfig:"DEBUG" default:"false"`
	LogLevel  string `envconfig:"LOG_LEVEL" default:"info"`
	Port      int    `envconfig:"PORT" default:"8090"`
	Migration Migration
	Database  struct {
		PostgresDB struct {
			DSN string `envconfig:"DATABASE_DSN"`
		}
	}
}

// Migration config
type Migration struct {
	Version uint   `envconfig:"DATABASE_VERSION"`
	Dir     string `envconfig:"MIGRATION_DIR"`
}

// LoadEnv loads config variables into Specification
func LoadEnv() (*Specification, error) {
	var conf Specification
	err := envconfig.Process("", &conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}
