package config

import (
	"errors"
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

var configs map[string]Config

func Read(dir string) error {
	configs = make(map[string]Config)

	files, err := os.ReadDir(dir)
	if err != nil {
		error_message := fmt.Sprintf("Error reading config directory: %s", err)
		log.Println(error_message)
		return errors.New(error_message)
	}

	for _, file := range files {
		if !file.Type().IsRegular() {
			continue
		}

		content, err := os.ReadFile(dir + file.Name())
		if err != nil {
			error_message := fmt.Sprintf("Error reading config file %s: %s", file.Name(), err)
			log.Println(error_message)
			return errors.New(error_message)
		}

		cfg := Config{}
		err = yaml.Unmarshal(content, &cfg)
		if err != nil {
			error_message := fmt.Sprintf("Error yaml-parsing file %s: %s", file.Name(), err)
			log.Println(error_message)
			return errors.New(error_message)
		}

		configs[file.Name()] = cfg
	}

	return nil
}

func CheckSrcFolder(dir string) error {
	_, err := os.ReadDir(dir)
	if err != nil {
		error_message := fmt.Sprintf("Error reading src directory: %s", err)
		log.Println(error_message)
		return errors.New(error_message)
	}

	for key, _ := range configs {
		if _, err := os.Stat(dir + key); os.IsNotExist(err) {
			error_message := fmt.Sprintf("Error reading Cisco-config from %s%s", dir, key)
			log.Println(error_message)
			return errors.New(error_message)
		}
	}

	return nil
}

func Get() map[string]Config {
	return configs
}
