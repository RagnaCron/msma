package runtime

import (
	"testing"

	"github.com/ragnacron/msma/internal/config"
)

func TestNewUnsupportedExporter(t *testing.T) {
	cfg := &config.Config{
		IntervalSeconds: 1,
		QueueSize:       10,
		Exporter: config.ExporterConfig{
			Type: "kafka",
		},
	}

	_, err := New(cfg)
	if err == nil {
		t.Errorf("expected error for unsupported exporter type")
	}
}

func TestNewStdout(t *testing.T) {
	cfg := &config.Config{
		IntervalSeconds: 1,
		QueueSize:       10,
		Exporter: config.ExporterConfig{
			Type: "stdout",
		},
	}
	app, err := New(cfg)
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
		Exporter: config.ExporterConfig{
			Type: "http",
			HTTP: config.HTTPConfig{
				Endpoint: "http://localhost:8080",
			},
		},
	}
	app, err := New(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if app == nil {
		t.Fatal("expected app to be created")
	}
}
