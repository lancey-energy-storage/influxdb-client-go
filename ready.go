package influxdb

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"
)

func (c *Client) Ready (ctx context.Context) (*ReadyResult, error) {
	log.Printf("[DEBUG] Pinging resource ")
	pingUrl, _ := url.Parse(c.url.String())
	pingUrl.Path = "/ready"

	req, err := http.NewRequest(http.MethodGet, pingUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	defer req.Body.Close()
	readyResult := &ReadyResult{}
	if err := json.NewDecoder(req.Body).Decode(readyResult); err != nil {
		return nil, err
	}
	return readyResult, nil
}

type ReadyResult struct {
	Started *time.Time `json:"started,omitempty"`
	Status  *string    `json:"status,omitempty"`
	Up      *string    `json:"up,omitempty"`
}
