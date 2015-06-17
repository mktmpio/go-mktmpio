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
	user, _ := user.Current()
	cfgFile, err := ioutil.ReadFile(user.HomeDir + "/.mktmpio.yml")
	err = yaml.Unmarshal(cfgFile, &config)
	return err, config
}
