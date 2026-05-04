// Package exporter
package exporter

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/ragnacron/msma/internal/model"
)

type StdoutExporter struct{}

func New() *StdoutExporter {
	return &StdoutExporter{}
}

func (*StdoutExporter) Export(metrics []model.Metric) error {
	data, err := json.Marshal(metrics)
	if err != nil {
		return err
	}

	buffer := bufio.NewWriter(os.Stdout)
	if _, err := buffer.Write(data); err != nil {
		return err
	}
	if err := buffer.WriteByte('\n'); err != nil {
		return err
	}
	if err := buffer.Flush(); err != nil {
		return err
	}

	return nil
}
