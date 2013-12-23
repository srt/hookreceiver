package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Addr         string
	Repositories map[string]RepositoryConfig
}

type RepositoryConfig struct {
	Branch  string
	Command string
	Dir     string
}

func (c *Config) findRepositoryConfig(n Notification) (repositoryConfig RepositoryConfig, found bool) {
	if repositoryConfig, found = c.Repositories[n.RepositoryUrl()]; !found {
		return
	}
	if repositoryConfig.Branch != "" {
		if _, found = n.Branches()[repositoryConfig.Branch]; !found {
			repositoryConfig = RepositoryConfig{}
			return
		}
	}
	return
}

func readConfig(filename string) (config Config, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()

	err = json.NewDecoder(f).Decode(&config)
	if err != nil {
		err = fmt.Errorf("Can't parse config file %q: %v", filename, err)
		return
	}

	return
}
