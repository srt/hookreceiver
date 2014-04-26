package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// See https://confluence.atlassian.com/display/BITBUCKET/POST+hook+management
type BitbucketNotification struct {
	CanonURL   string `json:"Canon_url"`
	User       string
	Repository BitbucketRepository
	Truncated  bool
	Commits    []BitbucketCommit
}

type BitbucketRepository struct {
	AbsoluteURL string `json:"Absolute_url"`
	Fork        bool
	IsPrivate   bool `json:"Is_private"`
	Name        string
	Owner       string
	Scm         string
	Slug        string
	Website     string
}

type BitbucketCommit struct {
	Author       string
	Branches     []string
	Branch       string
	Files        []BitbucketFile
	Message      string
	Node         string
	Parents      []string
	RawAuthor    string `json:"Raw_author"`
	RawNode      string `json:"Raw_node"`
	Revision     int
	Size         int
	Timestamp    string
	UTCTimestamp string
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

	return bitbucketParseBytes(bytes)
}

func bitbucketParseBytes(bytes []byte) (n BitbucketNotification, err error) {
	err = json.Unmarshal(bytes, &n)
	return
}

func (n BitbucketNotification) RepositoryURL() (repositoryURL string) {
	repositoryURL = n.CanonURL + n.Repository.AbsoluteURL
	if repositoryURL[len(repositoryURL)-1] == '/' {
		repositoryURL = repositoryURL[:len(repositoryURL)-1]
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
