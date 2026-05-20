package model

import (
	"encoding/json"
	"testing"
)

func TestMetricJSONSerialization(t *testing.T) {
	m := Metric{
		Host: "test",
		CPU:  []CPU{{Core: 0, Usage: 50.0}},
		Tags: map[string]string{"env": "test"},
	}

	data, err := json.Marshal(m)
	if err != nil {
		t.Fatalf("Marshal() failed: %v", err)
	}

	var metric Metric
	err = json.Unmarshal(data, &metric)
	if err != nil {
		t.Fatalf("Unmarshal() failed: %v", err)
	}

	if metric.Host != m.Host {
		t.Errorf("expected host %s, got %s", m.Host, metric.Host)
	}
	if len(metric.CPU) != 1 {
		t.Errorf("expected 1 CPU metric, got %d", len(metric.CPU))
	}
	if metric.Tags["env"] != "test" {
		t.Errorf("expected tag env=test, got %v", metric.Tags)
	}
}
