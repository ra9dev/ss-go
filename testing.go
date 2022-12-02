package ssgo

import (
	"context"
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

		if err := srv.Run(); err != nil {
			t.Errorf("failed to run srv: %+v", err)

			return
		}
	}()

	go func() {
		defer wg.Done()

		<-ctx.Done()

		if err := srv.Shutdown(); err != nil {
			t.Errorf("failed to shutdown srv: %+v", err)

			return
		}
	}()

	return wg.Wait
}
