package mktmpio

import (
	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// Config contains the user config options used for accessing the mktmpio API.
type Config struct {
	Token string `yaml:"token"`
}

// LoadConfig loads the configuration stored in `~/.mktmpio.yml`, returning it
// as a Config type instance.
func LoadConfig() (Config, error) {
	config := Config{}
	cfgPath, err := homedir.Expand("~/.mktmpio.yml")
	if err != nil {
		return config, err
	}
	cfgFile, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		return config, err
	}
	err = yaml.Unmarshal(cfgFile, &config)
	return config, err
}
