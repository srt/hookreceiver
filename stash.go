package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// See https://confluence.atlassian.com/display/STASH/POST+service+webhook+for+Stash
type StashNotification struct {
	Repository StashRepository
	RefChanges []StashRefChange
	Changesets interface{}
}

type StashRepository struct {
	ID            int
	Slug          string
	Name          string
	ScmID         string
	State         string
	StatusMessage string
	Forkable      bool
	Project       StashProject
	Public        bool
}

type StashProject struct {
	ID         int
	Key        string
	Name       string
	Public     bool
	Type       string
	IsPersonal bool
	Owner      StashUser
}

type StashRefChange struct {
	RefID    string
	FromHash string
	ToHash   string
	Type     string
}

type StashUser struct {
	Name         string
	EmailAddress string
	ID           int
	DisplayName  string
	Active       bool
	Slug         string
	Type         string
}

func StashParse(r *http.Request) (n Notification, err error) {
	defer r.Body.Close()
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	return stashParseBytes(bytes)
}

func stashParseBytes(bytes []byte) (n StashNotification, err error) {
	err = json.Unmarshal(bytes, &n)
	return
}

func (n StashNotification) RepositoryURL() (repositoryURL string) {
	repositoryURL = n.Repository.Project.Key + "/" + n.Repository.Name
	if repositoryURL[len(repositoryURL)-1] == '/' {
		repositoryURL = repositoryURL[:len(repositoryURL)-1]
	}
	return
}

// not supported
func (n StashNotification) Branches() (branches map[string]bool) {
	branches = make(map[string]bool)
	return
}
