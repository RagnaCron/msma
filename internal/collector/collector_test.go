package collector

import (
	"testing"
)

func TestCollect(t *testing.T) {
	col, err := New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	metric, err := col.Collect()
	if err != nil {
		t.Fatalf("Collect() failed: %v", err)
	}

	if metric.Host == "" {
		t.Error("expected host to be set")
	}
	if len(metric.CPU) == 0 {
		t.Error("expected CPU metrics to be collected")
	}
}
