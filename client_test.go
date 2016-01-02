// Copyright 2015 Datajin Technologies, Inc. All rights reserved.
// Use of this source code is governed by an Artistic-2
// license that can be found in the LICENSE file.

package mktmpio

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	testConfig   = &Config{Token: "1234-5678-90abcdef"}
	badURLConfig = &Config{URL: "https://bad-host/api/invalid-encoding:%b%a%d"}
)

func server(t *testing.T, status int, body string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("User-Agent") != "go-mktmpio" {
			t.Error("default UserAgent not used")
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		w.Write([]byte(body))
	}))
}

func TestNewClient(t *testing.T) {
	client, err := NewClient(testConfig)
	if err != nil {
		t.Error("NewClient returned an error:", err)
	}
	if client == nil {
		t.Error("NewClient returned a nil client")
	}
	if len(client.token) < 10 {
		t.Error("client.token too short:", client.token)
	}
	if len(client.url) < 10 {
		t.Error("client.url too short:", client.url)
	}
	if client.UserAgent != "go-mktmpio" {
		t.Error("client.UserAgent is not default:", client.UserAgent)
	}
}

func TestClientRequest(t *testing.T) {
	client, _ := NewClient(badURLConfig)
	req, err := client.newRequest("", "")
	if err == nil || req != nil {
		t.Error("client.newRequest should error when client has bad url", err, req)
	}
}

func TestClientCreate(t *testing.T) {
	client, err := NewClient(testConfig)
	if err != nil {
		t.Error("NewClient returned an error")
	}
	ts := server(t, 400, `{"error": "unsupported type"}`)
	defer ts.Close()
	client.url = ts.URL
	instance, err := client.Create("definitely unsupported")
	if err == nil {
		t.Error("client.Create did not return an error")
	}
	if instance != nil {
		t.Error("client.Create returned an instance:", instance)
	}
}

func TestClientOptions(t *testing.T) {
	client, err := NewClient(testConfig)
	if err != nil {
		t.Error("NewClient returned an error")
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("User-Agent") != "my custom user agent" {
			t.Error("custome UserAgent not used")
		}
	}))
	defer ts.Close()
	client.url = ts.URL
	client.UserAgent = "my custom user agent"
	client.Create("whatever")
}

func TestBadCredentialsClient(t *testing.T) {
	client, err := NewClient(testConfig)
	if err != nil {
		t.Error("NewClient returned an error")
	}
	if client == nil {
		t.Error("NewClient returned a nil client")
	}
	client.token = "this is a bad token"
	ts := server(t, 401, `{"error": "Authentication required"}`)
	defer ts.Close()
	client.url = ts.URL
	instance, err := client.Create("valid")
	if err == nil {
		t.Error("client.Create did not return an error")
	}
	if instance != nil {
		t.Error("client.Create returned an instance:", instance)
	}
}

func TestClientServerGone(t *testing.T) {
	ts := server(t, 200, `[]`)
	client, _ := NewClient(testConfig)
	client.url = ts.URL
	ts.Close()
	_, err := client.Create("doesn't matter")
	if err == nil {
		t.Error("client.Create did not return an error")
	}
}

func TestClientBadJSON(t *testing.T) {
	ts := server(t, 200, `{"omg this isn't even JSON!"}`)
	defer ts.Close()
	client, _ := NewClient(testConfig)
	client.url = ts.URL
	_, err := client.Create("valid")
	if err == nil {
		t.Error("client.Create did not return an error")
	}
}

func TestCreateDestroy(t *testing.T) {
	mockToken := "abcdefg"
	mockCreate := func(w http.ResponseWriter, r *http.Request) {
		body := []byte(`{
			"id": "12345678",
			"host": "1.2.3.4",
			"port": 12345,
			"remoteShell": {
				"cmd": ["dbshell", "-h", "1.2.3.4"],
				"env": {"DB_USER": "foo", "DB_PASS": "bar"}
			},
			"type": "db",
			"username": "foo",
			"password": "bar"
		}`)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(201)
		w.Write(body)
	}
	mockDestroy := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(204)
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if token := r.Header.Get("X-Auth-Token"); token != mockToken {
			t.Errorf("Invalid token '%s' in request: %s %s", token, r.Method, r.URL)
		}
		if r.Method == "POST" {
			if r.URL.Path != "/new/db" {
				t.Errorf("Create used wrong URL: %s", r.URL)
			}
			mockCreate(w, r)
		} else if r.Method == "DELETE" {
			if r.URL.Path != "/i/12345678" {
				t.Errorf("Create used wrong URL: %s", r.URL)
			}
			mockDestroy(w, r)
		} else {
			t.Fatal("Unexpected request", r)
		}
	}))
	defer ts.Close()
	client, _ := NewClient(testConfig)
	client.url = ts.URL
	client.token = mockToken
	instance, err := client.Create("db")
	if err != nil {
		t.Error("Error creating instance:", err)
	}
	if err = instance.Destroy(); err != nil {
		t.Error("Error destroying instance:", err)
	}
}

func TestAttach(t *testing.T) {
	cfg := LoadConfig()
	if cfg.Token == "" {
		t.Skipf("requires a real token to connect to the real service")
		return
	}
	msg := make([]byte, 64)
	client, err := NewClient(cfg)
	redis, err := client.Create("redis")
	if err != nil {
		t.Fatal("Error creating redis instance for attach test", err)
	}
	defer redis.Destroy()
	rw, err := client.Attach(redis.ID)
	if err != nil {
		t.Fatal("Error attaching to redis instance", err)
	}
	// <ESC>[6n - requests a report cursor position
	wsReadT(msg, rw, t)
	// <ESC>[0;0R - respond with location 0,0
	wsWriteT([]byte("\x1b[1;1R"), rw, t)
	// <ESC>[999C<ESC>[6n - move cursor 999, report position
	wsReadT(msg, rw, t)
	// <ESC>[0;0R - report we are at 39,12
	wsWriteT([]byte("\x1b[12;39R"), rw, t)
	wsReadT(msg, rw, t)
	wsWriteT([]byte("scan 0\r\n"), rw, t)
	wsReadT(msg, rw, t)
	wsWriteT([]byte("exit\r\n"), rw, t)
	wsReadT(msg, rw, t)
	if err := rw.Close(); err != nil {
		t.Error("Did not close cleanly", err)
	}
}

func TestAttachStdio(t *testing.T) {
	cfg := LoadConfig()
	if cfg.Token == "" {
		t.Skipf("requires a real token to connect to the real service")
		return
	}
	msg := make([]byte, 64)
	client, err := NewClient(cfg)
	redis, err := client.Create("redis")
	if err != nil {
		t.Fatal("Error creating redis instance for attach test", err)
	}
	defer redis.Destroy()
	stdin, stdout, _, err := client.AttachStdio(redis.ID)
	if err != nil {
		t.Fatal("Error attaching to redis instance", err)
	}
	wsWriteT([]byte("scan 0\n"), stdin, t)
	if err := stdin.Close(); err != nil {
		t.Error("Did not close cleanly", err)
	}
	wsReadT(msg, stdout, t)
}

func wsWriteT(b []byte, w io.Writer, t *testing.T) {
	if n, err := w.Write(b); err != nil {
		t.Error("Failed to write command to websocket", err)
	} else if n == 0 {
		t.Error("Nowthing was written to websocket")
	} else {
		if testing.Verbose() {
			println("wrote:", tty(b))
		}
		t.Log("Wrote:", tty(b))
	}
}

func wsReadT(msg []byte, r io.Reader, t *testing.T) {
	if n, err := r.Read(msg); err != nil {
		t.Error("Failed to read prompt from websocket", err)
	} else if n == 0 {
		t.Error("Nothing was readable from websocket")
	} else {
		if testing.Verbose() {
			println("read:", tty(msg[:n]))
		}
		t.Log("Read:", n, tty(msg[:n]))
	}
}

func tty(b []byte) string {
	s := ""
	for _, c := range b {
		if c == 0x1B {
			s += "<ESC>"
		} else if c == 0x0D {
			s += "<CR>"
		} else if c == 0x0A {
			s += "<LF>"
		} else if c > 31 && c < 127 {
			s += fmt.Sprintf("%c", c)
		} else {
			s += fmt.Sprintf("<%02X>", c)
		}
	}
	return s
}
