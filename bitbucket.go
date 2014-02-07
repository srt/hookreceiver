package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// See https://confluence.atlassian.com/display/BITBUCKET/POST+hook+management
type BitbucketNotification struct {
	Canon_url  string
	User       string
	Repository BitbucketRepository
	Truncated  bool
	Commits    []BitbucketCommit
}

type BitbucketRepository struct {
	Absolute_url string
	Fork         bool
	Is_private   bool
	Name         string
	Owner        string
	Scm          string
	Slug         string
	Website      string
}

type BitbucketCommit struct {
	Author       string
	Branches     []string
	Branch       string
	Files        []BitbucketFile
	Message      string
	Node         string
	Parents      []string
	Raw_author   string
	Raw_node     string
	Revision     int
	Size         int
	Timestamp    string
	Utctimestamp string
}

type BitbucketFile struct {
	File string
	Type string
}

func BitbucketParse(r *http.Request) (n Notification, err error) {
	r.ParseForm()
	payload := r.Form.Get("payload")

	log.Printf("Received payload %s", payload)

	bytes := []byte(payload)

	return parseBytes(bytes)
}

func parseBytes(bytes []byte) (n BitbucketNotification, err error) {
	err = json.Unmarshal(bytes, &n)
	return
}

func (n BitbucketNotification) RepositoryUrl() (repositoryUrl string) {
	repositoryUrl = n.Canon_url + n.Repository.Absolute_url
	if repositoryUrl[len(repositoryUrl)-1] == '/' {
		repositoryUrl = repositoryUrl[:len(repositoryUrl)-1]
	}
	return
}

func (n BitbucketNotification) Branches() (branches map[string]bool) {
	branches = make(map[string]bool)
	for _, commit := range n.Commits {
		if commit.Branch != "" {
			branches[commit.Branch] = true
		}
		for _, branch := range commit.Branches {
			branches[branch] = true
		}
	}
	return
}
