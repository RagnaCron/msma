package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ragnacron/msma/internal/model"
)

func TestExport(t *testing.T) {
	var receivedBody []model.Metric
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", ct)
		}
		_ = json.NewDecoder(r.Body).Decode(&receivedBody)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	exporter := NewHTTP(server.URL, 10)
	err := exporter.Export([]model.Metric{{Host: "test"}})
	if err != nil {
		t.Fatalf("Export() failed: %v", err)
	}

	if len(receivedBody) != 1 {
		t.Errorf("expected 1 metric, got %d", len(receivedBody))
	}
}

func TestExportError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	exporter := NewHTTP(server.URL, 10)
	err := exporter.Export([]model.Metric{{Host: "test"}})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestExportInvalidURL(t *testing.T) {
	exporter := NewHTTP("://invalid-url", 10)
	err := exporter.Export([]model.Metric{{Host: "test"}})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
