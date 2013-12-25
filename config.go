package main

import (
	"encoding/json"
	"fmt"
	"github.com/DisposaBoy/JsonConfigReader"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	Addr         string
	Repositories []RepositoryConfig
}

type RepositoryConfig struct {
	Url     string
	Branch  string
	Command string
	Dir     string
}

func (c *Config) FindRepositoryConfig(n Notification) (repositoryConfig RepositoryConfig, found bool) {
	for _, repositoryConfig = range c.Repositories {
		if repositoryConfig.Url != n.RepositoryUrl() {
			continue
		}
		if repositoryConfig.Branch != "" {
			if _, found = n.Branches()[repositoryConfig.Branch]; !found {
				continue
			}
		}
		found = true
		return
	}
	repositoryConfig = RepositoryConfig{}
	found = false
	return
}

func makeConfigPathVisitor(config *Config) func(path string, f os.FileInfo, err error) error {

	return func(path string, f os.FileInfo, err error) error {
		if f.Mode().IsRegular() {
			var currentConfig Config
			var file *os.File

			file, err = os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			reader := JsonConfigReader.New(file)
			err = json.NewDecoder(reader).Decode(&currentConfig)
			if err != nil {
				err = fmt.Errorf("Can't parse config file %q: %v", path, err)
				return err
			}

			if currentConfig.Addr != "" {
				config.Addr = currentConfig.Addr
			}
			config.Repositories = append(config.Repositories, currentConfig.Repositories...)
		}

		return nil
	}
}

func readConfig(filename string) (config Config, err error) {
	filepath.Walk(filename, makeConfigPathVisitor(&config))
	log.Printf("Config loaded from %q: %#v", filename, config)
	return
}
