package mktmpio

import (
	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Token string `yaml:"token"`
}

func LoadConfig() (error, Config) {
	config := Config{}
	cfgPath, err := homedir.Expand("~/.mktmpio.yml")
	if err != nil {
		return err, config
	}
	cfgFile, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		return err, config
	}
	err = yaml.Unmarshal(cfgFile, &config)
	return err, config
}
