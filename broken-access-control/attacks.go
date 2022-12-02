package broken_access_control

import (
	"fmt"
	ssgo "github.com/ra9dev/ss-go"
	"io"
	"log"
	"net/http"
)

// block of vars to force interface implementation
var (
	_ ssgo.Hacker = (*URLHacker)(nil)
	_ ssgo.Hacker = (*QueryHacker)(nil)
)

type URLHacker struct {
	urls []string
}

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

	log.Printf("Data has been stolen from %s: %s", url, data)

	return nil
}

func (h URLHacker) Attack() error {
	for _, url := range h.urls {
		if err := h.attack(url); err != nil {
			return err
		}
	}

	return nil
}

type QueryHacker struct {
	url    string
	params map[string]string
}

func NewQueryHacker(url string, params map[string]string) QueryHacker {
	return QueryHacker{
		url:    url,
		params: params,
	}
}

func (q QueryHacker) Attack() error {
	req, err := http.NewRequest(http.MethodGet, q.url, nil)
	if err != nil {
		return fmt.Errorf("failed to prepare req: %w", err)
	}

	query := req.URL.Query()

	for key, val := range q.params {
		query.Set(key, val)
	}

	req.URL.RawQuery = query.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch data from %s: %w", req.URL.String(), err)
	}

	defer func() { _ = resp.Body.Close() }()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read data from %s: %w", req.URL.String(), err)
	}

	log.Printf("Data has been stolen from %s: %s", req.URL.String(), data)

	return nil
}
