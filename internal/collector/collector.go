// Package collector
package collector

import (
	"os"
	"time"

	"github.com/ragnacron/msma/internal/model"
)

type Collector struct {
	Host string
} // Hostname?

func New() (*Collector, error) {
	host, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	return &Collector{
		Host: host,
	}, nil
}

func (c *Collector) Collect() (model.Metric, error) {
	payload, err := gatherMetricsPayload()
	if err != nil {
		return model.Metric{}, err
	}

	return model.Metric{
		Timestamp: time.Now().UTC(),
		Host:      c.Host,
		Metrics:   payload,
	}, nil
}

func gatherMetricsPayload() (model.MetricsPayload, error) {
	cpu, err := getCPUMetrics()
	if err != nil {
		return model.MetricsPayload{}, err
	}

	memory, err := getMemoryMetrics()
	if err != nil {
		return model.MetricsPayload{}, err
	}

	disk, err := getDiskMetrics()
	if err != nil {
		return model.MetricsPayload{}, err
	}

	sysInfo, err := getSystemInfoMetrics()
	if err != nil {
		return model.MetricsPayload{}, err
	}

	return model.MetricsPayload{
		CPU:    cpu,
		Memory: memory,
		Disk:   disk,
		System: sysInfo,
	}, nil
}
