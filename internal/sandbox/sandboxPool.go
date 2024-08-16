package sandbox

import (
	"errors"
	"fmt"
	"sync"
)

type SandboxPool struct {
	mu        sync.Mutex
	sandboxes map[int]*Sandbox
	maxSize   int
	created   int
	cond      *sync.Cond
}

func NewSandboxPool(maxSize int) *SandboxPool {
	s := &SandboxPool{
		sandboxes: make(map[int]*Sandbox),
		maxSize:   maxSize,
	}
	s.cond = sync.NewCond(&s.mu)
	return s
}

func (s *SandboxPool) GetSandbox() (*Sandbox, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for len(s.sandboxes) == 0 && s.created >= s.maxSize {
		fmt.Printf("wait len:%d ed %d\n", len(s.sandboxes), s.created)
		s.cond.Wait()
	}

	for id, sandbox := range s.sandboxes {
		delete(s.sandboxes, id)
		return sandbox, nil
	}

	if s.created < s.maxSize {
		s.created++
		return newSandbox(s.created)
	}

	return nil, errors.New("no sandbox available")
}

func (s *SandboxPool) ReleaseSandbox(sandbox *Sandbox) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.sandboxes) < s.maxSize {
		s.sandboxes[sandbox.ID] = sandbox
	}
	s.cond.Signal()
}
