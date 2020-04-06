package influxdb

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

func (c *Client) CreateDashboard(description string, name string, orgID string) (*Dashboard, error) {
	if name == "" {
		return nil, errors.New("a name is required")
	}
	if orgID == "" {
		return nil, errors.New("an org id is required")
	}

	log.Printf("[DEBUG] Creating a dashboard")

	inputData, err := json.Marshal(NewDashboard{
		Name:        name,
		Description: description,
		OrgID:       orgID,
	})

	req, err := http.NewRequest(http.MethodPost, c.url.String()+"/dashboards", bytes.NewBuffer(inputData))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", c.authorization)
	resp, err := c.httpClient.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		return nil, errors.New(resp.Status)
	}

	dashboard := &Dashboard{}
	if err := json.NewDecoder(resp.Body).Decode(dashboard); err != nil {
		return nil, err
	}

	return dashboard, nil
}

type Dashboard struct {
	OrgID       string `json:"orgID"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Links       struct {
		Self    string `json:"self"`
		Cells   string `json:"cells"`
		Owners  string `json:"owners"`
		Members string `json:"members"`
		Logs    string `json:"logs"`
		Labels  string `json:"labels"`
		Org     string `json:"org"`
	} `json:"links"`
	Id   string `json:"id"`
	Meta struct {
		CreatedAt string `json:"createdAt"`
		UpdatedAt string `json:"updatedAt"`
	} `json:"meta"`
	Cells []struct {
		Id    string `json:"id"`
		Links struct {
			Self string `json:"self"`
			View string `json:"view"`
		} `json:"links"`
		X      int    `json:"x"`
		Y      int    `json:"y"`
		W      int    `json:"w"`
		H      int    `json:"h"`
		ViewID string `json:"view ID"`
	} `json:"cells"`
	Labels []struct {
		Id         string `json:"id"`
		OrgID      string `json:"orgID"`
		Name       string `json:"name"`
		Properties struct {
			Color       string `json:"color"`
			Description string `json:"description"`
		} `json:"properties"`
	} `json:"labels"`
}

type NewDashboard struct {
	Description string `json:"description"`
	Name        string `json:"name"`
	OrgID       string `json:"orgID"`
}
