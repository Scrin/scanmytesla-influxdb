package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type InfluxDB struct {
	Url         string `yaml:"url"`
	AuthToken   string `yaml:"auth_token"`
	Org         string `yaml:"org"`
	Bucket      string `yaml:"bucket"`
	Measurement string `yaml:"measurement"`
}

type Config struct {
	InfluxDB InfluxDB `yaml:"influxdb,omitempty"`
}

func ReadConfig() (Config, error) {
	file := "./config.yml"
	if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
		ex, err := os.Executable()
		if err != nil {
			return Config{}, err
		}
		file = filepath.Dir(ex) + "/config.yml"
		if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
			return Config{}, fmt.Errorf(fmt.Sprintf(`No config found! Tried to open "./config.yml" or "%s"`, file))
		}
	}

	f, err := os.Open(file)
	if err != nil {
		return Config{}, err
	}
	defer f.Close()

	var conf Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&conf)

	return conf, err
}
