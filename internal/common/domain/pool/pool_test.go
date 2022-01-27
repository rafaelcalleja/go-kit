package pool

import (
	"context"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
	"time"
)

func NewMockSinglePool() SinglePool {
	return NewSingleObject(func() interface{} {
		return "expected"
	})
}

func TestPool(t *testing.T) {
	ctx := context.TODO()
	ctx, _ = context.WithTimeout(ctx, 500000*time.Microsecond)

	t.Run("pool is unlocked to current routine", func(t *testing.T) {
		pool := NewMockSinglePool()
		a, b, c := pool.Get(ctx), pool.Get(ctx), pool.Get(ctx)
		require.Equal(t, "expected", a.(string))
		require.Equal(t, "expected", b.(string))
		require.Equal(t, "expected", c.(string))
	})

	t.Run("pool is locked to non-current routine", func(t *testing.T) {
		pool := NewMockSinglePool()
		wg := sync.WaitGroup{}
		wg.Add(1)

		pool.Get(ctx)

		go func() {
			defer func() {
				err := recover()
				require.NotNil(t, err)
				wg.Done()
			}()

			pool.Get(ctx)
		}()

		wg.Wait()
	})

	t.Run("pool released is available to other routine", func(t *testing.T) {
		pool := NewMockSinglePool()
		wg := sync.WaitGroup{}
		wg.Add(1)

		pool.Get(ctx)
		pool.Release()
		go func() {
			defer func() {
				err := recover()
				require.Nil(t, err)
				wg.Done()
			}()

			pool.Get(ctx)
		}()

		wg.Wait()
	})
	t.Run("error releasing a released pool", func(t *testing.T) {
		pool := NewMockSinglePool()
		defer func() {
			err := recover()
			require.NotNil(t, err)
		}()

		pool.Get(ctx)
		pool.Release()
		pool.Release()
	})
}
