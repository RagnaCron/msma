// Package config
package config

import (
	"errors"
	"os"
	"strconv"

	"go.yaml.in/yaml/v4"
)

const defaultConfigPath = "./.config/msma/config.yaml"

func Load() (*Config, error) {
	cfg := defaultConfig()

	if err := loadFromFile(cfg); err != nil {
		return nil, err
	}

	if err := applyEnvOverrides(cfg); err != nil {
		return nil, err
	}

	if err := Validate(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func defaultConfig() *Config {
	return &Config{
		IntervalSeconds: 5,
		QueueSize:       10,
		Exporter: ExporterConfig{
			Type: "stdout",
		},
	}
}

func loadFromFile(cfg *Config) error {
	data, err := os.ReadFile(defaultConfigPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	return yaml.Unmarshal(data, cfg)
}

func applyEnvOverrides(cfg *Config) error {
	if interval := os.Getenv("AGENT_INTERVAL"); interval != "" {
		integer, err := strconv.Atoi(interval)
		if err != nil {
			return err
		}
		cfg.IntervalSeconds = integer
	}

	if size := os.Getenv("AGENT_QUEUE_SIZE"); size != "" {
		integer, err := strconv.Atoi(size)
		if err != nil {
			return err
		}
		cfg.QueueSize = integer
	}

	if eType := os.Getenv("AGENT_EXPORTER_TYPE"); eType != "" {
		cfg.Exporter.Type = eType
	}

	if endpoint := os.Getenv("AGENT_HTTP_ENDPOINT"); endpoint != "" {
		cfg.Exporter.HTTP.Endpoint = endpoint
	}

	if timeout := os.Getenv("AGENT_HTTP_TIMEOUT"); timeout != "" {
		integer, err := strconv.Atoi(timeout)
		if err != nil {
			return err
		}
		cfg.Exporter.HTTP.Timeout = int64(integer)
	}

	return nil
}
