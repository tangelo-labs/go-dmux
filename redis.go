package dmux

import (
	"context"
	"time"

	goredislib "github.com/go-redis/redis"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis"
)

// RedisConfig configuration needed by Redis factory.
type RedisConfig struct {
	DSN string
	// Lock expiration. Mutex instances from the factory will lock the key with expiration. 0 means no expiration.
	Expiration time.Duration
	// Retries to acquire the lock. If set to 0 default is 32.
	Retries uint
}

type redisFactory struct {
	rs   *redsync.Redsync
	conf RedisConfig
}

func (r *redisFactory) NewMutex(_ context.Context, name string) (Mutex, error) {
	var opts []redsync.Option

	opts = append(opts, redsync.WithExpiry(r.conf.Expiration))

	if r.conf.Retries == 0 {
		r.conf.Retries = 1
	}

	opts = append(opts, redsync.WithTries(int(r.conf.Retries)))

	return &redisDMux{
		rmu: r.rs.NewMutex(name, opts...),
	}, nil
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

	return &redisFactory{
		rs:   rs,
		conf: cfg,
	}, nil
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
