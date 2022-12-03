package broken_access_control

import (
	"fmt"
	"io"
	"log"
	"net/http"

	ssgo "github.com/ra9dev/ss-go"
)

var _ ssgo.Hacker = (*CORSHacker)(nil)

// CORSHacker steals data from an url with an unprotected query param
type CORSHacker struct {
	url     string
	origins []string
}

// NewCORSHacker constructor
func NewCORSHacker(url string, origins ...string) CORSHacker {
	return CORSHacker{
		url:     url,
		origins: origins,
	}
}

func (c CORSHacker) attack(origin string) error {
	url := c.url

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to prepare req: %w", err)
	}

	req.Header.Set(originHeaderKey, origin)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch data from %s: %w", url, err)
	}

	defer func() { _ = resp.Body.Close() }()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read data from %s: %w", url, err)
	}

	log.Printf("%s is vulnerable for origin %s, data has been stolen: %s", url, origin, data)

	return nil
}

// Attack implementation of ssgo.Hacker
func (c CORSHacker) Attack() error {
	for _, origin := range c.origins {
		if err := c.attack(origin); err != nil {
			return err
		}
	}

	return nil
}
