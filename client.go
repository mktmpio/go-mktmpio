// Package mktmpio provides easy access to the database servier provisioning
// services at https://mktmp.io/
package mktmpio

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Client provides authenticated API access for creating, listing, and destorying
// database servers.
type Client struct {
	token string
	url   string
}

// Root API url for the current version of the mktmpio HTTP API
const MktmpioURL = "https://mktmp.io/api/v1"

// NewClient creates a mktmpio Client using credentials loaded from the user
// config stored in ~/.mktmpio.yml
func NewClient() (*Client, error) {
	cfg, err := LoadConfig()
	if err != nil {
		return nil, err
	}
	client := &Client{
		token: cfg.Token,
		url:   MktmpioURL,
	}
	return client, nil
}

// Create creates a server of the type specified by `service`.
func (c Client) jsonRequest(method, path string) ([]byte, error) {
	reqURL := c.url + path
	req, _ := http.NewRequest(method, reqURL, nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Auth-Token", c.token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

// Create creates a server of the type specified by `service`.
func (c Client) Create(service string) (*Instance, error) {
	instance := &Instance{client: c}
	reqURL := "/new/" + service
	body, err := c.jsonRequest("POST", reqURL)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(body, instance); err != nil {
		fmt.Printf("Error reading JSON %v, %s", err, body)
		return nil, err
	}
	if len(instance.Error) > 0 {
		return nil, errors.New(instance.Error)
	}
	return instance, nil
}

// Destroy shuts down and deletes the server identified by `id`.
func (c Client) Destroy(id string) error {
	path := "/i/" + id
	_, err := c.jsonRequest("DELETE", path)
	return err
}
