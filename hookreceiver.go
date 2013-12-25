package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
)

type Server struct {
	Parse  Parse
	Config Config
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	notification, err := s.Parse(r)
	if err != nil {
		log.Printf("Unable to parse request: %s\n", err)
		return
	}

	repo := notification.RepositoryUrl()
	branches := notification.Branches()

	fmt.Fprintf(w, "Received notification for repository %q branches %q", repo, branches)
	log.Printf("Received notification for repository %q branches %q", repo, branches)

	if repositoryConfig, found := s.Config.FindRepositoryConfig(notification); found {
		log.Printf("Executing command %q in %q", repositoryConfig.Command, repositoryConfig.Dir)
		cmd := exec.Command("/bin/sh", "-c", repositoryConfig.Command)
		cmd.Dir = repositoryConfig.Dir
		out, err := cmd.CombinedOutput()
		log.Printf("Command output: %q", string(out))
		if err != nil {
			log.Printf("Command exited with error: %s\n", err)
		}
	} else {
		log.Printf("Repo/branch not configured.")
	}

}

var configFileName string

func init() {
	flag.StringVar(&configFileName, "c", "/etc/hookreceiver.conf.d", "Config file or directory name")
}

func main() {
	flag.Parse()
	os.Exit(run())
}

func run() int {
	config, err := readConfig(configFileName)
	if err != nil {
		log.Fatal(err)
		return 1
	}

	http.Handle("/hooks/bitbucket/", Server{BitbucketParse, config})
	if err := http.ListenAndServe(config.Addr, nil); err != nil {
		log.Fatal(err)
		return 1
	}

	return 0
}
