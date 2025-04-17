package linter

import "sync/atomic"

type summary struct {
	numError     atomic.Int32
	numWarning   atomic.Int32
	numJob       atomic.Int32
	numProcessed atomic.Int32
}

func newSummary() *summary {
	s := &summary{
		numError:     atomic.Int32{},
		numWarning:   atomic.Int32{},
		numJob:       atomic.Int32{},
		numProcessed: atomic.Int32{},
	}
	return s
}
