// Package app
package app

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ragnacron/msma/internal/config"
	"github.com/ragnacron/msma/internal/exporter"
	"github.com/ragnacron/msma/internal/exporter/http"
	"github.com/ragnacron/msma/internal/exporter/stdout"
	"github.com/ragnacron/msma/internal/model"
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

func (a *App) Run() error {
	ch := a.Queue.Channel()
	var wgExporter sync.WaitGroup
	var wgCollector sync.WaitGroup
	stop := make(chan struct{})
	sigCh := make(chan os.Signal, 1)

	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	wgExporter.Add(1)
	wgExporter.Go(func() { // Exporter goroutine
		for metric := range ch {
			// todo: this will need some loggin at some point in time
			// error is currently droped
			_ = a.Exporter.Export([]model.Metric{metric})
		}
	})

	ticker := time.NewTicker(time.Duration(a.Config.IntervalSeconds) * time.Second)

	wgCollector.Add(1)
	wgCollector.Go(func() { // Collector goroutine
		defer ticker.Stop()
		for {
			metric := model.Metric{}
			select {
			case <-stop:
				return
			case <-ticker.C:
				select {
				case ch <- metric:
				default:
					// DROP (todo: loggin will be a thing at one point in time)
				}
			}
		}
	})

	<-sigCh // Block

	close(stop) // Stop Collector
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

	return nil
}
