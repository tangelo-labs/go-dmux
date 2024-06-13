package dmux_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/require"
	"github.com/tangelo-labs/go-dmux"
)

func TestRedisMuxFactory(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	t.Run("GIVEN a redis mutex shared across 100 goroutines", func(t *testing.T) {
		redisDSN := startRedisServer(ctx, t)
		factory, err := dmux.NewRedisFactory(dmux.RedisConfig{DSN: redisDSN})

		require.NoError(t, err)
		require.NotNil(t, factory)

		mu, err := factory.NewMutex(ctx, "test")
		require.NoError(t, err)

		t.Run("WHEN each goroutine tries to acquire the mutex to increase a counter", func(t *testing.T) {
			sharedCounter := 0
			ready := make(chan struct{})
			wg := sync.WaitGroup{}

			for i := 0; i < 100; i++ {
				wg.Add(1)

				go func() {
					<-ready

					defer wg.Done()

					if lErr := mu.Lock(ctx); lErr != nil {
						t.Log(lErr.Error())
					}

					defer func() {
						if uErr := mu.Unlock(ctx); uErr != nil {
							t.Log(uErr.Error())
						}
					}()

					sharedCounter++
				}()
			}

			t.Run("THEN the counter final value is 100 once every goroutine completes", func(t *testing.T) {
				close(ready)
				wg.Wait()

				require.Equal(t, 100, sharedCounter)
			})
		})
	})
}

func startRedisServer(ctx context.Context, t *testing.T) string {
	t.Helper()

	server, err := miniredis.Run()
	require.NoError(t, err)
	require.NotNil(t, server)

	delta := 50 * time.Millisecond
	ticker := time.NewTicker(delta)

	go func() {
		for {
			select {
			case <-ctx.Done():
				server.Close()
			case <-ticker.C:
				server.FastForward(delta)
			}
		}
	}()

	return fmt.Sprintf("redis://%s/0", server.Addr())
}
