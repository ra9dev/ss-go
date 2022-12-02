package broken_access_control

import (
	"context"
	"fmt"
	ssgo "github.com/ra9dev/ss-go"
	"sync"
	"testing"
)

const defaultServerPort = 8080

func runTestServer(t *testing.T, ctx context.Context, opts ...ssgo.ServerOpt) (waitFunc func()) {
	t.Helper()

	srv := ssgo.NewServer(
		defaultServerPort,
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

func TestURLDataStealer_Attack(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wait := runTestServer(t, ctx, ssgo.ServerWithRoute(NewURLAttackTarget()))
	defer wait()

	baseURL := fmt.Sprintf("http://localhost:%d", defaultServerPort)

	hackers := []URLDataStealer{
		NewURLDataStealer(baseURL + appPath + publicAppInfoPath),
		NewURLDataStealer(baseURL + appPath + privateAppInfoPath),
	}

	for _, hacker := range hackers {
		if err := hacker.Attack(); err != nil {
			t.Errorf("failed to attack %s: %+v", hacker.url, err)

			continue
		}
	}

	cancel()
}
