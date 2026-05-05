// Package collector
package collector

import (
	"os"
	"time"

	"github.com/ragnacron/msma/internal/model"
)

type Collector struct{} // Hostname?

func New() *Collector {
	return &Collector{}
}

func (c *Collector) Collect() (model.Metric, error) {
	host, err := os.Hostname()
	if err != nil {
		return model.Metric{}, err
	}

	payload, err := gatherMetricsPayload()
	if err != nil {
		return model.Metric{}, err
	}

	return model.Metric{
		Timestamp: time.Now().UTC(),
		Host:      host,
		Metrics:   payload,
	}, nil
}

func gatherMetricsPayload() (model.MetricsPayload, error) {
	cpu, err := getCPUMetrics()
	if err != nil {
		return model.MetricsPayload{}, err
	}

	memory, err := getMemoryMetric()
	if err != nil {
		return model.MetricsPayload{}, err
	}

	disk, err := getDiskMetrics()
	if err != nil {
		return model.MetricsPayload{}, err
	}

	sysInfo, err := getSystemInfoMetric()
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
