// Package stdout
package stdout

import (
	"encoding/json"
	"os"

	"github.com/ragnacron/msma/internal/model"
)

type StdoutExporter struct{}

func New() *StdoutExporter {
	return &StdoutExporter{}
}

func (e *StdoutExporter) Export(metrics []model.Metric) error {
	data, err := json.Marshal(metrics)
	if err != nil {
		return err
	}

	if _, err := os.Stdout.Write(data); err != nil {
		return err
	}
	if _, err := os.Stdout.Write([]byte("\n")); err != nil {
		return err
	}

	return nil
}
