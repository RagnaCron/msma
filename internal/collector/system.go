// Package collector
package collector

import (
	"github.com/ragnacron/msma/internal/model"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
)

func getCPUMetrics() ([]model.CPU, error) {
	return nil, nil
}

func getMemoryMetric() (model.Memory, error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		return model.Memory{}, err
	}

	return model.Memory{
		Total:       v.Total,
		Used:        v.Used,
		UsedPercent: v.UsedPercent,
		Free:        v.Free,
	}, nil
}

func getDiskMetrics() ([]model.Disk, error) {
	return nil, nil
}

func getSystemInfoMetric() (model.System, error) {
	uptime, err := host.Uptime()
	if err != nil {
		return model.System{}, err
	}

	kernel, err := host.KernelVersion()
	if err != nil {
		return model.System{}, err
	}

	return model.System{
		Uptime: uptime,
		Kernel: kernel,
	}, nil
}
