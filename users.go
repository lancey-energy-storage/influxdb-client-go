package influxdb

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
)

func (c *Client) GetAllUsers() (*UserList, error) {
	log.Printf("[DEBUG] Get all users")

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

func (c *Client) CreateUser(name string, oauthID string, status string) (*User, error) {
	if name == "" {
		return nil, errors.New("a name is requried")
	}

	log.Printf("[DEBUG] Creation of a user with name %s", name)

	inputData, err := json.Marshal(NewUser{
		Name:    name,
		OauthID: oauthID,
		Status:  status,
	})
	req, err := http.NewRequest(http.MethodPost, c.url.String()+"/users", bytes.NewBuffer(inputData))
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

	user := &User{}
	if err := json.NewDecoder(resp.Body).Decode(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (c *Client) GetUserById(userID string) (*User, error) {
	if userID == "" {
		return nil, errors.New("a user id is required")
	}

	log.Printf("[DEBUG] Get user informations with id %s", userID)

	req, err := http.NewRequest(http.MethodGet, c.url.String()+"/users/"+userID, nil)
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

	user := &User{}
	if err := json.NewDecoder(resp.Body).Decode(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (c *Client) UpdateUser(userID string, name string, oauthID string, status string) (*User, error) {
	if userID == "" {
		return nil, errors.New("a user id is required")
	}
	if name == "" {
		return nil, errors.New("a name is required")
	}
	log.Printf("[DEBUG] Updating user informations with id %s", userID)

	inputData, err := json.Marshal(NewUser{
		Name:    name,
		OauthID: oauthID,
		Status:  status,
	})

	req, err := http.NewRequest(http.MethodPatch, c.url.String()+"/users/"+userID, bytes.NewBuffer(inputData))
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

	user := &User{}
	if err := json.NewDecoder(resp.Body).Decode(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (c *Client) DeleteUser(userID string) error {
	if userID == "" {
		return errors.New("a user id is required")
	}

	log.Printf("[DEBUG] Deleting the user with id %s", userID)

	req, err := http.NewRequest(http.MethodDelete, c.url.String()+"/users/"+userID, nil)
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

func (c *Client) UpdatePasswordOfUser(userID string, password string) error {
	if userID == "" {
		return errors.New("a user id is required")
	}
	if password == "" {
		return errors.New("a password is required")
	}

	log.Printf("[DEBUG] updating password for user with id %s", userID)

	inputData, err := json.Marshal(Password{
		Password: password,
	})
	req, err := http.NewRequest(http.MethodPut, c.url.String()+"/users/"+userID+"/password", bytes.NewBuffer(inputData))
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

func (c *Client) GetLogsOfUser(userID string, limit int, offset int) (*UserLogs, error) {
	if userID == "" {
		return nil, errors.New("a user id is required")
	}
	if limit == 0 || limit > 100 {
		return nil, errors.New("limit needs to be between [ 1 ... 100 ]")
	}
	if offset < 0 {
		return nil, errors.New("offset needs to be granter or equal to 0")
	}

	log.Printf("[DEBUG] Get logs of user with id: %s", userID)

	params := "?limit=" + strconv.Itoa(limit) + "&offset=" + strconv.Itoa(offset)
	req, err := http.NewRequest(http.MethodGet, c.url.String()+"/users/"+userID+"/logs"+params, nil)
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

	userLogs := &UserLogs{}
	if err := json.NewDecoder(resp.Body).Decode(userLogs); err != nil {
		return nil, err
	}

	return userLogs, nil
}

type UserList struct {
	Links struct {
		Self string `json:"self"`
	} `json:"links"`
	Users []User `json:"users"`
}

type User struct {
	Id      string `json:"id"`
	OauthID string `json:"oauthID"`
	Name    string `json:"name"`
	Status  string `json:"status"`
	Links   struct {
		Self string `json:"self"`
		Logs string `json:"logs"`
	} `json:"links"`
}

type NewUser struct {
	Name    string `json:"name"`
	OauthID string `json:"oauthID"`
	Status  string `json:"status"`
}

type Password struct {
	Password string `json:"password"`
}

type UserLogs struct {
	Logs []struct {
		Description string `json:"description"`
		Time        string `json:"time"`
		UserID      string `json:"userID"`
		Links       struct {
			User string `json:"user"`
		} `json:"links"`
	} `json:"logs"`
	Links struct {
		Next string `json:"next"`
		Self string `json:"self"`
		Prev string `json:"prev"`
	} `json:"links"`
}
