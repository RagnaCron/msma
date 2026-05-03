# Software Design Document (SDD)

## Minimal System Metrics Agent (Go)

---

## 1. Introduction

### 1.1 Purpose

This document describes the system design for the Minimal System Metrics Agent.
It defines the architecture, components, interfaces, and runtime behavior
required to satisfy the SRS.

### 1.2 Scope

The system is a lightweight daemon that collects system metrics and Exporters
them using a decoupled, queue-based architecture.

### 1.3 Definitions

- SRS: Software Requirements Specification
- SDD: Software Design Document

### 1.4 References

- SRS Document
- ISO/IEC/IEEE 1016:2009

---

## 2. Overall Design

### 2.1 Product Perspective

The system follows a producer–consumer architecture:

Collector → Channel → Exporter

- Collector produces metrics
- Channel buffers metrics
- Exporter consumes and transmits

---

### 2.2 Design Principles

- Minimal complexity
- Deterministic behavior
- Explicit backpressure handling
- Future extensibility (Kafka, batching)

---

## 3. Configuration Design

### 3.1 Configuration Sources

Configuration is loaded in this order:

1. Defaults
2. YAML file (`./config.yaml`, optional)
3. Environment variables (override all)

---

### 3.2 Configuration Schema

```yaml
interval: 5
queue_size: 10

exporter:
  type: http # or stdout
  http:
    endpoint: "http://localhost:8080/metrics"
```

---

### 3.3 Validation

Validation rules:

- queue_size > 0
- exporter type is valid
- HTTP exporter requires endpoint

Failure → log error → terminate

---

## 4. Functional Design

### 4.1 Metric Collection

Collector gathers:

- CPU usage
- memory stats
- disk usage
- uptime
- host + kernel

---

### 4.2 Metric Flow

Channel:

- type: chan Metric
- bounded
- capacity = queue_size

Non-blocking send:

```go
select {
case ch <- metric:
default:
    // drop
}
```

Dropped metrics are logged at info level.

---

### 4.3 Export Design

Exporter interface:

```go
type Exporter interface {
    Export(metrics []Metric) error
}
```

Batching:

- occurs inside exporter
- v1 always sends 1 metric per batch

---

### 4.4 Exporters

Stdout:

- prints JSON array per line

HTTP:

- blocking POST
- no retries
- no timeout config
- sends JSON array

---

## 5. Data Design

Metric structure includes:

- timestamp
- host
- nested metrics
- tags

Serialization:

- JSON array
- RFC3339 timestamps

---

## 6. Non-Functional Design

- CPU < 1%
- Memory < 50MB
- metric loss allowed
- fail-fast config

---

## 7. Concurrency Design

- 1 collector goroutine
- 1 exporter goroutine
- channel communication

---

## 8. Lifecycle

Startup:

- load config
- validate
- init components

Runtime:

- collect + export loop

Shutdown:

- stop collector
- drain channel with timeout
- flush remaining metrics

---

## 9. Future

- Kafka exporter
- batching
- retries
