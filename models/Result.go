package models

import (
	"sync"
)

type result struct {
	latest Statistics
	lock   sync.Mutex
}

type Result interface {
	Get() Statistics
	Combine(stats Statistics)
}

func NewResult() Result {
	return &result{}
}

func (r *result) Get() Statistics {
	r.lock.Lock()
	defer r.lock.Unlock()

	return r.latest
}

func (r *result) Combine(stats Statistics) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.latest = Combine(r.latest, stats)
}
