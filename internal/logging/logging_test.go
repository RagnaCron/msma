package logging

import (
	"bytes"
	"strings"
	"testing"
)

func TestParseLevel(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want Level
	}{
		{
			name: "debug",
			in:   "debug",
			want: LevelDebug,
		},
		{
			name: "Debug uppercase",
			in:   "Debug",
			want: LevelDebug,
		},
		{
			name: "error",
			in:   "error",
			want: LevelError,
		},
		{
			name: "ERROR uppercase",
			in:   "ERROR",
			want: LevelError,
		},
		{
			name: "unknown defaults to error",
			in:   "trace",
			want: LevelError,
		},
		{
			name: "empty defaults to error",
			in:   "",
			want: LevelError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseLevel(tt.in)
			if got != tt.want {
				t.Errorf("parseLevel(%q) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

func TestLoggerErrorPrints(t *testing.T) {
	var buf bytes.Buffer
	logger := New("error", &buf)

	logger.Error("test %s", "message")

	got := buf.String()
	if !strings.Contains(got, "ERROR") {
		t.Errorf("expected ERROR prefix in output, got: %q", got)
	}
	if !strings.Contains(got, "test message") {
		t.Errorf("expected 'test message' in output, got: %q", got)
	}
}

func TestLoggerDebugSilent(t *testing.T) {
	var buf bytes.Buffer
	logger := New("error", &buf)

	logger.Debug("test message")

	if got := buf.String(); got != "" {
		t.Errorf("expected no output at error level, got: %q", got)
	}
}

func TestLoggerDebugVerbose(t *testing.T) {
	var buf bytes.Buffer
	logger := New("debug", &buf)

	logger.Debug("test message")

	got := buf.String()
	if !strings.Contains(got, "DEBUG") {
		t.Errorf("expected DEBUG prefix in output, got: %q", got)
	}
	if !strings.Contains(got, "test message") {
		t.Errorf("expected 'test message' in output, got: %q", got)
	}
}

func TestLoggerErrorSilent(t *testing.T) {
	var buf bytes.Buffer
	logger := New("debug", &buf)

	// Error should still print at debug level
	logger.Error("test error")
	if got := buf.String(); got == "" {
		t.Error("expected error output at debug level")
	}
}
