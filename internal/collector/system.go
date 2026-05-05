// Package collector
package collector

import (
	"github.com/ragnacron/msma/internal/model"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
)

func getCPUMetrics() ([]model.CPU, error) {
	cpuData, err := cpu.Percent(0, true)
	if err != nil {
		return nil, err
	}

	cpus := make([]model.CPU, 0, len(cpuData))
	for i, v := range cpuData {
		cpus = append(cpus, model.CPU{
			Core:  i,
			Usage: v,
		})
	}

	return cpus, nil
}

func getMemoryMetrics() (model.Memory, error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		return model.Memory{}, err
	}

	return model.Memory{
		Total: v.Total,
		Used:  v.Used,
		Usage: v.UsedPercent,
		Free:  v.Free,
	}, nil
}

func getDiskMetrics() ([]model.Disk, error) {
	part, err := disk.Partitions(false)
	if err != nil {
		return nil, err
	}

	diskP := make([]model.Disk, 0, len(part))

	for _, stat := range part {
		diskStat, err := disk.Usage(stat.Mountpoint)
		if err != nil {
			continue
		}
		diskP = append(diskP, model.Disk{
			Mount: stat.Mountpoint,
			Usage: diskStat.UsedPercent,
		})
	}

	return diskP, nil
}

func getSystemInfoMetrics() (model.System, error) {
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
