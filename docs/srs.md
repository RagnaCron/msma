# Software Requirements Specification (SRS)

## Minimal System Metrics Agent (Go)

## 1. Introduction

### 1.1 Purpose

This document specifies the requirements for a Minimal System Metrics Agent
implemented in Go. The agent collects basic system metrics from a host machine
and exports them to a configurable backend.

### 1.2 Scope

The system is a lightweight, standalone daemon that:

- collects host-level system metrics (CPU, memory, disk, and OS information)
- emits metrics at a fixed interval
- exports metrics via a configurable backend (stdout or HTTP)
- uses a decoupled architecture between metric collection and export

Kafka support is defined as a future extension and is not part of version 1.

### 1.3 Definitions, Acronyms, and Abbreviations

- SRS: Software Requirements Specification
- SDD: Software Design Document
- FR: Functional Requirement
- NFR: Non-Functional Requirement
- CFG: Configuration Requirement
- DATA: Data Requirement
- VER: Verification Requirement

### 1.4 References

- ISO/IEC/IEEE 29148:2018
- Go Programming Language Specification

## 2. Overall Description

### 2.1 Product Perspective

The system is a standalone agent running on a host machine. It separates metric
collection from metric export using a bounded internal queue. Metric export is
batch-oriented, and version 1 uses batches of size one.

### 2.2 Product Functions

- periodic metric collection
- metric serialization to JSON
- metric buffering via bounded queue
- metric export via configured backend
- configuration loading and validation
- controlled handling of backpressure and shutdown

### 2.3 Operating Environment

#### REQ-OE-001

The system SHALL run on Linux (amd64).

#### REQ-OE-002

The system SHOULD support macOS (amd64) on a best-effort basis.

#### REQ-OE-003

The system SHALL be distributed as a single standalone binary.

## 3. Specific Requirements

### 3.1 Configuration Requirements

#### REQ-CFG-001

The system SHALL support configuration via a YAML file.

#### REQ-CFG-002

The system SHALL support configuration via environment variables.

#### REQ-CFG-003

Environment variables SHALL override default values and values defined in the
YAML configuration.

#### REQ-CFG-004

The system SHALL use `./.config/msma/config.yaml` as the default configuration path.

#### REQ-CFG-005

The configuration file SHALL be optional.

#### REQ-CFG-006

The system SHALL validate configuration on startup.

#### REQ-CFG-007

The system SHALL terminate execution if configuration validation fails.

#### REQ-CFG-008

The system SHALL support configuration of the internal queue size.

#### REQ-CFG-009

The system SHALL use a default internal queue size of 10 if not specified.

#### REQ-CFG-010

The system SHALL validate that the configured queue size is greater than zero.

### 3.2 Functional Requirements

#### Metric Collection

#### REQ-FR-001

The system SHALL collect total CPU usage as a percentage.

#### REQ-FR-002

The system SHALL collect memory metrics including total, used, and free memory.

#### REQ-FR-003

The system SHALL collect disk usage percentage per mounted filesystem.

#### REQ-FR-004

The system SHALL collect system uptime in seconds.

#### REQ-FR-005

The system SHALL collect host identification including hostname and kernel version.

#### Metric Flow

#### REQ-FR-006

The system SHALL collect metrics at a fixed configurable interval independent of
export operations.

#### REQ-FR-007

The system SHALL use a bounded internal queue with configurable capacity.

#### REQ-FR-008

The system SHALL ensure metric collection is non-blocking with respect to export
operations.

#### REQ-FR-009

The system SHALL drop metrics when the queue is full.

#### REQ-FR-010

The system SHALL continue operation after dropping metrics.

### 3.3 Export Semantics

#### REQ-FR-011

The system SHALL provide metrics to the exporter in batches.

#### REQ-FR-012

The system SHALL support batches containing one or more metrics.

#### REQ-FR-013

The system SHALL ensure that, in version 1, each batch contains exactly one metric.

#### REQ-FR-014

The system SHALL perform batching within the exporter component.

#### REQ-FR-015

The system SHALL support a stdout exporter.

#### REQ-FR-016

The system SHALL support an HTTP exporter.

#### REQ-FR-017

The system SHALL allow only one active exporter at a time.

#### REQ-FR-018

The stdout exporter SHALL output one JSON array per line.

#### REQ-FR-019

The HTTP exporter SHALL send metrics via HTTP POST.

#### REQ-FR-020

The HTTP exporter SHALL send metrics as a JSON array.

#### REQ-FR-021

The HTTP exporter SHALL include the JSON array in the request body.

### 3.4 Failure Handling

#### REQ-FR-022

The system SHALL log an error if metric export fails.

#### REQ-FR-023

The system SHALL continue operation after export failures.

#### REQ-FR-024

The system SHALL log every dropped metric event.

### 3.5 Shutdown Requirements

#### REQ-FR-025

The system SHALL stop metric collection on shutdown.

#### REQ-FR-026

The system SHALL attempt a best-effort flush of queued metrics before exit.

### 3.6 Data Requirements

#### REQ-DATA-001

The system SHALL structure exported data as JSON arrays of metric objects.

#### REQ-DATA-002

Each metric object SHALL include a timestamp, a host identifier, and metrics data.

#### REQ-DATA-003

The system SHALL include an optional tags field.

#### REQ-DATA-004

The system SHALL serialize exported timestamps in RFC3339 UTC format.

#### REQ-DATA-005

The system SHALL use the same JSON array structure for stdout and HTTP exports.

### 3.7 Non-Functional Requirements

#### REQ-NFR-001

The system SHALL maintain CPU usage below 1%.

#### REQ-NFR-002

The system SHALL maintain memory usage below 50 MB.

#### REQ-NFR-003

The system SHALL tolerate metric loss under backpressure.

#### REQ-NFR-004

The system SHALL keep the runtime behavior deterministic under configuration
failure by terminating on invalid configuration.

### 3.8 Future Requirements

#### REQ-FUT-001

The system SHOULD support a Kafka exporter in future versions.

#### REQ-FUT-002

The system SHOULD support metric aggregation beyond one metric per batch in
future versions.

#### REQ-FUT-003

The system SHOULD support retry and buffering mechanisms in future versions.

#### REQ-FUT-004

The system SHOULD support non-empty tagging in future versions.

## 4. Verification Requirements

### REQ-VER-001

The system SHALL emit valid JSON arrays at each interval.

### REQ-VER-002

The system SHALL drop metrics when the queue is full.

### REQ-VER-003

The system SHALL continue collecting metrics when the exporter is slow.

### REQ-VER-004

The system SHALL honor environment-variable override precedence over YAML and defaults.

### REQ-VER-005

The system SHALL continue operating when the configuration file is absent and
defaults or environment variables are sufficient.

### REQ-VER-006

The system SHALL perform a best-effort flush during shutdown.

## 5. Dependencies

- Go runtime
- gopsutil
