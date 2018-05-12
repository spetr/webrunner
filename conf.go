package main

import (
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

func parseConfig() {
	configFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatal("Configuration: Can not read configuration")
	}
	err = yaml.Unmarshal([]byte(configFile), &conf)
	if err != nil {
		log.Fatalf("Configuration: Parsing error: %v", err)
	}
}
