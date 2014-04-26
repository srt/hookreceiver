package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/DisposaBoy/JsonConfigReader"
)

type Config struct {
	Addr         string
	Repositories []RepositoryConfig
}

type RepositoryConfig struct {
	Name    string
	URL     string
	Branch  string
	Command string
	Dir     string
}

func (c *Config) FindRepositoryConfig(name string, n Notification) (RepositoryConfig, bool) {
	for _, repositoryConfig := range c.Repositories {
		if repositoryConfig.URL != "" && repositoryConfig.URL != n.RepositoryURL() {
			continue
		}
		if repositoryConfig.Name != "" && repositoryConfig.Name != name {
			continue
		}
		if repositoryConfig.Branch != "" {
			if _, found := n.Branches()[repositoryConfig.Branch]; !found {
				continue
			}
		}
		return repositoryConfig, true
	}
	return RepositoryConfig{}, false
}

func makeConfigPathWalkFunc(config *Config) func(path string, f os.FileInfo, err error) error {
	return func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if f.Mode().IsRegular() {
			var currentConfig Config
			var file *os.File

			file, err = os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			reader := JsonConfigReader.New(file)
			jsonErr := json.NewDecoder(reader).Decode(&currentConfig)
			if jsonErr != nil {
				log.Printf("Can't parse config file %q: %v", path, jsonErr)
				return nil
			}

			if currentConfig.Addr != "" {
				config.Addr = currentConfig.Addr
			}
			config.Repositories = append(config.Repositories, currentConfig.Repositories...)
		}
		return nil
	}
}

func ReadConfig(filename string) (Config, error) {
	config := Config{}
	err := filepath.Walk(filename, makeConfigPathWalkFunc(&config))
	if err != nil {
		return config, err
	}
	log.Printf("Config loaded from %q: %#v", filename, config)
	return config, nil
}
