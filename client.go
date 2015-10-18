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
	"net/url"
)

// go-1.2 does not automatically load SHA2-384, which is what parts of
// mktmp.io's cert chain use.
// See also: http://bridge.grumpy-troll.org/2014/05/golang-tls-comodo/
import _ "crypto/sha512" // side-effect only

// Client provides authenticated API access for creating, listing, and destorying
// database servers.
type Client struct {
	token     string
	url       string
	UserAgent string
}

// NewClient creates a mktmpio Client using credentials loaded from the user
// config stored in ~/.mktmpio.yml
func NewClient(cfg *Config) (*Client, error) {
	client := &Client{
		token:     cfg.Token,
		url:       cfg.URL,
		UserAgent: "go-mktmpio",
	}
	if client.url == "" {
		client.url = MktmpioURL
	}
	return client, nil
}

// Create creates a server of the type specified by `service`.
func (c Client) jsonRequest(method, path string, instance *Instance) error {
	reqURL := c.url + path
	req, _ := http.NewRequest(method, reqURL, nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)
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
func (c Client) Attach(id string) (io.ReadWriteCloser, error) {
	wsURL, err := url.Parse(c.url)
	if err != nil {
		return nil, err
	}
	if wsURL.Scheme == "https" {
		wsURL.Scheme = "wss"
		wsURL.Host = "mktmp.io:8443"
	} else {
		wsURL.Scheme = "ws"
	}
	wsURL.Path = "/ws"
	wsURL.RawQuery = "id=" + id
	cfg, err := websocket.NewConfig(wsURL.String(), "http://localhost/")
	if err != nil {
		return nil, err
	}
	cfg.Header.Set("Accept", "application/json")
	cfg.Header.Set("User-Agent", "go-mktmpio")
	cfg.Header.Set("X-Auth-Token", c.token)
	return websocket.DialConfig(cfg)
}
