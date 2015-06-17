package mktmpio

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Client struct {
	token string
	url   string
}

func NewClient() (error, *Client) {
	err, cfg := LoadConfig()
	if err != nil {
		return err, nil
	}
	client := Client{
		token: cfg.Token,
		url:   "https://mktmp.io/",
	}
	return nil, &client
}

func (c Client) Create(service string) (error, *Instance) {
	instance := Instance{client: c}
	req, _ := http.NewRequest("POST", c.url+"api/v1/new/"+service,
		strings.NewReader(""))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Auth-Token", c.token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err, nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err, nil
	}
	err = json.Unmarshal(body, &instance)
	if err != nil {
		fmt.Printf("Error reading JSON %v, %s", err, body)
		return err, nil
	}
	if len(instance.Error) > 0 {
		return errors.New(instance.Error), nil
	}
	return nil, &instance
}

func (c Client) Destroy(id string) error {
	url := c.url + "api/v1/i/" + id
	req, _ := http.NewRequest("DELETE", url, strings.NewReader(""))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Auth-Token", c.token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return err
}
