package ssgo

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
)

const testServerRoutinesNum = 2

var testServerPort uint32 = 9999

// RunTestServer listens http on incrementing port for testing purposes
func RunTestServer(t *testing.T, ctx context.Context, opts ...ServerOpt) (baseURL string, waitFunc func()) {
	t.Helper()

	srv := NewServer(
		uint(atomic.LoadUint32(&testServerPort)),
		opts...,
	)

	atomic.AddUint32(&testServerPort, 1) // for paralleltest

	wg := new(sync.WaitGroup)
	wg.Add(testServerRoutinesNum)

	go func() {
		defer wg.Done()

		err := srv.Run()
		require.NoError(t, err)
	}()

	go func() {
		defer wg.Done()

		<-ctx.Done()

		err := srv.Shutdown()
		require.NoError(t, err)
	}()

	return fmt.Sprintf("http://localhost:%d", srv.port), wg.Wait
}
