// Package mktmpio provides easy access to the database servier provisioning
// services at https://mktmp.io/
package mktmpio

import (
	"encoding/json"
	"errors"
	"golang.org/x/net/websocket"
	"io"
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

// Root WS url for current version of the mktmpio HTTP API
const MktmpioWSURL = "wss://mktmp.io:8443/ws"

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
func (c Client) jsonRequest(method, path string, instance *Instance) error {
	reqURL := c.url + path
	req, _ := http.NewRequest(method, reqURL, nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Auth-Token", c.token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if instance != nil {
		return json.Unmarshal(body, instance)
	}
	return nil
}

// Create creates a server of the type specified by `service`.
func (c Client) Create(service string) (*Instance, error) {
	instance := &Instance{client: c}
	reqURL := "/new/" + service
	if err := c.jsonRequest("POST", reqURL, instance); err != nil {
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
	return c.jsonRequest("DELETE", path, nil)
}

// Attach creates a remote shell for the instance identified by `id` and then
// returns a Reader and a Writer for communicating with it.
func (c Client) Attach(id string) (io.Reader, io.Writer, error) {
	url := MktmpioWSURL + "?id=" + id
	cfg, err := websocket.NewConfig(url, "http://localhost/")
	if err != nil {
		return nil, nil, err
	}
	cfg.Header.Set("Accept", "application/json")
	cfg.Header.Set("User-Agent", "mktmpio/cli")
	cfg.Header.Set("X-Auth-Token", c.token)
	ws, err := websocket.DialConfig(cfg)
	return ws, ws, err
}
