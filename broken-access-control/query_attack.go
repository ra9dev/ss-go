package broken_access_control

import (
	"fmt"
	"io"
	"log"
	"net/http"

	ssgo "github.com/ra9dev/ss-go"
)

var _ ssgo.Hacker = (*QueryHacker)(nil)

// QueryHacker steals data from an url with an unprotected query param
type QueryHacker struct {
	url    string
	params map[string]string
}

// NewQueryHacker constructor
func NewQueryHacker(url string, params map[string]string) QueryHacker {
	return QueryHacker{
		url:    url,
		params: params,
	}
}

// Attack implementation of ssgo.Hacker
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
