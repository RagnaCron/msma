// Package queue
package queue

import "github.com/ragnacron/msma/internal/model"

type Queue struct {
	ch chan model.Metric
}

func NewQueue(size int) *Queue {
	return &Queue{ch: make(chan model.Metric, size)}
}

func (q *Queue) Channel() chan model.Metric {
	return q.ch
}
