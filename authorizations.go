package influxdb

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

func (c *Client) GetAllAuthorizations(org string, orgID string, user string, userID string) (*AuthorizationsList, error) {
	log.Printf("[DEBUG] Get all authorizations")

	params := "?org=" + org + "&orgID=" + orgID + "&user=" + user + "&userID=" + userID
	req, err := http.NewRequest(http.MethodGet, c.url.String()+"/authorizations"+params, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", c.authorization)
	resp, err := c.httpClient.Do(req)

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	authorizationsList := &AuthorizationsList{}
	if err := json.NewDecoder(resp.Body).Decode(authorizationsList); err != nil {
		return nil, err
	}

	return authorizationsList, nil
}

type AuthorizationsList struct {
	Links struct {
		Next string `json:"next"`
		Self string `json:"self"`
		Prev string `json:"prev"`
	} `json:"links"`
	Authorizations []struct {
		Status      string `json:"status"`
		Description string `json:"description"`
		CreatedAt   string `json:"createdAt"`
		UpdatedAt   string `json:"updatedAt"`
		OrgID       string `json:"orgID"`
		Permissions []struct {
			Action   string `json:"action"`
			Resource struct {
				Type  string `json:"type"`
				Id    string `json:"id"`
				Name  string `json:"name"`
				OrgID string `json:"orgID"`
				Org   string `json:"org"`
			} `json:"resource"`
		} `json:"permissions"`
		Id     string `json:"id"`
		Token  string `json:"token"`
		UserID string `json:"userID"`
		User   string `json:"user"`
		Org    string `json:"org"`
		Links  struct {
			Self string `json:"self"`
			User string `json:"user"`
		} `json:"links"`
	} `json:"authorizations"`
}
