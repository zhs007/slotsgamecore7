package simserv

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Config - configuration
type Config struct {
	GameCode    string `yaml:"gamecode"`
	BindAddr    string `yaml:"bindaddr"`
	IsDebugMode bool   `yaml:"isdebugmode"`
	LogLevel    string `yaml:"loglevel"`
}

// LoadConfig - load configuration
func LoadConfig(fn string) (*Config, error) {
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
