// Package collector
package collector

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ragnacron/msma/internal/model"
)

type Collector struct {
	Host string
}

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
	payload, err := collectMetrics()
	if err != nil {
		return model.Metric{}, err
	}

	return model.Metric{
		Timestamp: time.Now().UTC(),
		Host:      c.Host,
		Metrics:   *payload,
	}, nil
}

type collector struct {
	name string
	fn   func(*model.MetricsPayload) error
}

type collectorError struct {
	name string
	err  error
}

type MetricsError struct {
	errs []collectorError
}

func (e *MetricsError) Error() string {
	var parts []string

	for _, err := range e.errs {
		parts = append(parts, fmt.Sprintf("%s: %v", err.name, err.err))
	}

	return "metric collection failures: " + strings.Join(parts, ", ")
}

func collectMetrics() (*model.MetricsPayload, error) {
	var (
		wg     sync.WaitGroup
		mu     sync.Mutex
		metric model.MetricsPayload
		errs   []collectorError
	)

	collectors := getCollectors()

	for _, c := range collectors {
		wg.Go(func() {
			err := c.fn(&metric)

			mu.Lock()
			defer mu.Unlock()

			if err != nil {
				errs = append(errs, collectorError{
					name: c.name,
					err:  err,
				})
				return
			}
		})
	}

	wg.Wait()

	successes := len(collectors) - len(errs)

	if successes == 0 {
		return nil, &MetricsError{
			errs: errs,
		}
	}

	return &metric, nil
}

func getCollectors() []collector {
	return []collector{
		{
			name: "cpu",
			fn:   collectCPU,
		},
		{
			name: "memory",
			fn:   collectMem,
		},
		{
			name: "disk",
			fn:   collectDisk,
		},
		{
			name: "sysinfo",
			fn:   collectSysInfo,
		},
	}
}

func collectCPU(m *model.MetricsPayload) error {
	cpu, err := getCPUMetrics()
	if err != nil {
		return err
	}

	m.CPU = cpu

	return nil
}

func collectMem(m *model.MetricsPayload) error {
	memory, err := getMemoryMetrics()
	if err != nil {
		return err
	}

	m.Memory = &memory

	return nil
}

func collectDisk(m *model.MetricsPayload) error {
	disk, err := getDiskMetrics()
	if err != nil {
		return err
	}

	m.Disk = disk

	return nil
}

func collectSysInfo(m *model.MetricsPayload) error {
	sysInfo, err := getSystemInfoMetrics()
	if err != nil {
		return err
	}

	m.System = &sysInfo

	return nil
}
