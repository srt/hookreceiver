package main

import (
	"bytes"
	"reflect"
	"testing"
)

var repositoryConfig1 = RepositoryConfig{
	URL:    "https://bitbucket.org/user/repo1",
	Branch: "master",
}

var repositoryConfig2 = RepositoryConfig{
	URL: "https://bitbucket.org/user/repo2",
}

type mockNotification struct {
	repositoryURL string
	branches      []string
}

func (n mockNotification) RepositoryURL() string {
	return n.repositoryURL
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
		if actual.Branch != test.expected.Branch {
			t.Errorf("%d. got %#v, want %#v", i, actual.Branch, test.expected.Branch)
		}
		if actual.Name != test.expected.Name {
			t.Errorf("%d. got %#v, want %#v", i, actual.Name, test.expected.Name)
		}
		if actual.Dir != test.expected.Dir {
			t.Errorf("%d. got %#v, want %#v", i, actual.Dir, test.expected.Dir)
		}
		if actual.URL != test.expected.URL {
			t.Errorf("%d. got %#v, want %#v", i, actual.URL, test.expected.URL)
		}
		if !reflect.DeepEqual(actual.Command, test.expected.Command) {
			t.Errorf("%d. got %#v, want %#v", i, actual.Command, test.expected.Command)
		}
	}
}

type appendConfigTest struct {
	input    string
	expected Config
}

var appendConfigTests = []appendConfigTest{
	{`{
  "Addr": ":8081",
  "Secret": "t0ps3cr3t",
  "Repositories": [
    {
      "URL": "https://bitbucket.org/srt/foo",
      "Command": [
        "git pull"
      ],
      "Dir": "/var/www/foo"
    },
    {
      "Name": "bar",
      "Command": [
        "git pull"
      ],
      "Dir": "/var/www/bar"
    }
  ]
}`, Config{
		Addr:   ":8081",
		Secret: "t0ps3cr3t",
		Repositories: []RepositoryConfig{
			RepositoryConfig{
				URL: "https://bitbucket.org/srt/foo",
				Command: []string{
					"git pull",
				},
				Dir: "/var/www/foo",
			},
			RepositoryConfig{
				Name: "bar",
				Command: []string{
					"git pull",
				},
				Dir: "/var/www/bar",
			},
		}},
	}}

func TestAppendConfig(t *testing.T) {
	actual := Config{}
	actual.Addr = ":8080"

	for _, test := range appendConfigTests {
		err := appendConfig(&actual, bytes.NewReader([]byte(test.input)))
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(actual, test.expected) {
			t.Errorf("got %#v, want %#v", actual, test.expected)
		}
	}
}
