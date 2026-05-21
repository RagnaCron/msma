// Package config
package config

type Config struct {
	IntervalSeconds int            `yaml:"interval"`
	QueueSize       int            `yaml:"queue_size"`
	LogLevel        string         `yaml:"log_level"`
	Exporter        ExporterConfig `yaml:"exporter"`
}

type ExporterConfig struct {
	Type string     `yaml:"type"`
	HTTP HTTPConfig `yaml:"http"`
}

type HTTPConfig struct {
	Endpoint string `yaml:"endpoint"`
	Timeout  int64  `yaml:"timeout"`
}
