package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/srt/hookreceiver/bitbucket"
	"log"
	"net/http"
	"os"
	"os/exec"
)

type Server struct {
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var notificaton bitbucket.Notification

	if err := notificaton.Parse(r); err != nil {
		log.Printf("Unable to parse request: %s\n", err)
		return
	}

	repo := notificaton.RepositoryUrl()

	fmt.Fprintln(w, repo)
	log.Printf("Received notification for repository %q", repo)

	if repositoryConfig, found := config.Repositories[repo]; found {
		log.Printf("Executing command %q", repositoryConfig.Command)
		cmd := exec.Command("/bin/sh", "-c", repositoryConfig.Command)
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			log.Printf("Command exited with error: %s\n", err)
		} else {
			log.Printf("Command result: %q", out.String())
		}
	} else {
		log.Printf("Repo not configured")
	}

}

var configFileName string
var config Config

func init() {
	flag.StringVar(&configFileName, "c", "hookreceiver.conf", "Config file name")
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
	server := Server{}

	var err error
	config, err = readConfig(configFileName)
	if err != nil {
		log.Fatal(err)
		return 1
	}

	http.Handle("/", server)
	if err := http.ListenAndServe(config.Addr, nil); err != nil {
		log.Fatal(err)
		return 1
	}

	return 0
}
