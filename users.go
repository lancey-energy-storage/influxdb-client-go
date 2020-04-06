package influxdb

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

func (c *Client) GetAllUsers() (*UserList, error) {
	log.Printf("[DEBUG] Get all authorizations")

	req, err := http.NewRequest(http.MethodGet, c.url.String()+"/users", nil)
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

	userList := &UserList{}
	if err := json.NewDecoder(resp.Body).Decode(userList); err != nil {
		return nil, err
	}

	return userList, nil
}

type UserList struct {
	Links struct {
		Self string `json:"self"`
	} `json:"links"`
	Users []struct {
		Id      string `json:"id"`
		OauthID string `json:"oauthID"`
		Name    string `json:"name"`
		Status  string `json:"status"`
		Links   struct {
			Self string `json:"self"`
			Logs string `json:"logs"`
		} `json:"links"`
	} `json:"users"`
}
