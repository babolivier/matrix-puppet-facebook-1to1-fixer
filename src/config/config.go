package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Config represents the full configuration.
type Config struct {
	Matrix MatrixConfig `yaml:"matrix"`
}

// MatrixConfig represents the Matrix part of the configuration.
type MatrixConfig struct {
	// HomeserverURL is the full URL of the homeserver to contact
	// (e.g. https://matrix.org/).
	HomeserverURL string `yaml:"homeserver_url"`
	// ServerName is the server name of the homeserver to contact
	// (e.g. matrix.org)0.
	ServerName string `yaml:"server_name"`
	// Localpart is the local part from the bridge user's Matrix ID (e.g. Alice).
	Localpart string `yaml:"localpart"`
	// AccessToken is the bridge user's access token.
	AccessToken string `yaml:"access_token"`
}

// Parse reads the file located at the provided path then proceeds to create and
// fill in a Config instance.
func Parse(path string) (cfg *Config, err error) {
	cfg = new(Config)

	// Load the content from the configuration file.
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	// Parse the YAML content from the configuration file.
	err = yaml.Unmarshal(content, cfg)
	return
}
