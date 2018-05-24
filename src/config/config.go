package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Matrix MatrixConfig `yaml:"matrix"`
}

type MatrixConfig struct {
	HomeserverURL string `yaml:"homeserver_url"`
	ServerName    string `yaml:"server_name"`
	Localpart     string `yaml:"localpart"`
	AccessToken   string `yaml:"access_token"`
}

func Parse(path string) (cfg *Config, err error) {
	cfg = new(Config)

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(content, cfg)
	return
}
