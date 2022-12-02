package ssgo

import (
	"context"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
)

func RunTestServer(t *testing.T, ctx context.Context, opts ...ServerOpt) (waitFunc func()) {
	t.Helper()

	srv := NewServer(
		DefaultServerPort,
		opts...,
	)

	wg := new(sync.WaitGroup)
	wg.Add(2)

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

	return wg.Wait
}
