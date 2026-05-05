// Package app
package app

import (
	"github.com/ragnacron/msma/internal/config"
	"github.com/ragnacron/msma/internal/exporter"
	"github.com/ragnacron/msma/internal/exporter/http"
	"github.com/ragnacron/msma/internal/exporter/stdout"
	"github.com/ragnacron/msma/internal/queue"
)

type App struct {
	Config   *config.Config
	Queue    *queue.Queue
	Exporter exporter.Exporter
}

func New(cfg *config.Config) *App {
	var e exporter.Exporter

	switch cfg.Exporter.Type {
	case "stdout":
		e = stdout.New()
	case "http":
		e = http.NewHTTP(cfg.Exporter.HTTP.Endpoint)
	default:
		panic("unsupported exporter type") // todo: return error in New -> App
	}

	return &App{
		Config:   cfg,
		Queue:    queue.New(cfg.QueueSize),
		Exporter: e,
	}
}

func Run() error {
	return nil
}
