package dmux

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/Avalanche-io/counter"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
	"github.com/tangelo-labs/go-dotenv"
)

type environment struct {
	// TODO: refactor this to use a test database
	ReceiptsMYSQLDSN string `env:"MARKETPLACE_PAYMENTS_RECEIPTS_DB_MYSQL_DSN"`
}

func TestMysqlLocksDMutex(t *testing.T) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancelFunc()

	var envVars environment

	require.NoError(t, dotenv.LoadAndParse(&envVars))

	db, err := sql.Open("mysql", envVars.ReceiptsMYSQLDSN)
	require.NoError(t, err)

	db.SetMaxIdleConns(10)
	db.SetConnMaxIdleTime(20 * time.Second)
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(20 * time.Second)

	t.Run("GIVEN a mysql dmutex", func(t *testing.T) {
		factoryMtx := NewMySQLFactory(db, 10)

		t.Run("WHEN several goroutines try to acquire the same lock", func(t *testing.T) {
			startChan := make(chan struct{})

			lockErrCounter := counter.NewUnsigned()
			unlockErrCounter := counter.NewUnsigned()
			locks := counter.NewUnsigned()
			mtx, nErr := factoryMtx.NewMutex(ctx, "test")

			if nErr != nil {
				fmt.Printf("error: %s", nErr)
				return
			}

			wg := sync.WaitGroup{}
			numGoroutines := 500

			for i := 0; i < numGoroutines; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()

					<-startChan

					if gErr := mtx.Lock(ctx); gErr != nil {
						lockErrCounter.Up()

						return
					}

					locks.Up()

					if gErr := mtx.Unlock(ctx); gErr != nil {
						unlockErrCounter.Up()
					}
				}()
			}

			close(startChan)
			wg.Wait()

			t.Run("THEN all locks are acquired sequentially", func(t *testing.T) {
				require.Zero(t, lockErrCounter.Get())
				require.Zero(t, unlockErrCounter.Get())
				require.EqualValues(t, numGoroutines, locks.Get())
			})
		})
	})
}
