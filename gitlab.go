package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type GitlabNotification struct {
	Before            string
	After             string
	Ref               string
	UserID            int    `json:"User_Id"`
	UserName          string `json:"User_Name"`
	ProjectID         int    `json:"Project_Id"`
	Repository        GitlabRepository
	Commits           []GitlabCommit
	TotalCommitsCount int `json:"Total_Commits_Count"`
}

type GitlabRepository struct {
	Name        string
	URL         string
	Description string
	Homepage    string
}

type GitlabCommit struct {
	ID        string
	Message   string
	Timestamp string
	URL       string
	Author    GitlabAuthor
}

type GitlabAuthor struct {
	Name  string
	Email string
}

func GitlabParse(r *http.Request) (n Notification, err error) {
	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)

	log.Printf("Received body %s", body)

	return gitlabParseBytes(body)
}

func gitlabParseBytes(bytes []byte) (n GitlabNotification, err error) {
	err = json.Unmarshal(bytes, &n)
	return
}

func (n GitlabNotification) RepositoryURL() (repositoryURL string) {
	repositoryURL = n.Repository.Homepage
	if repositoryURL[len(repositoryURL)-1] == '/' {
		repositoryURL = repositoryURL[:len(repositoryURL)-1]
	}
	return
}

// not supported by Gitlab push webhook
func (n GitlabNotification) Branches() (branches map[string]bool) {
	branches = make(map[string]bool)
	return
}
