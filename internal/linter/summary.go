package linter

import (
	"sync"
	"sync/atomic"

	"github.com/keenbytes/octo-linter/v2/internal/linter/glitch"
)

type summary struct {
	mu sync.Mutex
	numError     atomic.Int32
	numWarning   atomic.Int32
	numJob       atomic.Int32
	numProcessed atomic.Int32
	glitches []*glitch.Glitch
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

func (s *summary) addGlitch(g *glitch.Glitch) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.glitches = append(s.glitches, g)
}
