// Package logging
package logging

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type Level uint8

const (
	LevelError Level = iota
	LevelDebug
)

type Logger struct {
	*log.Logger
	level Level
}

func New(level Level, w io.Writer) *Logger {
	if w == nil {
		w = os.Stdout
	}

	return &Logger{
		Logger: log.New(w, "", log.LstdFlags),
		level:  level,
	}
}

func ParseLevel(s string) Level {
	switch strings.ToLower(s) {
	case "debug":
		return LevelDebug
	default:
		return LevelError
	}
}

func (l *Logger) Error(format string, v ...any) {
	if !l.enabled(LevelError) {
		return
	}

	l.printf("ERROR", format, v...)
}

func (l *Logger) Debug(format string, v ...any) {
	if !l.enabled(LevelDebug) {
		return
	}

	l.printf("DEBUG", format, v...)
}

func (l *Logger) enabled(level Level) bool {
	return l.level >= level
}

func (l *Logger) printf(prefix, format string, v ...any) {
	l.Printf("[%s] %s", prefix, fmt.Sprintf(format, v...))
}
