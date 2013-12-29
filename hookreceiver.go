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
	"time"
)

const (
	BUFFER_SIZE = 20
	TIMEOUT     = 500 * time.Millisecond
)

type HookReceiveServer struct {
	Parse               Parse
	NotificationChannel chan Notification
}

func handleNotifications(ch <-chan Notification) {
	for {
		select {
		case notification := <-ch:
			handleNotification(notification)
		}
	}
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

	select {
	case s.NotificationChannel <- notification:
		fmt.Fprintf(w, "Ok, thanks for the notification about repository %q branches %v\n", repo, branches)
		log.Printf("Received and dispatched notification for repository %q branches %v", repo, branches)
	case <-time.After(TIMEOUT):
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Sorry, can't handle this notifaction right now (too many notifications pending)\n")
		log.Printf("Received but discarded notification for repository %q branches %v (too many notifications pending)", repo, branches)
	}
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

	notificationChannel := make(chan Notification, BUFFER_SIZE)
	go handleNotifications(notificationChannel)
	http.Handle("/hooks/bitbucket/", HookReceiveServer{BitbucketParse, notificationChannel})

	server := &http.Server{Addr: config.Addr}
	listener, err := net.Listen("tcp", server.Addr)
	if err != nil {
		log.Fatal(err)
		return 1
	}

	go server.Serve(listener)

	log.Printf("HTTP server started")

	hupChannel := make(chan os.Signal, 1)
	signal.Notify(hupChannel, syscall.SIGHUP)
	go reloadConfig(hupChannel)

	killChannel := make(chan os.Signal, 1)
	signal.Notify(killChannel, os.Kill, os.Interrupt)

	<-killChannel
	log.Println("Exiting")
	listener.Close()
	// TODO: terminate handleNotifications()

	return 0
}
