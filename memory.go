package dmux

import (
	"context"
	"sync"
)

type inMemoryFactoryMutex struct {
	mutexes map[string]*inMemoryMutex
	sync.Mutex
}

// NewInMemoryFactory returns a new in-memory distributed mutex factory.
// This is not safe for production use, but is useful for testing.
func NewInMemoryFactory() MutexFactory {
	return &inMemoryFactoryMutex{}
}

func (i *inMemoryFactoryMutex) NewMutex(_ context.Context, name string) (Mutex, error) {
	i.Lock()
	defer i.Unlock()

	if i.mutexes == nil {
		i.mutexes = map[string]*inMemoryMutex{}
	}

	if _, ok := i.mutexes[name]; !ok {
		i.mutexes[name] = &inMemoryMutex{}
	}

	return i.mutexes[name], nil
}

type inMemoryMutex struct {
	mu sync.Mutex
}

func (i *inMemoryMutex) Lock(_ context.Context) error {
	i.mu.Lock()

	return nil
}

func (i *inMemoryMutex) Unlock(_ context.Context) error {
	i.mu.Unlock()

	return nil
}
