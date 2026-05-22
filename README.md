# Minimal System Metrics Agent

A lightweight Go daemon that collects host-level metrics (CPU, memory, disk,
uptime, kernel) and exports them via stdout or HTTP.

## Motivation

System metrics agents tend to be either heavy (Prometheus node_exporter with
dozens of collectors) or nonexistent. This project sits in the gap: a single
binary that collects CPU, memory, disk, and kernel info, exports them as JSON
via stdout or HTTP, and does it with no framework, no external dependencies
beyond gopsutil, and under 1000 lines of code. It's designed for environments
where you need to ship metrics to a custom endpoint without running a full
monitoring stack.

## Quick Start

```bash
# Install
go install github.com/RagnaCron/msma@latest
```

```bash
# Run with defaults
msma
```

```bash
# Run with HTTP exporter
cd ~/.config/msma
cat << EOF > config.yaml
exporter:
  type: http
  http:
    endpoint: http://localhost:8081/metrics
    timeout: 30
EOF
./msma

```

## Features

- Concurrent collection using four parallel subsystem readers
- Partial failure tolerance — a failing subsystem does not block others
- Bounded channel with non-blocking drop on overflow
- Two exporters: stdout and HTTP (POST, JSON)
- Configurable log level (error, debug)
- Configurable HTTP client timeout
- Graceful shutdown on SIGINT/SIGTERM

## Build

```bash
make build   # compiles to ./msma
make run     # build and run
make clean   # remove binary
```

## Configuration and Usage

Config is loaded from three sources in order:
defaults → YAML file → environment variables.
The YAML file path is `~/.config/msma/config.yaml`.
Missing config file is silently ignored.

### YAML Config

```yaml
interval: 5                    # collection interval in seconds
queue_size: 10                 # number of metrics buffered
log_level: error               # error, debug
exporter:
  type: stdout                 # stdout, http
  http:
    endpoint: ""               # required when type is http
    timeout: 30                # HTTP client timeout in seconds
```

### Environment Variables

| Variable                | Default | Description                              |
| ------------------------- | --------- | --------------------------------------- |
| `AGENT_INTERVAL`        | `5`     | Collection interval in seconds           |
| `AGENT_QUEUE_SIZE`      | `10`    | Metrics queue capacity                   |
| `AGENT_EXPORTER_TYPE`   | `stdout`| `stdout` or `http`                     |
| `AGENT_HTTP_ENDPOINT`   | `""`    | HTTP POST target (required when `http`)  |
| `AGENT_HTTP_TIMEOUT`    | `30`    | HTTP client timeout in seconds          |
| `AGENT_LOG_LEVEL`       | `error` | Log verbosity: `error` or `debug`       |

Environment variables override both defaults and the YAML file.

### Default Values

| Field               | Default  |
|---------------------|----------|
| `interval`          | `5`      |
| `queue_size`        | `10`     |
| `log_level`         | `error`  |
| `exporter.type`     | `stdout` |
| `exporter.http.timeout` | `30` |

## Collector

The collector runs four subsystem readers concurrently (CPU, memory, disk,
system info) and merges the results. If a subsystem fails, the others are still
collected — the returned metric contains whatever succeeded and the error
lists the failed subsystems.

- **Disk metrics**: silently skip unaccessible mount points
- **CPU metrics**: per-core usage percentage
- **Memory metrics**: total, used, free, used percentage
- **System metrics**: uptime in seconds, kernel version string

## Exporters

### Stdout

Writes one JSON array per collection to stdout. Each line contains a single `Metric` object.

### HTTP

Posts a JSON array of metrics to the configured endpoint with `Content-Type: application/json`. Returns an error if the status code is outside 200–299. The timeout is configurable via `exporter.http.timeout` or `AGENT_HTTP_TIMEOUT`.

## Metrics JSON Schema

```json
{
  "timestamp": "2026-05-21T12:00:00Z",
  "host": "my-machine",
  "cpu": [
    {"core": 0, "usage": 12.5}
  ],
  "memory": {
    "total": 8589934592,
    "used": 4294967296,
    "usage": 50.0,
    "free": 4294967296
  },
  "disk": [
    {"mount": "/", "usage": 45.2}
  ],
  "system": {
    "uptime": 3600,
    "kernel": "6.1.0"
  }
}
```

## Runtime

The agent starts two goroutines:

1. **Collector** — reads metrics on a ticker interval, sends to the bounded channel
2. **Exporter** — reads from the channel and exports

On SIGINT/SIGTERM: the collector stops first, then the channel is closed and the
exporter drains remaining metrics. The exporter has a 5-second shutdown timeout.

## Testing

```bash
make test   # runs all tests with -race flag
```

## CI

GitHub Actions runs `gofmt`, `go vet`, `go test -race`, and `go build`
on every push.

## 🤝 Contributing

### Clone the repo

```bash
git clone https://github.com/RagnaCron/msma.git
cd msma
```

### Run the test suite

```bash
make test
```

### Build the compiled binary

```bash
make build
```

### Submit a pull request

If you'd like to contribute, please fork the repository and open a pull request
to the `main` branch.
