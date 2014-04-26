package main

import (
	"net/http"
)

type NotificationRequest struct {
	Path         string
	Notification Notification
}

type Notification interface {
	RepositoryURL() string
	Branches() map[string]bool
}

type Parse func(*http.Request) (notification Notification, err error)
