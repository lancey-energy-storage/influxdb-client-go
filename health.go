package influxdb

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
)

func (c *Client) GetHealth(ctx context.Context) (*Health, error) {
	log.Printf("[DEBUG] Getting health of resource ")
	pingUrl, _ := url.Parse(c.url.String())
	pingUrl.Path = "/health"

	req, err := http.NewRequest(http.MethodGet, pingUrl.String(), nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	resp, err := c.httpClient.Do(req)

	defer resp.Body.Close()
	health := &Health{}

	if err := json.NewDecoder(resp.Body).Decode(health); err != nil {
		return nil, err
	}
	return health, nil
}

type Health struct {
	Name    string     `json:"name"`
	Message string     `json:"message"`
	Checks  []struct{} `json:"checks"`
	Status  string     `json:"status"`
}
