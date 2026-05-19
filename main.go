// Package main
package main

import (
	"fmt"
	"log"

	"github.com/ragnacron/msma/internal/config"
	"github.com/ragnacron/msma/internal/runtime"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalln(fmt.Errorf("error loading config: %w", err))
	}

	app, err := runtime.New(cfg)
	if err != nil {
		log.Fatalln(fmt.Errorf("error starting runtime: %w", err))
	}
	err = app.Run()
	if err != nil {
		log.Fatalln(fmt.Errorf("error running app: %w", err))
	}
}
