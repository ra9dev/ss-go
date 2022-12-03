package ssgo

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

const testServerRoutinesNum = 2

func RunTestServer(t *testing.T, ctx context.Context, opts ...ServerOpt) (baseURL string, waitFunc func()) {
	t.Helper()

	srv := NewServer(
		defaultServerPort,
		opts...,
	)

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
