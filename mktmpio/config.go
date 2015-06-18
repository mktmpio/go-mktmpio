package mktmpio

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os/user"
)

type Config struct {
	Token string `yaml:"token"`
}

func LoadConfig() (error, Config) {
	config := Config{}
	user, err := user.Current()
	if err != nil {
		return err, config
	}
	cfgFile, err := ioutil.ReadFile(user.HomeDir + "/.mktmpio.yml")
	if err != nil {
		return err, config
	}
	err = yaml.Unmarshal(cfgFile, &config)
	return err, config
}
