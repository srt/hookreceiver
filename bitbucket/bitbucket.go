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

func (n *Notification) Parse(r *http.Request) error {
	r.ParseForm()
	payload := r.Form.Get("payload")
	bytes := []byte(payload)

	return n.parseBytes(bytes)
}

func (n *Notification) parseBytes(bytes []byte) error {
	return json.Unmarshal(bytes, n)
}
