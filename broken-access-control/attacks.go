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
	_ ssgo.Hacker = (*URLDataStealer)(nil)
)

type URLDataStealer struct {
	url string
}

func NewURLDataStealer(url string) URLDataStealer {
	return URLDataStealer{url: url}
}

func (s URLDataStealer) Attack() error {
	resp, err := http.Get(s.url)
	if err != nil {
		return fmt.Errorf("failed to fetch data from %s: %w", s.url, err)
	}

	defer func() { _ = resp.Body.Close() }()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read data from %s: %w", s.url, err)
	}

	log.Printf("Data has been stolen from %s: %s\n", s.url, data)

	return nil
}
