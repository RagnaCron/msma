package collector

import (
	"errors"
	"strings"
	"sync"
	"testing"

	"github.com/ragnacron/msma/internal/model"
)

var errFake = errors.New("fake error")

func TestCollect(t *testing.T) {
	col, err := New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	metric, err := col.Collect()
	if err != nil {
		t.Fatalf("Collect() failed: %v", err)
	}

	if metric.Host == "" {
		t.Error("expected host to be set")
	}
	if len(metric.CPU) == 0 {
		t.Error("expected CPU metrics to be collected")
	}
}

func TestCollectPartialFailure(t *testing.T) {
	original := getCPUMetricsFn
	defer func() { getCPUMetricsFn = original }()

	getCPUMetricsFn = func() ([]model.CPU, error) {
		return nil, errFake
	}

	col, err := New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	metric, err := col.Collect()
	if err == nil {
		t.Fatal("expected error for partial failure")
	}

	if metric.Host == "" {
		t.Error("expected host to be set")
	}
	if metric.CPU != nil {
		t.Error("expected CPU to be nil after failure")
	}
	if metric.Memory == nil {
		t.Error("expected Memory to be collected despite CPU failure")
	}
	if len(metric.Disk) == 0 {
		t.Error("expected Disk to be collected despite CPU failure")
	}
	if metric.System == nil {
		t.Error("expected System to be collected despite CPU failure")
	}
}

func TestCollectAllFail(t *testing.T) {
	origCPU := getCPUMetricsFn
	origMem := getMemoryMetricsFn
	origDisk := getDiskMetricsFn
	origSys := getSystemInfoMetricsFn
	defer func() {
		getCPUMetricsFn = origCPU
		getMemoryMetricsFn = origMem
		getDiskMetricsFn = origDisk
		getSystemInfoMetricsFn = origSys
	}()

	getCPUMetricsFn = func() ([]model.CPU, error) { return nil, errFake }
	getMemoryMetricsFn = func() (model.Memory, error) { return model.Memory{}, errFake }
	getDiskMetricsFn = func() ([]model.Disk, error) { return nil, errFake }
	getSystemInfoMetricsFn = func() (model.System, error) { return model.System{}, errFake }

	col, err := New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	metric, err := col.Collect()
	if err == nil {
		t.Fatal("expected error when all collectors fail")
	}

	me, ok := err.(*MetricsError)
	if !ok {
		t.Fatalf("expected *MetricsError, got %T", err)
	}
	if len(me.errs) != 4 {
		t.Errorf("expected 4 errors, got %d", len(me.errs))
	}
	if metric.CPU != nil || metric.Memory != nil || len(metric.Disk) > 0 || metric.System != nil {
		t.Error("expected zero-valued payload when all collectors fail")
	}
}

func TestMetricsErrorFormat(t *testing.T) {
	me := &MetricsError{
		errs: []collectorError{
			{name: "cpu", err: errors.New("cpu failed")},
			{name: "memory", err: errors.New("memory failed")},
		},
	}

	got := me.Error()
	if got != "metric collection failures: cpu: cpu failed, memory: memory failed" {
		t.Errorf("unexpected error format: %q", got)
	}
}

func TestCollectConcurrency(t *testing.T) {
	col, err := New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Go(func() {
			_, err := col.Collect()
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
	wg.Wait()
}

func TestMetricsErrorAllFailed(t *testing.T) {
	me := &MetricsError{
		errs: []collectorError{
			{name: "cpu", err: errFake},
			{name: "memory", err: errFake},
			{name: "disk", err: errFake},
			{name: "sysinfo", err: errFake},
		},
	}

	got := me.Error()
	if !strings.HasPrefix(got, "metric collection failures:") {
		t.Errorf("expected error to start with 'metric collection failures:', got: %q", got)
	}
}
