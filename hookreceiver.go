package main

import (
	"bitbucket"
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
)

type Server struct {
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r)

	var notificaton bitbucket.Notification

	if err := notificaton.Parse(r); err != nil {
		log.Println(err)
	}

	log.Println(notificaton)
	fmt.Fprintln(w, notificaton)
	fmt.Fprintln(w, notificaton.Repository.Name)

	cmd := exec.Command("ls", "-al", notificaton.Repository.Name)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Printf("Command exited with error: %s\n", err)
	}
	fmt.Printf("Result: %q\n", out.String())
}

func main() {
	// Call realMain instead of doing the work here so we can use
	// `defer` statements within the function and have them work properly.
	// (defers aren't called with os.Exit)
	os.Exit(realMain())
}

// realMain is executed from main and returns the exit status to exit with.
func realMain() int {
	server := Server{}

	http.Handle("/", server)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
		return 1
	}
	return 0
}
