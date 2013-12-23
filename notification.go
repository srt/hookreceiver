package main

type Notification interface {
	RepositoryUrl() string
	Branches() map[string]bool
}
