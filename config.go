// Copyright 2015 Datajin Technologies, Inc. All rights reserved.
// Use of this source code is governed by an Artistic-2
// license that can be found in the LICENSE file.

package mktmpio

import (
	"io/ioutil"
	"os"

	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"
)

// MktmpioURL is the root of the current version of the mktmpio HTTP API
const MktmpioURL = "https://mktmp.io/api/v1"

// MKtmpioCfgFile is the path to the user's config
const MKtmpioCfgFile = "~/.mktmpio.yml"

// Config contains the user config options used for accessing the mktmpio API.
type Config struct {
	Token string
	URL   string `yaml:",omitempty"`
	err   error
}

func (c Config) String() string {
	bytes, _ := yaml.Marshal(c)
	return string(bytes)
}

// LoadConfig loads the configuration stored in `~/.mktmpio.yml`, returning it
// as a Config type instance.
func LoadConfig() *Config {
	config := Config{}
	defConf := DefaultConfig()
	file := FileConfig(ConfigPath())
	env := EnvConfig()
	return config.Apply(defConf).Apply(file).Apply(env)
}

// DefaultConfig returns a configuration with only the default values set
func DefaultConfig() *Config {
	config := new(Config)
	config.URL = MktmpioURL
	return config
}

// EnvConfig returns a configuration with only values provided by environment variables
func EnvConfig() *Config {
	config := new(Config)
	config.Token = os.Getenv("MKTMPIO_TOKEN")
	config.URL = os.Getenv("MKTMPIO_URL")
	return config
}

// FileConfig returns a configuration with any values provided by the given YAML config file
func FileConfig(cfgPath string) *Config {
	config := new(Config)
	cfgFile, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		config.err = err
	} else {
		config.err = yaml.Unmarshal(cfgFile, config)
	}
	return config
}

// Apply creates a new Config with non-empty values from the provided Config
// overriding the options of the base Config
func (c *Config) Apply(b *Config) *Config {
	newCfg := new(Config)
	if b.Token == "" {
		newCfg.Token = c.Token
	} else {
		newCfg.Token = b.Token
	}
	if b.URL == "" {
		newCfg.URL = c.URL
	} else {
		newCfg.URL = b.URL
	}
	return newCfg
}

// ConfigPath returns the path to the user config file
func ConfigPath() string {
	if path, err := homedir.Expand(MKtmpioCfgFile); err == nil {
		return path
	}
	return ""
}

// Save stores the given configuration in ~/.mktmpio.yml, overwriting the
// current contents if the file exists.
func (c *Config) Save(cfgPath string) error {
	cfgFile, err := yaml.Marshal(c)
	if err == nil {
		err = ioutil.WriteFile(cfgPath, cfgFile, 0600)
	}
	return err
}
