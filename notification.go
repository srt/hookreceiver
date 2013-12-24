package main

import (
	"net/http"
)

type Notification interface {
	RepositoryUrl() string
	Branches() map[string]bool
}

type Parse func(*http.Request) (notification Notification, err error)
