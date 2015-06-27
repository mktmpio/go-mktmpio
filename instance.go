package mktmpio

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Instance represents a server that has been created on the mktmpio service.
type Instance struct {
	ID          string
	Host        string
	Port        int
	Error       string
	RemoteShell shell
	Type        string
	Username    string
	Password    string
	client      Client
}

type shell struct {
	Cmd []string
	Env map[string]string
}

// Destroy the server on the mktmpio service
func (i *Instance) Destroy() {
	i.client.Destroy(i.ID)
}

// Cmd returns an exec.Cmd that is pre-populated with the command, arguments,
// and environment variables required for spawning a local shell connected to
// the remote server.
func (i *Instance) Cmd() *exec.Cmd {
	cmd := exec.Command(i.RemoteShell.Cmd[0], i.RemoteShell.Cmd[1:]...)
	if len(i.RemoteShell.Env) > 0 {
		cmd.Env = append(os.Environ(), envList(i.RemoteShell.Env)...)
	}
	return cmd
}

// LoadEnv modifies the current environment by setting environment variables
// that contain the host, port and credentials required for connecting to the
// remote server represented by the Instance.
func (i *Instance) LoadEnv() error {
	var err error
	setEnv := func(key, val string) {
		if err == nil {
			err = os.Setenv(envKey(i, key), val)
		}
	}
	setEnv("host", i.Host)
	setEnv("port", strconv.Itoa(i.Port))
	setEnv("username", i.Username)
	setEnv("password", i.Password)
	return err
}

func envKey(i *Instance, field interface{}) string {
	return strings.ToUpper(fmt.Sprintf("%s_%s", i.Type, field))
}

func envList(kv map[string]string) []string {
	env := make([]string, len(kv))
	for k, v := range kv {
		env = append(env, k+"="+v)
	}
	return env
}
