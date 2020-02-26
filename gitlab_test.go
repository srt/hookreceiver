package main

import (
	"reflect"
	"testing"
)

type gitlabParseBytesTest struct {
	input    string
	expected GitlabNotification
}

var gitlabParseBytesTests = []gitlabParseBytesTest{
	{`{
  "before": "95790bf891e76fee5e1747ab589903a6a1f80f22",
  "after": "da1560886d4f094c3e6c9ef40349f7d38b5d27d7",
  "ref": "refs/heads/master",
  "user_id": 4,
  "user_name": "John Smith",
  "project_id": 15,
  "repository": {
    "name": "Diaspora",
    "url": "git@localhost:diaspora.git",
    "description": "",
    "homepage": "http://localhost/diaspora"
  },
  "commits": [
    {
      "id": "b6568db1bc1dcd7f8b4d5a946b0b91f9dacd7327",
      "message": "Update Catalan translation to e38cb41.",
      "timestamp": "2011-12-12T14:27:31+02:00",
      "url": "http://localhost/diaspora/commits/b6568db1bc1dcd7f8b4d5a946b0b91f9dacd7327",
      "author": {
        "name": "Jordi Mallach",
        "email": "jordi@softcatala.org"
      }
    },
    {
      "id": "da1560886d4f094c3e6c9ef40349f7d38b5d27d7",
      "message": "fixed readme",
      "timestamp": "2012-01-03T23:36:29+02:00",
      "url": "http://localhost/diaspora/commits/da1560886d4f094c3e6c9ef40349f7d38b5d27d7",
      "author": {
        "name": "GitLab dev user",
        "email": "gitlabdev@dv6700.(none)"
      }
    }
  ],
  "total_commits_count": 4
}`,
		GitlabNotification{
			Before:    "95790bf891e76fee5e1747ab589903a6a1f80f22",
			After:     "da1560886d4f094c3e6c9ef40349f7d38b5d27d7",
			Ref:       "refs/heads/master",
			UserID:    4,
			UserName:  "John Smith",
			ProjectID: 15,

			Repository: GitlabRepository{
				Name:        "Diaspora",
				URL:         "git@localhost:diaspora.git",
				Description: "",
				Homepage:    "http://localhost/diaspora"},
			Commits: []GitlabCommit{
				{
					ID:        "b6568db1bc1dcd7f8b4d5a946b0b91f9dacd7327",
					Message:   "Update Catalan translation to e38cb41.",
					Timestamp: "2011-12-12T14:27:31+02:00",
					URL:       "http://localhost/diaspora/commits/b6568db1bc1dcd7f8b4d5a946b0b91f9dacd7327",
					Author: GitlabAuthor{
						Name:  "Jordi Mallach",
						Email: "jordi@softcatala.org",
					}},
				{
					ID:        "da1560886d4f094c3e6c9ef40349f7d38b5d27d7",
					Message:   "fixed readme",
					Timestamp: "2012-01-03T23:36:29+02:00",
					URL:       "http://localhost/diaspora/commits/da1560886d4f094c3e6c9ef40349f7d38b5d27d7",
					Author: GitlabAuthor{
						Name:  "GitLab dev user",
						Email: "gitlabdev@dv6700.(none)",
					}}},
			TotalCommitsCount: 4,
		}},

	{`{"before":"0000000000000000000000000000000000000000","after":"8dd3b755cf933430b30835d40c13c903cb27d576","ref":"refs/heads/master","user_id":22,"user_name":"Stefan Reuter","project_id":61,"repository":{"name":"Hilde Hirsch","url":"git@gitlab.example.com:stefan.reuter/hirsch.git","description":"","homepage":"https://gitlab.example.com/stefan.reuter/hirsch"},"commits":[],"total_commits_count":0}`,
		GitlabNotification{
			Before:    "0000000000000000000000000000000000000000",
			After:     "8dd3b755cf933430b30835d40c13c903cb27d576",
			Ref:       "refs/heads/master",
			UserID:    22,
			UserName:  "Stefan Reuter",
			ProjectID: 61,
			Repository: GitlabRepository{
				Name:        "Hilde Hirsch",
				URL:         "git@gitlab.example.com:stefan.reuter/hirsch.git",
				Description: "",
				Homepage:    "https://gitlab.example.com/stefan.reuter/hirsch",
			},
			Commits:           []GitlabCommit{},
			TotalCommitsCount: 0,
		}},
	{`{"before":"6de0da623948c631bcde27b1242268bf00913fb6","after":"bb598afb3d5c08d5e69d3eefcb8638d7b2e1790a","ref":"refs/heads/master","user_id":22,"user_name":"Stefan Reuter","project_id":61,"repository":{"name":"Hilde Hirsch","url":"git@gitlab.example.com:stefan.reuter/hirsch.git","description":"","homepage":"https://gitlab.example.com/stefan.reuter/hirsch"},"commits":[{"id":"bb598afb3d5c08d5e69d3eefcb8638d7b2e1790a","message":"Test","timestamp":"2014-04-22T16:29:41+02:00","url":"https://gitlab.example.com/stefan.reuter/hirsch/commit/bb598afb3d5c08d5e69d3eefcb8638d7b2e1790a","author":{"name":"Stefan Reuter","email":"stefan.reuter@example.com"}}],"total_commits_count":1}`,
		GitlabNotification{
			Before:    "6de0da623948c631bcde27b1242268bf00913fb6",
			After:     "bb598afb3d5c08d5e69d3eefcb8638d7b2e1790a",
			Ref:       "refs/heads/master",
			UserID:    22,
			UserName:  "Stefan Reuter",
			ProjectID: 61,
			Repository: GitlabRepository{
				Name:        "Hilde Hirsch",
				URL:         "git@gitlab.example.com:stefan.reuter/hirsch.git",
				Description: "",
				Homepage:    "https://gitlab.example.com/stefan.reuter/hirsch",
			},
			Commits: []GitlabCommit{
				{
					ID:        "bb598afb3d5c08d5e69d3eefcb8638d7b2e1790a",
					Message:   "Test",
					Timestamp: "2014-04-22T16:29:41+02:00",
					URL:       "https://gitlab.example.com/stefan.reuter/hirsch/commit/bb598afb3d5c08d5e69d3eefcb8638d7b2e1790a",
					Author: GitlabAuthor{
						Name:  "Stefan Reuter",
						Email: "stefan.reuter@example.com",
					}}},
			TotalCommitsCount: 1,
		}},
}

func TestGitlabParseBytes(t *testing.T) {
	for _, test := range gitlabParseBytesTests {
		actual, err := gitlabParseBytes([]byte(test.input))
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(actual, test.expected) {
			t.Errorf("got %#v, want %#v", actual, test.expected)
		}
	}
}
