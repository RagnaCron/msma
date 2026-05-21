// Package runtime
package runtime

import (
	"errors"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ragnacron/msma/internal/collector"
	"github.com/ragnacron/msma/internal/config"
	"github.com/ragnacron/msma/internal/exporter"
	"github.com/ragnacron/msma/internal/exporter/http"
	"github.com/ragnacron/msma/internal/exporter/stdout"
	"github.com/ragnacron/msma/internal/logging"
	"github.com/ragnacron/msma/internal/model"
	"github.com/ragnacron/msma/internal/queue"
)

type App struct {
	Config   *config.Config
	Queue    *queue.Queue
	Logger   *logging.Logger
	Exporter exporter.Exporter
}

func New(cfg *config.Config, l *logging.Logger) (*App, error) {
	var e exporter.Exporter

	switch cfg.Exporter.Type {
	case "stdout":
		e = stdout.New()
	case "http":
		e = http.NewHTTP(cfg.Exporter.HTTP.Endpoint, cfg.Exporter.HTTP.Timeout, l)
	default:
		return nil, errors.New("unsupported exporter type")
	}

	return &App{
		Config:   cfg,
		Queue:    queue.New(cfg.QueueSize),
		Logger:   l,
		Exporter: e,
	}, nil
}

func (a *App) Run() error {
	col, err := collector.New()
	if err != nil {
		return err
	}
	ch := a.Queue.Channel()
	var wgExporter sync.WaitGroup
	var wgCollector sync.WaitGroup
	stop := make(chan struct{})
	sigCh := make(chan os.Signal, 1)

	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	wgExporter.Go(func() { // Exporter goroutine
		for metric := range ch {
			if err := a.Exporter.Export([]model.Metric{metric}); err != nil {
				a.Logger.Error("export failed: %v\n", err)
			}
		}
	})

	ticker := time.NewTicker(time.Duration(a.Config.IntervalSeconds) * time.Second)
	wgCollector.Go(func() { // Collector goroutine
		defer ticker.Stop()
		for {
			select {
			case <-stop:
				return
			case <-ticker.C:
				metrics, err := col.Collect()
				if err != nil {
					continue
				}
				select {
				case ch <- *metrics:
				default:
					a.Logger.Debug("metric dropped: queue full")
				}
			}
		}
	})

	<-sigCh // Block
	a.Logger.Debug("shutdown signal received")

	close(stop) // Stop Collector
	a.Logger.Debug("draining queue...")
	wgCollector.Wait()

	close(ch) // Stop Channel
	done := make(chan struct{})
	go func() {
		wgExporter.Wait()
		close(done)
	}()
	select {
	case <-done:
		// clean shutdown
	case <-time.After(5 * time.Second):
		// timeout -> exit anyway
	}
	a.Logger.Debug("shutdown complete")

	return nil
}
