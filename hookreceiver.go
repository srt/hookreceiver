package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

type HookReceiveServer struct {
	Parse Parse
}

func handleNotification(notification Notification) {
	if repositoryConfig, found := config.FindRepositoryConfig(notification); found {
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

func (s HookReceiveServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	notification, err := s.Parse(r)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		fmt.Fprintf(w, "Unable to parse request: %s\n", err)
		log.Printf("Unable to parse request: %s\n", err)
		return
	}

	repo := notification.RepositoryUrl()
	branches := notification.Branches()

	fmt.Fprintf(w, "Ok, thanks for the notification about repository %q branches %v\n", repo, branches)
	log.Printf("Received notification for repository %q branches %v", repo, branches)

	go handleNotification(notification)
}

var configPath string
var config Config

func init() {
	flag.StringVar(&configPath, "c", "/etc/hookreceiver.conf.d", "Config path (file or directory)")
}

func reloadConfig(c <-chan os.Signal) {
	for s := range c {
		log.Printf("Got %s signal: Reloading configuration", s)
		newConfig, err := ReadConfig(configPath)
		if err == nil {
			config = newConfig
		} else {
			log.Println(err)
		}
	}
}

func main() {
	flag.Parse()
	os.Exit(run())
}

func run() int {
	var err error

	config, err = ReadConfig(configPath)
	if err != nil {
		log.Fatal(err)
		return 1
	}

	http.Handle("/hooks/bitbucket/", HookReceiveServer{BitbucketParse})

	server := &http.Server{Addr: config.Addr}
	listener, err := net.Listen("tcp", server.Addr)
	if err != nil {
		log.Fatal(err)
		return 1
	}

	go server.Serve(listener)
	// if err := server.Serve(listener); err != nil {
	// 	log.Fatal(err)
	// 	return 1
	// }

	log.Printf("HTTP server started")

	hupChannel := make(chan os.Signal, 1)
	signal.Notify(hupChannel, syscall.SIGHUP)
	go reloadConfig(hupChannel)

	killChannel := make(chan os.Signal, 1)
	signal.Notify(killChannel, os.Kill, os.Interrupt)

	<-killChannel

	return 0
}
