// Package exporter
package exporter

import "github.com/ragnacron/msma/internal/model"

type Exporter interface {
	Export([]model.Metric) error
}
