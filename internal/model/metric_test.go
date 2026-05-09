package model

import (
	"encoding/json"
	"testing"
)

func TestMetricJSONSerialization(t *testing.T) {
	m := Metric{
		Host: "test",
		Metrics: MetricsPayload{
			CPU: []CPU{{Core: 0, Usage: 50.0}},
		},
		Tags: map[string]string{"env": "test"},
	}

	data, err := json.Marshal(m)
	if err != nil {
		t.Fatalf("Marshal() failed: %v", err)
	}

	var decoded Metric
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Unmarshal() failed: %v", err)
	}

	if decoded.Host != m.Host {
		t.Errorf("expected host %s, got %s", m.Host, decoded.Host)
	}
	if len(decoded.Metrics.CPU) != 1 {
		t.Errorf("expected 1 CPU metric, got %d", len(decoded.Metrics.CPU))
	}
	if decoded.Tags["env"] != "test" {
		t.Errorf("expected tag env=test, got %v", decoded.Tags)
	}
}