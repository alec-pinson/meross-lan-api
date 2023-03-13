package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Debug   bool
	File    string
	Devices []Device `yaml:"devices"`
	Key     string
}

type Device struct {
	Name    string `yaml:"name"`
	IP      string `yaml:"ip"`
	Channel int    `yaml:"channel"`
	Status  string
}

func (config Config) Load() Config {
	log.Println("Loading configuration...")

	// debug
	if strings.ToLower(os.Getenv("DEBUG")) == "true" {
		config.Debug = true
	} else {
		config.Debug = false
	}

	// config file path
	config.File = os.Getenv("CONFIG_FILE")
	if config.File == "" {
		config.File = "configuration.yaml"
	}

	// secret key
	config.Key = os.Getenv("KEY")
	if config.Key == "" {
		log.Fatal("Environment variable 'KEY' must be specified.")
	}

	// load yaml from file
	yamlFile, err := ioutil.ReadFile(config.File)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Configuration loaded...")

	return config
}
