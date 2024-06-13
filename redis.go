package dmux

import (
	"context"

	goredislib "github.com/go-redis/redis"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis"
)

// RedisConfig configuration needed by Redis factory.
type RedisConfig struct {
	DSN string
}

type redisFactory struct {
	rs *redsync.Redsync
}

func (r *redisFactory) NewMutex(_ context.Context, name string) (Mutex, error) {
	return &redisDMux{rmu: r.rs.NewMutex(name)}, nil
}

// NewRedisFactory builds a new distributed mutex factory that uses Redis as
// backend to provide distributed mutexes.
func NewRedisFactory(cfg RedisConfig) (MutexFactory, error) {
	opts, err := goredislib.ParseURL(cfg.DSN)
	if err != nil {
		return nil, err
	}

	client := goredislib.NewClient(opts)
	pool := goredis.NewPool(client)
	rs := redsync.New(pool)

	return &redisFactory{rs: rs}, nil
}

type redisDMux struct {
	rmu *redsync.Mutex
}

func (r *redisDMux) Lock(ctx context.Context) error {
	return r.rmu.LockContext(ctx)
}

func (r *redisDMux) Unlock(ctx context.Context) error {
	_, err := r.rmu.UnlockContext(ctx)

	return err
}
