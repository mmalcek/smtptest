package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type tConfig struct {
	Server   string `yaml:"server"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	TLS      string `yaml:"TLS"`
	TLSvalid bool   `yaml:"TLSvalid"`
	TLSmin   string `yaml:"TLSmin"`
	TLSmax   string `yaml:"TLSmax"`
	Auth     string `yaml:"auth"`
	From     string `yaml:"from"`
	To       string `yaml:"to"`
	Subject  string `yaml:"subject"`
	Body     string `yaml:"body"`
}

var config tConfig

func loadConfig() error {
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(yamlFile, &config); err != nil {
		return err
	}
	return nil
}
