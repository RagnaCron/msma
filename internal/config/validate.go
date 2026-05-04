// Package config
package config

import (
	"errors"
)

func Validate(cfg *Config) error {
	if cfg.QueueSize < 1 {
		return errors.New("invalid queue_size: must be > 0")
	}

	if cfg.IntervalSeconds < 1 {
		return errors.New("invalid interval: must be > 0")
	}

	if cfg.Exporter.Type != "stdout" && cfg.Exporter.Type != "http" {
		return errors.New("invalid type: must be stdout or http")
	}

	if cfg.Exporter.Type == "http" && cfg.Exporter.HTTP.Endpoint == "" {
		return errors.New("invalid endpoint: must not be empty")
	}

	return nil
}
