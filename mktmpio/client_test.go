package mktmpio

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	err, client := NewClient()
	if err != nil {
		t.Error("NewClient returned an error:", err)
	}
	if len(client.token) < 10 {
		t.Error("client.token too short:", client.token)
	}
	if len(client.url) < 10 {
		t.Error("client.url too short:", client.url)
	}
}

func TestClientCreate(t *testing.T) {
	err, client := NewClient()
	if err != nil {
		t.Error("NewClient returned an error")
	}
	if client == nil {
		t.Error("NewClient returned a nil client")
	}
	err, instance := client.Create("definitely unsupported")
	if err == nil {
		t.Error("client.Create did not return an error")
	}
	if instance != nil {
		t.Error("client.Create returned an instance:", instance)
	}
}
