package main

import (
	"encoding/json"
	"net/http"
)

// See https://confluence.atlassian.com/display/BITBUCKET/POST+hook+management
type BitbucketNotification struct {
	Canon_url  string
	User       string
	Repository struct {
		Absolute_url string
		Fork         bool
		Is_private   bool
		Name         string
		Owner        string
		Scm          string
		Slug         string
		Website      string
	}
	Commits []struct {
		Author string
		Branch string
		Files  []struct {
			File string
			Type string
		}
		Message      string
		Node         string
		Parents      []string
		Raw_author   string
		Raw_node     string
		Revision     string
		Size         int
		Timestamp    string
		Utctimestamp string
	}
}

func (n *BitbucketNotification) Parse(r *http.Request) error {
	r.ParseForm()
	payload := r.Form.Get("payload")
	bytes := []byte(payload)

	return json.Unmarshal(bytes, n)
}
