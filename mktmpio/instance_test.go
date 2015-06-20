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

func TestCmdNoEnv(t *testing.T) {
	instance := Instance{
		ID:       "someId",
		Host:     "some-host",
		Port:     1234,
		Type:     "mktmpdb",
		Username: "user",
		Password: "pass",
		RemoteShell: shell{
			Cmd: []interface{}{"tmpdbcli", "-h", "some-host", "-p", 1234},
		},
	}
	cmd := instance.Cmd()
	if cmd.Path != "tmpdbcli" {
		t.Error("cmd.Path incorrect:", cmd.Path)
	}
	if cmd.Args[0] != cmd.Path {
		t.Error("cmd.Path and cmd.Args[0] should be the same", cmd.Args[0])
	}
	if len(cmd.Args) != 5 {
		t.Error("cmd.Args wrong length:", len(cmd.Args))
	}
	if cmd.Args[4] != "1234" {
		t.Error("int argument was not stringified correctly:", cmd.Args[4])
	}
	if len(cmd.Env) != 0 {
		t.Error("no env variables should be set", cmd.Env)
	}
}

func TestCmdWithEnv(t *testing.T) {
	instance := Instance{
		ID:       "someId",
		Host:     "some-host",
		Port:     1234,
		Type:     "mktmpdb",
		Username: "user",
		Password: "pass",
		RemoteShell: shell{
			Cmd: []interface{}{"tmpdbcli", "-h", "some-host", "-p", 1234},
			Env: map[string]string{
				"MKTMPIO-DBPASS": "pass",
			},
		},
	}
	cmd := instance.Cmd()
	if len(cmd.Args) != 5 {
		t.Error("cmd.Args wrong length:", len(cmd.Args))
	}
	if len(cmd.Env) < 1 {
		t.Error("cmd.Env should be populated:", cmd.Env)
	}
	if os.Getenv("MKTMPIO-DBPASS") != "" {
		t.Error("real environment should not already contain MKTMPIO-DBPASS")
	}
	if cmd.Env[len(cmd.Env)-1] != "MKTMPIO-DBPASS=pass" {
		t.Error("required shell env var not set:", cmd.Env[len(cmd.Env)-1])
	}
}
