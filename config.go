package mktmpio

import (
	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

// Root API url for the current version of the mktmpio HTTP API
const MktmpioURL = "https://mktmp.io/api/v1"

// Config contains the user config options used for accessing the mktmpio API.
type Config struct {
	Token string
	URL   string
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
	if os.Getenv("MKTMPIO_TOKEN") != "" {
		config.Token = os.Getenv("MKTMPIO_TOKEN")
	}
	if os.Getenv("MKTMPIO_URL") != "" {
		config.URL = os.Getenv("MKTMPIO_URL")
	}
	if config.URL == "" {
		config.URL = MktmpioURL
	}
	return config, err
}
