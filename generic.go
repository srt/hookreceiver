package main

import (
	"net/http"
)

type GenericNotification struct {
}

func GenericParse(r *http.Request) (n Notification, err error) {
	return GenericNotification{}, nil
}

func (n GenericNotification) RepositoryUrl() (repositoryUrl string) {
	repositoryUrl = ""
	return
}

func (n GenericNotification) Branches() (branches map[string]bool) {
	branches = make(map[string]bool)
	return
}
