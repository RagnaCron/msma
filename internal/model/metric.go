// Package model
package model

import "time"

type Metric struct {
	Timestamp time.Time         `json:"timestamp"`
	Host      string            `json:"host"`
	Metrics   MetricsPayload    `json:"metrics"`
	Tags      map[string]string `json:"tags,omitempty"`
}

type MetricsPayload struct {
	CPU    []CPU  `json:"cpu"`
	Memory Memory `json:"memory"`
	Disk   []Disk `json:"disk"`
	System System `json:"system"`
}

type CPU struct {
	Core  int     `json:"core"`
	Usage float64 `json:"usage"`
}

type Memory struct {
	Total uint64  `json:"total"`
	Used  uint64  `json:"used"`
	Usage float64 `json:"usage"`
	Free  uint64  `json:"free"`
}

type Disk struct {
	Mount string  `json:"mount"`
	Usage float64 `json:"usage"`
}

type System struct {
	Uptime uint64 `json:"uptime"`
	Kernel string `json:"kernel"`
}
