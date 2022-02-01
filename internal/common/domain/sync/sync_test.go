package sync

import (
	"context"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
	"time"
)

type depA struct {
	mux *ChanSync
	dep *depC
}

type depB struct {
	mux *ChanSync
	dep *depC
}

type depC struct {
	val string
}

func TestNewChanSync(t *testing.T) {
	wg := sync.WaitGroup{}
	wg2 := sync.WaitGroup{}
	wg3 := sync.WaitGroup{}

	mux := NewChanSync()

	z := &depC{"foo"}

	x := depA{mux: mux, dep: z}
	y := depB{mux: mux, dep: z}

	require.Same(t, x.mux, y.mux)

	ctx := context.Background()
	ctx = x.mux.Lock(ctx)
	require.Equal(t, "foo", z.val)

	for i := 0; i < 100; i++ {
		wg3.Add(1)
		go func() {
			defer wg3.Done()

			x.mux.CWait(context.Background())
			x.mux.CUnlock()

			_ = x.mux.Lock(context.Background())
			time.Sleep(5000 * time.Microsecond)
			x.mux.Unlock()
		}()
	}

	for i := 0; i < 100; i++ {
		wg2.Add(1)
		go func() {
			wg2.Add(2)
			defer wg2.Done()

			go func() {
				defer wg2.Done()
				y.mux.CWait(context.Background()) //Other context
				defer y.mux.CUnlock()

				require.Equal(t, "bar", z.val)
			}()

			go func() {
				defer wg2.Done()
				y.mux.CWait(context.WithValue(ctx, ctxSyncKey, y.mux))
				defer y.mux.CUnlock()

				require.Equal(t, "foo", z.val)
			}()
		}()
	}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(5000 * time.Microsecond)

			y.mux.CWait(ctx)
			defer y.mux.CUnlock()

			require.Equal(t, "foo", z.val)
		}()
	}

	wg.Wait()
	require.Equal(t, "foo", z.val)
	x.mux.Unlock()
	z.val = "bar"

	wg3.Wait()
	wg2.Wait()
	require.Equal(t, "bar", z.val)
}
