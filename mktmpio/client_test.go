package mktmpio

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	err, client := NewClient()
	if err != nil {
		t.Fail()
	}
	if len(client.token) < 10 {
		t.Fail()
	}
	if len(client.url) < 10 {
		t.Fail()
	}
	t.Logf("config: %v\n", client)
}
