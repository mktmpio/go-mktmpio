// Copyright 2015 Datajin Technologies, Inc. All rights reserved.
// Use of this source code is governed by an Artistic-2
// license that can be found in the LICENSE file.

// Package mktmpio provides easy access to the database servier provisioning
// services at https://mktmp.io/
package mktmpio

import (
	"encoding/json"
	"errors"
	"github.com/docker/docker/pkg/stdcopy"
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

// NewRequest creates an http.Request based on the Client's configuration. The
// created request object is suitable for passing to http.Client.Do()
func (c Client) newRequest(method, path string) (*http.Request, error) {
	req, err := http.NewRequest(method, c.url+path, nil)
	if req != nil {
		req.Header.Set("Accept", "application/json")
		req.Header.Set("User-Agent", c.UserAgent)
		req.Header.Set("X-Auth-Token", c.token)
	}
	return req, err
}

// Create creates a server of the type specified by `service`.
func (c Client) rawRequest(method, path string) ([]byte, error) {
	req, _ := c.newRequest(method, path)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

// Create creates a server of the type specified by `service`.
func (c Client) jsonRequest(method, path string, instance interface{}) error {
	body, err := c.rawRequest(method, path)
	if err != nil || instance == nil {
		return err
	}
	return json.Unmarshal(body, instance)
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

// AttachStdio creates a remote shell for the instance identified by `id` and
// returns an io.WriteCloser for that shell's stdin and an io.Reader for each of
// stdout and stderr on that shell. This is for non-interactive shells, like one
// would use for piping a script into a shell or for piping the output from.
func (c Client) AttachStdio(id string) (io.WriteCloser, io.Reader, io.Reader, error) {
	conn, err := c.attachWS(id, true)
	if err != nil {
		return nil, nil, nil, err
	}
	inReader, inWriter := io.Pipe()
	outReader, outWriter := io.Pipe()
	errReader, errWriter := io.Pipe()
	go func() {
		// stdcopy is Docker's demuxer for their stdout/stderr multiplexed stream
		stdcopy.StdCopy(outWriter, errWriter, conn)
		errWriter.Close()
		outWriter.Close()
	}()
	go func() {
		io.Copy(conn, inReader)
		// A cheap hack sentinel value to indicate EOF to the server without closing
		// the actual connection. This would be so much easier with plain TCP :-(
		conn.Write([]byte{255, 255, 255, 255})
	}()
	return inWriter, outReader, errReader, err
}

// Attach creates a remote shell for the instance identified by `id` and then
// returns a Reader and a Writer for communicating with it via a psuedo-TTY. The
// bytes read from the channel will include TTY control sequences. This type of
// connection is most appropriate for connecting directly to a local TTY.
func (c Client) Attach(id string) (io.ReadWriteCloser, error) {
	return c.attachWS(id, false)
}

func (c Client) attachWS(id string, stdio bool) (*websocket.Conn, error) {
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
	params := url.Values{}
	params.Set("id", id)
	if stdio {
		params.Set("stdio", "true")
	} else {
		params.Set("stdio", "false")
	}
	wsURL.RawQuery = params.Encode()
	cfg, err := websocket.NewConfig(wsURL.String(), "http://localhost/")
	if err != nil {
		return nil, err
	}
	cfg.Header.Set("Accept", "application/json")
	cfg.Header.Set("User-Agent", "go-mktmpio")
	cfg.Header.Set("X-Auth-Token", c.token)
	conn, err := websocket.DialConfig(cfg)
	if err == nil {
		conn.PayloadType = websocket.BinaryFrame
	}
	return conn, err
}
