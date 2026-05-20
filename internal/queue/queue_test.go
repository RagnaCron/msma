package queue

import (
	"testing"

	"github.com/ragnacron/msma/internal/model"
)

func TestNew(t *testing.T) {
	q := New(5)
	if cap(q.ch) != 5 {
		t.Errorf("expected channel capacity 5, got %d", cap(q.ch))
	}
}

func TestChannelSendReceive(t *testing.T) {
	q := New(2)
	ch := q.Channel()
	m1 := model.Metric{Host: "test1"}
	m2 := model.Metric{Host: "test2"}

	go func() {
		q.ch <- m1
		q.ch <- m2
	}()

	received := make([]model.Metric, 0, 2)
	for range 2 {
		m, ok := <-ch
		if !ok {
			t.Fatalf("unexpected channel close")
		}
		received = append(received, m)
	}

	if received[0].Host != "test1" || received[1].Host != "test2" {
		t.Errorf("unexpected metrics received: %+v", received)
	}
}
