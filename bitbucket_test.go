package main

import (
	"reflect"
	"testing"
)

type parseBytesTest struct {
	input    string
	expected BitbucketNotification
}

var parseBytesTests = []parseBytesTest{
	{`{"repository": {"website": "", "fork": false, "name": "test", "scm": "git", "owner": "srt", "absolute_url": "/srt/test/", "slug": "test", "is_private": true}, "truncated": false, "commits": [{"node": "9d8a38ea7756", "files": [{"type": "modified", "file": "README.md"}], "branch": "master", "utctimestamp": "2013-12-22 02:54:39+00:00", "timestamp": "2013-12-22 03:54:39", "raw_node": "9d8a38ea7756a40405dc9bc8f7803700937b58d7", "message": "New date\n", "size": -1, "author": "srt", "parents": ["b8b2e71c4ecd"], "raw_author": "Stefan Reuter <stefan.reuter@example.com>", "revision": null}], "canon_url": "https://bitbucket.org", "user": "srt"}`,
		BitbucketNotification{
			Canon_url: "https://bitbucket.org",
			User:      "srt",
			Repository: BitbucketRepository{
				Absolute_url: "/srt/test/",
				Fork:         false,
				Is_private:   true,
				Name:         "test",
				Owner:        "srt",
				Scm:          "git",
				Slug:         "test",
				Website:      ""},
			Commits: []BitbucketCommit{
				BitbucketCommit{
					Author: "srt",
					Branch: "master",
					Files: []BitbucketFile{
						BitbucketFile{
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
					Utctimestamp: "2013-12-22 02:54:39+00:00"}}}},
	{`{"repository": {"website": "", "fork": false, "name": "hirsch-forum", "scm": "git", "owner": "hilde", "absolute_url": "/hilde/hirsch-forum/", "slug": "hirsch-forum", "is_private": true}, "truncated": false, "commits": [{"node": "a0e92085d5ee", "files": [{"type": "modified", "file": "styles/merge/template/overall_header.html"}], "branches": ["dev3.0", "merge", "dev", "mobile"], "raw_author": "hans <hans@example.com>", "utctimestamp": "2013-12-26 00:43:36+00:00", "author": "hans", "timestamp": "2013-12-26 01:43:36", "raw_node": "a0e92085d5eeb3fe346d22e90039263520825819", "parents": ["cad6e3b7ce49"], "branch": null, "message": "Add missing base tag\n", "revision": null, "size": -1}, {"node": "bc44c7f1428f", "files": [{"type": "modified", "file": "styles/merge/template/overall_header.html"}], "raw_author": "hans <hans@example.com>", "utctimestamp": "2013-12-26 00:43:52+00:00", "author": "hans", "timestamp": "2013-12-26 01:43:52", "raw_node": "bc44c7f1428f8d24a7d2979fe1b09b4141f5d494", "parents": ["32f4be438ae8", "a0e92085d5ee"], "branch": "master", "message": "Merge branch 'dev'\n", "revision": null, "size": -1}], "canon_url": "https://bitbucket.org", "user": "hans"}`,
		BitbucketNotification{
			Canon_url: "https://bitbucket.org",
			User:      "hans",
			Repository: BitbucketRepository{
				Absolute_url: "/hilde/hirsch-forum/",
				Fork:         false,
				Is_private:   true,
				Name:         "hirsch-forum",
				Owner:        "hilde",
				Scm:          "git",
				Slug:         "hirsch-forum",
				Website:      ""},
			Commits: []BitbucketCommit{
				BitbucketCommit{
					Author:   "hans",
					Branches: []string{"dev3.0", "merge", "dev", "mobile"},
					Branch:   "",
					Files: []BitbucketFile{
						BitbucketFile{
							File: "styles/merge/template/overall_header.html",
							Type: "modified"}},
					Message:      "Add missing base tag\n",
					Node:         "a0e92085d5ee",
					Parents:      []string{"cad6e3b7ce49"},
					Raw_author:   "hans <hans@example.com>",
					Raw_node:     "a0e92085d5eeb3fe346d22e90039263520825819",
					Revision:     0,
					Size:         -1,
					Timestamp:    "2013-12-26 01:43:36",
					Utctimestamp: "2013-12-26 00:43:36+00:00"},
				BitbucketCommit{Author: "hans",
					Branches: []string(nil),
					Branch:   "master",
					Files: []BitbucketFile{
						BitbucketFile{
							File: "styles/merge/template/overall_header.html",
							Type: "modified"}},
					Message:      "Merge branch 'dev'\n",
					Node:         "bc44c7f1428f",
					Parents:      []string{"32f4be438ae8", "a0e92085d5ee"},
					Raw_author:   "hans <hans@example.com>",
					Raw_node:     "bc44c7f1428f8d24a7d2979fe1b09b4141f5d494",
					Revision:     0,
					Size:         -1,
					Timestamp:    "2013-12-26 01:43:52",
					Utctimestamp: "2013-12-26 00:43:52+00:00"}}}},
}

func TestParseBytes(t *testing.T) {
	for _, test := range parseBytesTests {
		actual, err := parseBytes([]byte(test.input))
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(actual, test.expected) {
			t.Errorf("got %#v, want %#v", actual, test.expected)
		}
	}
}

type branchesTest struct {
	notification BitbucketNotification
	branches     []string
}

var branchesTests = []branchesTest{
	// one commit, one branch
	{BitbucketNotification{Commits: []BitbucketCommit{
		BitbucketCommit{Branch: "master"}}},
		[]string{"master"}},
	// three commits, one without a branch and two each with one branch
	{BitbucketNotification{Commits: []BitbucketCommit{
		BitbucketCommit{Branch: "master"},
		BitbucketCommit{},
		BitbucketCommit{Branch: "dev"}}},
		[]string{"master", "dev"}},
	// two commits, one with two branches and one with 1 branch
	{BitbucketNotification{Commits: []BitbucketCommit{
		BitbucketCommit{Branches: []string{"master", "bastard"}},
		BitbucketCommit{Branch: "dev"}}},
		[]string{"master", "bastard", "dev"}},
	{BitbucketNotification{Commits: []BitbucketCommit{
		BitbucketCommit{Branches: []string{"master", "bastard"}},
		BitbucketCommit{Branch: "master"}}},
		[]string{"master", "bastard"}},
	{parseBytesTests[0].expected, []string{"master"}},
	{parseBytesTests[1].expected, []string{"master", "dev3.0", "merge", "dev", "mobile"}},
}

func TestBranches(t *testing.T) {
	for _, test := range branchesTests {
		for _, expectedBranch := range test.branches {
			value, found := test.notification.Branches()[expectedBranch]
			if !found {
				t.Errorf("For branch %q no map entry found", expectedBranch)
			}
			if !value {
				t.Errorf("For branch %q got value %v, want %v", expectedBranch, value, true)
			}
		}
	}
}
