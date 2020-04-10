package influxdb

import (
	"bytes"
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

func (c *Client) CreateAuthorization(description string, orgID string, permissions []Permissions, status string) (*AuthorizationCreated, error) {
	if orgID == "" {
		return nil, errors.New("an org id is required")
	}
	if len(permissions) == 0 {
		return nil, errors.New("a list of permissions is required")
	}

	log.Printf("[DEBUG] Posting a new permission")

	inputData, err := json.Marshal(SetupNewAuthorization{
		Description: description,
		OrgID:       orgID,
		Permissions: permissions,
		Status:      status,
	})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, c.url.String()+"/authorizations", bytes.NewBuffer(inputData))
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

	authorizationCreated := &AuthorizationCreated{}
	if err := json.NewDecoder(resp.Body).Decode(authorizationCreated); err != nil {
		return nil, err
	}

	return authorizationCreated, nil
}

func (c *Client) GetAuthorizationById(authID string) (*AuthorizationDetails, error) {
	if authID == "" {
		return nil, errors.New("a auth id is required")
	}

	log.Printf("[DEBUG] Getting the authorization with id %s", authID)

	req, err := http.NewRequest(http.MethodGet, c.url.String()+"/authorizations/"+authID, nil)
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

	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	authorizationDetails := &AuthorizationDetails{}
	if err := json.NewDecoder(resp.Body).Decode(authorizationDetails); err != nil {
		return nil, err
	}

	return authorizationDetails, nil
}

func (c *Client) UpdateAnAuthorizationStatus(authID string, description string, status string) (*AuthorizationDetails, error) {
	if authID == "" {
		return nil, errors.New("an auth id is required")
	}

	log.Printf("[DEBUG] Updating the status of the authorization with id %s", authID)

	inputData, err := json.Marshal(Status{
		Description: description,
		Status:      status,
	})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPatch, c.url.String()+"/authorizations/"+authID, bytes.NewBuffer(inputData))
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

	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	authorizationDetails := &AuthorizationDetails{}
	if err := json.NewDecoder(resp.Body).Decode(authorizationDetails); err != nil {
		return nil, err
	}

	return authorizationDetails, nil
}

func (c *Client) DeleteAnAuthorization(authID string) error {
	if authID == "" {
		return errors.New("an auth id is required")
	}

	log.Printf("[DEBUG] Delete the authorization with id %s", authID)

	req, err := http.NewRequest(http.MethodDelete, c.url.String()+"/authorizations/"+authID, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", c.authorization)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 204 {
		return errors.New(resp.Status)
	}

	return nil
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

type Permissions struct {
	Action   string   `json:"action"`
	Resource Resource `json:"resource"`
}

type Resource struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Org   string `json:"org"`
	OrgID string `json:"orgID"`
	Type  string `json:"type"`
}

type Status struct {
	Description string `json:"description"`
	Status      string `json:"status"`
}

type AuthorizationCreated struct {
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
			OrgId string `json:"orgID"`
			Org   string `json:"org"`
		} `json:"resource"`
	} `json:"permissions"`
	Id     string `json:"id"`
	Token  string `json:"token"`
	UserID string `json:"userID"`
	User   string `json:"user"`
	Org    string `json:"user"`
	Links  struct {
		Self string `json:"self"`
		User string `json:"user"`
	} `json:"links"`
}

type SetupNewAuthorization struct {
	Description string        `json:"description"`
	OrgID       string        `json:"orgID"`
	Permissions []Permissions `json:"permissions"`
	Status      string        `json:"status"`
}

type AuthorizationDetails struct {
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
}
