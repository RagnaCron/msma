package config

import (
	"os"
	"testing"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Config
		wantErr bool
	}{
		{
			name: "valid config",
			cfg: &Config{
				IntervalSeconds: 5,
				QueueSize:       10,
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
				Exporter: ExporterConfig{
					Type: "kafka",
				},
			},
			wantErr: true,
		},
		{
			name: "http without endpoint",
			cfg: &Config{
				IntervalSeconds: 5,
				QueueSize:       10,
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
				Exporter: ExporterConfig{
					Type: "http",
					HTTP: HTTPConfig{
						Endpoint: "http://localhost:8080",
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoadEnvOverrides(t *testing.T) {
	// Ensure no config file exists
	os.Remove(defaultConfigPath)
	defer os.Remove(defaultConfigPath)

	t.Setenv("AGENT_INTERVAL", "10")
	t.Setenv("AGENT_QUEUE_SIZE", "20")
	t.Setenv("AGENT_EXPORTER_TYPE", "http")
	t.Setenv("AGENT_HTTP_ENDPOINT", "http://test")

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
	if cfg.Exporter.Type != "http" {
		t.Errorf("expected exporter type http, got %s", cfg.Exporter.Type)
	}
	if cfg.Exporter.HTTP.Endpoint != "http://test" {
		t.Errorf("expected endpoint http://test, got %s", cfg.Exporter.HTTP.Endpoint)
	}
}