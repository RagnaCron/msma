package stdout

import (
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/ragnacron/msma/internal/model"
)

func TestExport(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	e := New()
	err := e.Export([]model.Metric{{Host: "test1"}, {Host: "test2"}})
	if err != nil {
		t.Fatalf("Export() failed: %v", err)
	}

	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = oldStdout

	// Verify output is valid JSON array
	var metrics []model.Metric
	err = json.Unmarshal(out, &metrics)
	if err != nil {
		t.Fatalf("Output is not valid JSON: %v", err)
	}

	if len(metrics) != 2 {
		t.Errorf("expected 2 metrics, got %d", len(metrics))
	}
}