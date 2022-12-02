package ssgo

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
)

func RunTestServer(t *testing.T, ctx context.Context, opts ...ServerOpt) (baseURL string, waitFunc func()) {
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

	return fmt.Sprintf("http://localhost:%d", srv.port), wg.Wait
}
