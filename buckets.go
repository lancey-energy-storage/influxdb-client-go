package influxdb

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func (c *Client) GetBucketsInSource(id string) (*BucketSource, error) {
	log.Printf("[DEBUG] Get buckets from source with id %s ", id)

	req, err := http.NewRequest(http.MethodGet, c.url.String()+"/sources/"+id+"/buckets", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Authorization", c.authorization)
	resp, err := c.httpClient.Do(req)

	defer resp.Body.Close()
	bucketSource := &BucketSource{}
	if err := json.NewDecoder(resp.Body).Decode(bucketSource); err != nil {
		return nil, err
	}

	return bucketSource, nil
}

func (c *Client) GetBuckets(limit int, name string, offset int, org string, orgID string) (*BucketSource, error) {
	if limit == 0 || limit > 100 {
		return nil, errors.New("limit needs to be between [ 1 ... 100 ]")
	}
	if offset < 0 {
		return nil, errors.New("offset needs to be granter or equal to 0")
	}
	log.Printf("[DEBUG] Get buckets")

	params := "/buckets?limit=" + strconv.Itoa(limit) + "&name=" + name + "&offset=" + strconv.Itoa(offset) + "&org=" + org + "&orgID=" + orgID + ""
	req, err := http.NewRequest(http.MethodGet, c.url.String()+params, nil)
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

	bucketSource := &BucketSource{}
	if err := json.NewDecoder(resp.Body).Decode(bucketSource); err != nil {
		return nil, err
	}

	return bucketSource, nil
}

func (c *Client) CreateBucket(description string, name string, orgID string, retentionRules []RetentionRules, rp string) (*BucketCreate, error) {
	if name == "" {
		return nil, errors.New("name is needed to create a new bucket")
	}
	if len(retentionRules) == 0 {
		return nil, errors.New("retentions rules is/are needed to create a new bucket")
	}
	log.Printf("[DEBUG] Creation of a new bucket")

	inputData, err := json.Marshal(SetupCreateBucket{
		Description:    description,
		Name:           name,
		OrgID:          orgID,
		RetentionRules: retentionRules,
		Rp:             rp,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.url.String()+"/buckets", bytes.NewBuffer(inputData))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", c.authorization)
	resp, err := c.httpClient.Do(req)

	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		return nil, errors.New(resp.Status)
	}

	bucketCreate := &BucketCreate{}
	if err := json.NewDecoder(resp.Body).Decode(bucketCreate); err != nil {
		return nil, err
	}

	return bucketCreate, nil
}

func (c *Client) GetBucketByID(bucketID string) (*SimpleBucket, error) {
	if bucketID == "" {
		return nil, errors.New("a bucket id is required")
	}

	log.Printf("[DEBUG] Get bucket with id: %s", bucketID)

	req, err := http.NewRequest(http.MethodGet, c.url.String()+"/buckets/"+bucketID, nil)
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

	simpleBucket := &SimpleBucket{}
	if err := json.NewDecoder(resp.Body).Decode(simpleBucket); err != nil {
		return nil, err
	}

	return simpleBucket, nil
}

func (c *Client) UpdateABucket(bucketID string, description string, labels []Labels, name string, orgID string, retentionRules []RetentionRules, rp string) (*SimpleBucket, error) {
	if name == "" {
		return nil, errors.New("name should be specified")
	}
	if len(retentionRules) == 0 {
		return nil, errors.New("retention rules should be specified")
	}

	log.Printf("[DEBUG] Updating the bucket with id: %s", bucketID)

	inputData, err := json.Marshal(SetupUpdateBucket{
		Description:    description,
		Labels:         labels,
		Name:           name,
		OrgID:          orgID,
		RetentionRules: retentionRules,
		Rp:             rp,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", c.url.String()+"/buckets/"+bucketID, bytes.NewBuffer(inputData))
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

	updateBucket := &SimpleBucket{}
	if err := json.NewDecoder(resp.Body).Decode(updateBucket); err != nil {
		return nil, err
	}

	return updateBucket, nil
}

func (c *Client) DeleteABucket(bucketID string) error {
	if bucketID == "" {
		return errors.New("a bucketID should be specified")
	}

	log.Printf("[DEBUG] Deleting bucket with id: %s", bucketID)

	req, err := http.NewRequest("DELETE", c.url.String()+"/buckets/"+bucketID, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", c.authorization)
	resp, err := c.httpClient.Do(req)

	defer resp.Body.Close()

	if resp.StatusCode != 204 {
		return errors.New(resp.Status)
	}

	return nil
}

func (c *Client) ListLabelsForABucket(bucketID string) (*LabelsOfBucket, error) {
	if bucketID == "" {
		return nil, errors.New("a bucket id is required")
	}

	log.Printf("[DEBUG] Getting labels for bucket with id: %s", bucketID)

	req, err := http.NewRequest("GET", c.url.String()+"/buckets/"+bucketID+"/labels", nil)
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

	labelsOfBucket := &LabelsOfBucket{}
	if err := json.NewDecoder(resp.Body).Decode(labelsOfBucket); err != nil {
		return nil, err
	}

	return labelsOfBucket, nil
}

func (c *Client) AddLabelToBucket(bucketID string, labelID string) (*LabelsOfBucket, error) {
	if bucketID == "" {
		return nil, errors.New("a bucket id is required to add label to it")
	}
	if len(labelID) == 0 {
		return nil, errors.New("an array of one label id is required")
	}

	log.Printf("[DEBUG] Adding label to the bucket with the id: %s", bucketID)

	inputData := fmt.Sprintf("{\"labelID\": \"%s\"}", labelID)

	req, err := http.NewRequest(http.MethodPost, c.url.String()+"/buckets/"+bucketID+"/labels", bytes.NewBuffer([]byte(inputData)))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", c.authorization)
	resp, err := c.httpClient.Do(req)

	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		return nil, errors.New(resp.Status)
	}

	labelsOfBucket := &LabelsOfBucket{}
	if err := json.NewDecoder(resp.Body).Decode(labelsOfBucket); err != nil {
		return nil, err
	}

	return labelsOfBucket, nil
}

func (c *Client) DeleteALabelFromBucket(bucketID string, labelID string) error {
	if bucketID == "" {
		return errors.New("a bucket id is required")
	}
	if labelID == "" {
		return errors.New("a label id is required")
	}

	log.Printf("[DEBUG] Deleting label id %s of bucket id %s", labelID, bucketID)

	req, err := http.NewRequest("DELETE", c.url.String()+"/buckets/"+bucketID+"/labels/"+labelID, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", c.authorization)
	resp, err := c.httpClient.Do(req)

	defer resp.Body.Close()

	if resp.StatusCode != 204 {
		return errors.New(resp.Status)
	}

	return nil
}

type BucketSource struct {
	Links struct {
		Next string `json:"next"`
		Self string `json:"self"`
		Prev string `json:"prev"`
	} `json:"links"`
	Buckets []struct {
		Links struct {
			Labels  string `json:"labels"`
			Logs    string `json:"logs"`
			Members string `json:"members"`
			Org     string `json:"org"`
			Owners  string `json:"owners"`
			Self    string `json:"self"`
			Write   string `json:"write"`
		} `json:"links"`
		Id             string `json:"id"`
		Type           string `json:"type"`
		Name           string `json:"name"`
		Description    string `json:"description"`
		OrgId          string `json:"orgId"`
		Rp             string `json:"rp"`
		CreatedAt      string `json:"createdAt"`
		UpdatedAt      string `json:"updatedAt"`
		RetentionRules []struct {
			Type         string `json:"type"`
			EverySeconds int    `json:"everySeconds"`
		} `json:"retentionRules"`
	} `json:"buckets"`
	Labels []struct {
		Id          string `json:"id"`
		OrgId       string `json:"orgId"`
		Name        string `json:"name"`
		Properties  string `json:"properties"`
		Color       string `json:"color"`
		Description string `json:"description"`
	} `json:"labels"`
}

type BucketCreate struct {
	Links struct {
		Labels  string `json:"labels"`
		Logs    string `json:"logs"`
		Members string `json:"members"`
		Org     string `json:"org"`
		Owners  string `json:"owners"`
		Self    string `json:"self"`
		Write   string `json:"write"`
	} `json:"links"`
	Id             string `json:"id"`
	Type           string `json:"user"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	OrgID          string `json:"orgID"`
	Rp             string `json:"rp"`
	CreatedAt      string `json:"createdAt"`
	UpdatedAt      string `json:"updatedAt"`
	RetentionRules []struct {
		Type         string `json:"type"`
		EverySeconds int    `json:"everySeconds"`
	} `json:"retentionRules"`
	Labels []struct {
		Id         string `json:"id"`
		OrgID      string `json:"orgID"`
		Name       string `json:"name"`
		Properties string `json:"properties"`
	} `json:"labels"`
}

type RetentionRules struct {
	EverySeconds int
	Type         string
}

type SetupCreateBucket struct {
	Description    string `json:"description"`
	Name           string `json:"name"`
	OrgID          string `json:"orgID"`
	RetentionRules []RetentionRules
	Rp             string `json:"rp"`
}

type SimpleBucket struct {
	Links struct {
		Labels  string `json:"labels"`
		Logs    string `json:"logs"`
		Members string `json:"members"`
		Org     string `json:"org"`
		Owners  string `json:"owners"`
		Self    string `json:"self"`
		Write   string `json:"write"`
	} `json:"links"`
	Id             string `json:"id"`
	Type           string `json:"user"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	OrgID          string `json:"orgID"`
	Rp             string `json:"rp"`
	CreatedAt      string `json:"createdAt"`
	UpdatedAt      string `json:"updatedAt"`
	RetentionRules []struct {
		Type         string `json:"type"`
		EverySeconds int    `json:"everySeconds"`
	} `json:"retentionRules"`
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

type Labels struct {
	Name       string
	Properties string
}

type SetupUpdateBucket struct {
	Description    string `json:"description"`
	Labels         []Labels
	Name           string `json:"name"`
	OrgID          string `json:"orgID"`
	RetentionRules []RetentionRules
	Rp             string `json:"rp"`
}

type LabelsOfBucket struct {
	Labels []struct {
		Id         string `json:"id"`
		OrgId      string `json:"orgID"`
		Name       string `json:"name"`
		Properties struct {
			Color       string `json:"color"`
			Description string `json:"description"`
		} `json:"properties"`
	} `json:"labels"`
	Links struct {
		Next string `json:"next"`
		Self string `json:"self"`
		Prev string `json:"prev"`
	} `json:"links"`
}
