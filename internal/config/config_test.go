package config

import (
	"os"
	"testing"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name        string
		cfg         *Config
		wantErr     bool
		wantTimeout int64
	}{
		{
			name: "valid config",
			cfg: &Config{
				IntervalSeconds: 5,
				QueueSize:       10,
				LogLevel:        "error",
				Exporter: ExporterConfig{
					Type: "stdout",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid queue size",
			cfg: &Config{
				IntervalSeconds: 5,
				QueueSize:       0,
				LogLevel:        "error",
				Exporter: ExporterConfig{
					Type: "stdout",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid interval",
			cfg: &Config{
				IntervalSeconds: 0,
				QueueSize:       10,
				LogLevel:        "error",
				Exporter: ExporterConfig{
					Type: "stdout",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid exporter type",
			cfg: &Config{
				IntervalSeconds: 5,
				QueueSize:       10,
				LogLevel:        "error",
				Exporter: ExporterConfig{
					Type: "kafka",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid log level",
			cfg: &Config{
				IntervalSeconds: 5,
				QueueSize:       10,
				LogLevel:        "info",
				Exporter: ExporterConfig{
					Type: "stdout",
				},
			},
			wantErr: true,
		},
		{
			name: "http without endpoint",
			cfg: &Config{
				IntervalSeconds: 5,
				QueueSize:       10,
				LogLevel:        "error",
				Exporter: ExporterConfig{
					Type: "http",
					HTTP: HTTPConfig{},
				},
			},
			wantErr: true,
		},
		{
			name: "http with endpoint",
			cfg: &Config{
				IntervalSeconds: 5,
				QueueSize:       10,
				LogLevel:        "error",
				Exporter: ExporterConfig{
					Type: "http",
					HTTP: HTTPConfig{
						Endpoint: "http://localhost:8080",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "http negative timeout clamped",
			cfg: &Config{
				IntervalSeconds: 5,
				QueueSize:       10,
				LogLevel:        "error",
				Exporter: ExporterConfig{
					Type: "http",
					HTTP: HTTPConfig{
						Endpoint: "http://localhost:8080",
						Timeout:  -10,
					},
				},
			},
			wantErr:     false,
			wantTimeout: defaultHTTPTimeout,
		},
		{
			name: "http zero timeout clamped",
			cfg: &Config{
				IntervalSeconds: 5,
				QueueSize:       10,
				LogLevel:        "error",
				Exporter: ExporterConfig{
					Type: "http",
					HTTP: HTTPConfig{
						Endpoint: "http://localhost:8080",
						Timeout:  0,
					},
				},
			},
			wantErr:     false,
			wantTimeout: defaultHTTPTimeout,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantTimeout != 0 && tt.cfg.Exporter.HTTP.Timeout != tt.wantTimeout {
				t.Errorf("expected timeout %d, got %d", tt.wantTimeout, tt.cfg.Exporter.HTTP.Timeout)
			}
		})
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := defaultConfig()

	if cfg.IntervalSeconds != 5 {
		t.Errorf("expected interval 5, got %d", cfg.IntervalSeconds)
	}
	if cfg.QueueSize != 10 {
		t.Errorf("expected queue size 10, got %d", cfg.QueueSize)
	}
	if cfg.LogLevel != "error" {
		t.Errorf("expected log level to be 'error', got %s", cfg.LogLevel)
	}
	if cfg.Exporter.Type != "stdout" {
		t.Errorf("expected exporter type stdout, got %s", cfg.Exporter.Type)
	}
	if cfg.Exporter.HTTP.Timeout != defaultHTTPTimeout {
		t.Errorf("expected timeout %d, got %d", defaultHTTPTimeout, cfg.Exporter.HTTP.Timeout)
	}
}

func TestLoadEnvOverrides(t *testing.T) {
	// Ensure no config file exists
	_ = os.Remove(defaultConfigPath)
	defer func() {
		_ = os.Remove(defaultConfigPath)
	}()
	t.Setenv("AGENT_INTERVAL", "10")
	t.Setenv("AGENT_QUEUE_SIZE", "20")
	t.Setenv("AGENT_EXPORTER_TYPE", "http")
	t.Setenv("AGENT_HTTP_ENDPOINT", "http://test")
	t.Setenv("AGENT_HTTP_TIMEOUT", "60")
	t.Setenv("AGENT_LOG_LEVEL", "debug")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if cfg.IntervalSeconds != 10 {
		t.Errorf("expected interval 10, got %d", cfg.IntervalSeconds)
	}
	if cfg.QueueSize != 20 {
		t.Errorf("expected queue size 20, got %d", cfg.QueueSize)
	}
	if cfg.LogLevel != "debug" {
		t.Errorf("expected log level to be 'debug', got %s", cfg.LogLevel)
	}
	if cfg.Exporter.Type != "http" {
		t.Errorf("expected exporter type http, got %s", cfg.Exporter.Type)
	}
	if cfg.Exporter.HTTP.Endpoint != "http://test" {
		t.Errorf("expected endpoint http://test, got %s", cfg.Exporter.HTTP.Endpoint)
	}
	if cfg.Exporter.HTTP.Timeout != 60 {
		t.Errorf("expected timeout 60, got %d", cfg.Exporter.HTTP.Timeout)
	}
}
