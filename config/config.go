package config

import (
	"errors"
	"fmt"
	"os"

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
	if _, err := os.Stat("./config.yml"); errors.Is(err, os.ErrNotExist) {
		return Config{}, fmt.Errorf(fmt.Sprintf("No config found! Tried to open \"%s\"", "./config.yml"))
	}

	f, err := os.Open("./config.yml")
	if err != nil {
		return Config{}, err
	}
	defer f.Close()

	var conf Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&conf)

	return conf, err
}
