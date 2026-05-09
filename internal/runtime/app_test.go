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

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic for unsupported exporter type")
		}
	}()

	_ = New(cfg)
}

func TestNewStdout(t *testing.T) {
	cfg := &config.Config{
		IntervalSeconds: 1,
		QueueSize:       10,
		Exporter: config.ExporterConfig{
			Type: "stdout",
		},
	}
	app := New(cfg)
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
	app := New(cfg)
	if app == nil {
		t.Fatal("expected app to be created")
	}
}