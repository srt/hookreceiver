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
	Parse Parse
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	notification, err := s.Parse(r)
	if err != nil {
		log.Printf("Unable to parse request: %s\n", err)
		return
	}

	repo := notification.RepositoryUrl()

	fmt.Fprintln(w, repo)
	log.Printf("Received notification for repository %q", repo)

	if repositoryConfig, found := config.FindRepositoryConfig(notification); found {
		log.Printf("Executing command %q", repositoryConfig.Command)
		cmd := exec.Command("/bin/sh", "-c", repositoryConfig.Command)
		cmd.Dir = repositoryConfig.Dir
		out, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("Command exited with error: %s\n", err)
		} else {
			log.Printf("Command result: %q", string(out))
		}
	} else {
		log.Printf("Repo/branch not configured.")
	}

}

var configFileName string
var config Config

func init() {
	flag.StringVar(&configFileName, "c", "/etc/hookreceiver.conf.d", "Config file or directory name")
}

func main() {
	// Call realMain instead of doing the work here so we can use
	// `defer` statements within the function and have them work properly.
	// (defers aren't called with os.Exit)

	flag.Parse()
	os.Exit(realMain())
}

// realMain is executed from main and returns the exit status to exit with.
func realMain() int {
	var err error
	config, err = readConfig(configFileName)
	if err != nil {
		log.Fatal(err)
		return 1
	}

	http.Handle("/hooks/bitbucket", Server{Parse: BitbucketParse})
	if err := http.ListenAndServe(config.Addr, nil); err != nil {
		log.Fatal(err)
		return 1
	}

	return 0
}
