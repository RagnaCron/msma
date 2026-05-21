package runtime

import (
	"testing"

	"github.com/ragnacron/msma/internal/config"
	"github.com/ragnacron/msma/internal/logging"
)

func TestNewUnsupportedExporter(t *testing.T) {
	cfg := &config.Config{
		IntervalSeconds: 1,
		QueueSize:       10,
		LogLevel:        "debug",
		Exporter: config.ExporterConfig{
			Type: "kafka",
		},
	}

	l := logging.New(cfg.LogLevel, nil)

	_, err := New(cfg, l)
	if err == nil {
		t.Errorf("expected error for unsupported exporter type")
	}
}

func TestNewStdout(t *testing.T) {
	cfg := &config.Config{
		IntervalSeconds: 1,
		QueueSize:       10,
		LogLevel:        "debug",
		Exporter: config.ExporterConfig{
			Type: "stdout",
		},
	}
	l := logging.New(cfg.LogLevel, nil)
	app, err := New(cfg, l)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if app == nil {
		t.Fatal("expected app to be created")
	}
}

func TestNewHTTP(t *testing.T) {
	cfg := &config.Config{
		IntervalSeconds: 1,
		QueueSize:       10,
		LogLevel:        "debug",
		Exporter: config.ExporterConfig{
			Type: "http",
			HTTP: config.HTTPConfig{
				Endpoint: "http://localhost:8080",
			},
		},
	}
	l := logging.New(cfg.LogLevel, nil)
	app, err := New(cfg, l)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if app == nil {
		t.Fatal("expected app to be created")
	}
}
