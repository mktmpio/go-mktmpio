package mktmpio

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Instance struct {
	Id          string `json:"id,omitempty"`
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

func (i *Instance) Destroy() {
	i.client.Destroy(i.Id)
}

func (i *Instance) Cmd() *exec.Cmd {
	args := make([]string, len(i.RemoteShell.Cmd))
	for i, a := range i.RemoteShell.Cmd {
		args[i] = fmt.Sprintf("%v", a)
	}
	return exec.Command(args[0], args[1:]...)
}

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
