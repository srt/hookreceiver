package main

import (
	"testing"
)

var repositoryConfig1 = RepositoryConfig{
	Url:    "https://bitbucket.org/user/repo1",
	Branch: "master",
}

var repositoryConfig2 = RepositoryConfig{
	Url: "https://bitbucket.org/user/repo2",
}

type mockNotification struct {
	repositoryUrl string
	branches      []string
}

func (n mockNotification) RepositoryUrl() string {
	return n.repositoryUrl
}

func (n mockNotification) Branches() map[string]bool {
	branchMap := make(map[string]bool)
	for _, branch := range n.branches {
		branchMap[branch] = true
	}
	return branchMap
}

type findRepositoryConfigTest struct {
	config       Config
	notification Notification
	expectMatch  bool
	expected     RepositoryConfig
}

var config1 = Config{Repositories: []RepositoryConfig{repositoryConfig1, repositoryConfig2}}

var findRepositoryConfigTests = []findRepositoryConfigTest{
	{config1, mockNotification{}, false, RepositoryConfig{}},
	{config1, mockNotification{"https://bitbucket.org/user/repo1", []string{}}, false, RepositoryConfig{}},
	{config1, mockNotification{"https://bitbucket.org/user/repo1", []string{"dev"}}, false, RepositoryConfig{}},
	{config1, mockNotification{"https://bitbucket.org/user/repo1", []string{"dev", "master"}}, true, repositoryConfig1},
	{config1, mockNotification{"https://bitbucket.org/user/repo2", []string{}}, true, repositoryConfig2},
	{config1, mockNotification{"https://bitbucket.org/user/repo2", []string{"dev", "master"}}, true, repositoryConfig2},
}

func TestFindRepositoryConfig(t *testing.T) {
	for i, test := range findRepositoryConfigTests {
		actual, found := test.config.FindRepositoryConfig("", test.notification)
		if found != test.expectMatch {
			t.Errorf("%d. got %#v, want %#v", i, found, test.expectMatch)
		}
		if actual != test.expected {
			t.Errorf("%d. got %#v, want %#v", i, actual, test.expected)
		}
	}
}
