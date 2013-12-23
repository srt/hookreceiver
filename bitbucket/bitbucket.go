package bitbucket

import (
	"encoding/json"
	"net/http"
)

// See https://confluence.atlassian.com/display/BITBUCKET/POST+hook+management
type Notification struct {
	Canon_url  string
	User       string
	Repository Repository
	Commits    []Commit
}

type Repository struct {
	Absolute_url string
	Fork         bool
	Is_private   bool
	Name         string
	Owner        string
	Scm          string
	Slug         string
	Website      string
}

type Commit struct {
	Author       string
	Branch       string
	Files        []File
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

type File struct {
	File string
	Type string
}

func Parse(r *http.Request) (n Notification, err error) {
	r.ParseForm()
	payload := r.Form.Get("payload")
	bytes := []byte(payload)

	return parseBytes(bytes)
}

func parseBytes(bytes []byte) (n Notification, err error) {
	err = json.Unmarshal(bytes, &n)
	return
}

func (n Notification) RepositoryUrl() (repositoryUrl string) {
	repositoryUrl = n.Canon_url + n.Repository.Absolute_url
	if repositoryUrl[len(repositoryUrl)-1] == '/' {
		repositoryUrl = repositoryUrl[:len(repositoryUrl)-1]
	}
	return
}

func (n Notification) Branches() (branches map[string]bool) {
	branches = make(map[string]bool)
	for _, commit := range n.Commits {
		branches[commit.Branch] = true
	}
	return
}
