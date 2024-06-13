package dmux

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	gomysqllock "github.com/sanketplus/go-mysql-lock"
)

type mysqlLocksFactory struct {
	db             *sql.DB
	locker         *gomysqllock.MysqlLocker
	timeoutSeconds int
}

// NewMySQLFactory returns a new distributed mutex factory that uses MySQL as
// the backend for the locks.
func NewMySQLFactory(db *sql.DB, timeoutSeconds int) MutexFactory {
	return &mysqlLocksFactory{
		db:             db,
		locker:         gomysqllock.NewMysqlLocker(db),
		timeoutSeconds: timeoutSeconds,
	}
}

func (m *mysqlLocksFactory) NewMutex(_ context.Context, name string) (Mutex, error) {
	return &mysqlLocksMutex{
		lockName: name,
		locker:   m.locker,
		timeout:  m.timeoutSeconds,
	}, nil
}

type mysqlLocksMutex struct {
	lockName string
	locker   *gomysqllock.MysqlLocker
	lock     *gomysqllock.Lock
	mu       sync.Mutex
	timeout  int
}

func (m *mysqlLocksMutex) Lock(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.lock != nil {
		return nil
	}

	lock, err := m.locker.ObtainTimeoutContext(ctx, m.lockName, m.timeout)
	if err != nil {
		return fmt.Errorf("%w: failed to obtain lock `%s", err, m.lockName)
	}

	m.lock = lock

	return nil
}

func (m *mysqlLocksMutex) Unlock(_ context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.lock == nil {
		return nil
	}

	if err := m.lock.Release(); err != nil {
		return fmt.Errorf("%w: failed to release lock `%s", err, m.lockName)
	}

	m.lock = nil

	return nil
}
