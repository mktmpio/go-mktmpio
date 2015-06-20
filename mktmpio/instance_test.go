package mktmpio

import (
	"os"
	"testing"
)

func TestLoadEnv(t *testing.T) {
	instance := Instance{
		ID:       "someId",
		Host:     "some-host",
		Port:     1234,
		Type:     "mktmpdb",
		Username: "user",
		Password: "pass",
	}
	err := instance.LoadEnv()
	if err != nil {
		t.Error("LoadEnv returned an error:", err)
	}
	val := os.Getenv("MKTMPDB_HOST")
	if val != "some-host" {
		t.Error("host env var not 'some-host':", val)
	}
	val = os.Getenv("MKTMPDB_PORT")
	if val != "1234" {
		t.Error("port env var not '1234':", val)
	}
	val = os.Getenv("MKTMPDB_USERNAME")
	if val != "user" {
		t.Error("user env var not 'user':", val)
	}
	val = os.Getenv("MKTMPDB_PASSWORD")
	if val != "pass" {
		t.Error("pssword env var not 'pass':", val)
	}
}
