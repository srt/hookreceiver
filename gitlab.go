package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type GitlabNotification struct {
	Before              string
	After               string
	Ref                 string
	User_Id             int
	User_Name           string
	Project_Id          int
	Repository          GitlabRepository
	Commits             []GitlabCommit
	Total_Commits_Count int
}

type GitlabRepository struct {
	Name        string
	Url         string
	Description string
	Homepage    string
}

type GitlabCommit struct {
	Id        string
	Message   string
	Timestamp string
	Url       string
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

func (n GitlabNotification) RepositoryUrl() (repositoryUrl string) {
	repositoryUrl = n.Repository.Homepage
	if repositoryUrl[len(repositoryUrl)-1] == '/' {
		repositoryUrl = repositoryUrl[:len(repositoryUrl)-1]
	}
	return
}

// not supported by Gitlab push webhook
func (n GitlabNotification) Branches() (branches map[string]bool) {
	branches = make(map[string]bool)
	return
}
