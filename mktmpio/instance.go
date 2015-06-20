package mktmpio

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Instance represents a server that has been created on the mktmpio service.
type Instance struct {
	ID          string `json:"id,omitempty"`
	Host        string `json:"host,omitempty"`
	Port        int    `json:"port,omitempty"`
	Error       string `json:"error,omitempty"`
	RemoteShell shell  `json:"remoteShell,omitempty"`
	Type        string `json:"type"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	client      Client
}

type shell struct {
	Cmd []interface{} `json:"cmd,omitempty"`
	Env []interface{} `json:"env,omitempty"`
}

// Destroy the server on the mktmpio service
func (i *Instance) Destroy() {
	i.client.Destroy(i.ID)
}

// Cmd returns an exec.Cmd that is pre-populated with the command, arguments,
// and environment variables required for spawning a local shell connected to
// the remote server.
func (i *Instance) Cmd() *exec.Cmd {
	args := make([]string, len(i.RemoteShell.Cmd))
	for i, a := range i.RemoteShell.Cmd {
		args[i] = fmt.Sprintf("%v", a)
	}
	return exec.Command(args[0], args[1:]...)
}

// LoadEnv modifies the current environment by setting environment variables
// that contain the host, port and credentials required for connecting to the
// remote server represented by the Instance.
func (i *Instance) LoadEnv() error {
	err := os.Setenv(envKey(i, "host"), i.Host)
	if err != nil {
		return err
	}
	err = os.Setenv(envKey(i, "port"), fmt.Sprintf("%d", i.Port))
	if err != nil {
		return err
	}
	err = os.Setenv(envKey(i, "username"), i.Username)
	if err != nil {
		return err
	}
	err = os.Setenv(envKey(i, "password"), i.Password)
	if err != nil {
		return err
	}
	return err
}

func envKey(i *Instance, field interface{}) string {
	return strings.ToUpper(fmt.Sprintf("%s_%s", i.Type, field))
}
