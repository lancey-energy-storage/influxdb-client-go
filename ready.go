package influxdb

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
)

func (c *Client) Ready(ctx context.Context) (*ReadyResult, error) {
	log.Printf("[DEBUG] Pinging resource ")
	pingUrl, _ := url.Parse(c.url.String())
	pingUrl.Path = "/ready"

	req, err := http.NewRequest(http.MethodGet, pingUrl.String(), nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	resp, err := c.httpClient.Do(req)

	defer resp.Body.Close()
	readyResult := &ReadyResult{}

	if err := json.NewDecoder(resp.Body).Decode(readyResult); err != nil {
		return nil, err
	}
	return readyResult, nil
}

type ReadyResult struct {
	Started string `json:"started"`
	Status  string `json:"status"`
	Up      string `json:"up"`
}
