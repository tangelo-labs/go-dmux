package dmux

import "context"

// Mutex is an abstraction for a distributed mutex.
type Mutex interface {
	// Lock obtain a lock for this mutex. After this is successful, no one else
	// can obtain this lock until it is unlocked.
	Lock(ctx context.Context) error

	// Unlock release the lock so other processes or threads can obtain it.
	Unlock(ctx context.Context) error
}

// MutexFactory is an abstraction for a distributed mutex factory.
type MutexFactory interface {
	// NewMutex creates a new distributed mutex for the given key.
	NewMutex(ctx context.Context, name string) (Mutex, error)
}
