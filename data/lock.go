package data

import (
	"sync"
)

// IndexLock is a custom Lock object for ensuring that all operations on our data are safe.
type IndexLock interface {
	Lock()
	Unlock()
}

type SimpleLock struct {
	lock *sync.Mutex
}

func (s *SimpleLock) Lock() {
	s.lock.Lock()
}

func (s *SimpleLock) Unlock() {
	s.lock.Unlock()
}

func NewLock() IndexLock {
	return &SimpleLock{&sync.Mutex{}}
}
