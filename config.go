package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/DisposaBoy/JsonConfigReader"
)

type Config struct {
	Addr         string
	Secret       string
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

func appendConfig(config *Config, reader io.Reader) error {
	var currentConfig Config
	err := json.NewDecoder(reader).Decode(&currentConfig)
	if err != nil {
		return err
	}

	if currentConfig.Addr != "" {
		config.Addr = currentConfig.Addr
	}
	if currentConfig.Secret != "" {
		config.Secret = currentConfig.Secret
	}

	config.Repositories = append(config.Repositories, currentConfig.Repositories...)
	return nil
}

func makeConfigPathWalkFunc(config *Config) func(path string, f os.FileInfo, err error) error {
	return func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if f.Mode().IsRegular() {
			var file *os.File

			file, err = os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			err = appendConfig(config, JsonConfigReader.New(file))
			if err != nil {
				log.Printf("Can't parse config file %q: %v", path, err)
				return nil
			}
		}
		return nil
	}
}

func ReadConfig(filename string) (Config, error) {
	config := Config{}
	config.Addr = ":8080"
	err := filepath.Walk(filename, makeConfigPathWalkFunc(&config))
	if err != nil {
		return config, err
	}
	log.Printf("Config loaded from %q: %#v", filename, config)
	return config, nil
}
