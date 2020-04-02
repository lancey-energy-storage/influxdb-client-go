package influxdb

import (
	"encoding/json"
	"log"
	"net/http"
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
