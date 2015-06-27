package mktmpio

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func server(t *testing.T, status int, body string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		w.Write([]byte(body))
	}))
}

func TestNewClient(t *testing.T) {
	client, err := NewClient()
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
}

func TestClientCreate(t *testing.T) {
	client, err := NewClient()
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

func TestBadCredentialsClient(t *testing.T) {
	client, err := NewClient()
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
	client := Client{
		url:   ts.URL,
		token: mockToken,
	}
	instance, err := client.Create("db")
	if err != nil {
		t.Error("Error creating instance:", err)
	}
	if err = instance.Destroy(); err != nil {
		t.Error("Error destroying instance:", err)
	}
}
