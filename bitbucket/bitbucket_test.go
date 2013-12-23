package bitbucket

import (
	"reflect"
	"testing"
)

var expected = Notification{
	Canon_url: "https://bitbucket.org",
	User:      "srt",
	Repository: Repository{
		Absolute_url: "/srt/test/",
		Fork:         false,
		Is_private:   true,
		Name:         "test",
		Owner:        "srt",
		Scm:          "git",
		Slug:         "test",
		Website:      ""},
	Commits: []Commit{
		Commit{
			Author: "srt",
			Branch: "master",
			Files: []File{
				File{
					File: "README.md",
					Type: "modified"}},
			Message:      "New date\n",
			Node:         "9d8a38ea7756",
			Parents:      []string{"b8b2e71c4ecd"},
			Raw_author:   "Stefan Reuter <stefan.reuter@example.com>",
			Raw_node:     "9d8a38ea7756a40405dc9bc8f7803700937b58d7",
			Revision:     0,
			Size:         -1,
			Timestamp:    "2013-12-22 03:54:39",
			Utctimestamp: "2013-12-22 02:54:39+00:00"}}}

func TestParseBytes(t *testing.T) {
	actual, err := parseBytes([]byte(`{"repository": {"website": "", "fork": false, "name": "test", "scm": "git", "owner": "srt", "absolute_url": "/srt/test/", "slug": "test", "is_private": true}, "truncated": false, "commits": [{"node": "9d8a38ea7756", "files": [{"type": "modified", "file": "README.md"}], "branch": "master", "utctimestamp": "2013-12-22 02:54:39+00:00", "timestamp": "2013-12-22 03:54:39", "raw_node": "9d8a38ea7756a40405dc9bc8f7803700937b58d7", "message": "New date\n", "size": -1, "author": "srt", "parents": ["b8b2e71c4ecd"], "raw_author": "Stefan Reuter <stefan.reuter@example.com>", "revision": null}], "canon_url": "https://bitbucket.org", "user": "srt"}`))
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v, want %v", actual, expected)
	}
}

func TestBranches(t *testing.T) {
	value, found := expected.Branches()["master"]

	if !found {
		t.Errorf("For branch 'master' no map entry found")
	}
	if !value {
		t.Errorf("For branch 'master' got value %v, want %v", value, true)
	}
}
