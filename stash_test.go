package main

import (
	"reflect"
	"testing"
)

type stashParseBytesTest struct {
	input    string
	expected StashNotification
}

var stashParseBytesTests = []stashParseBytesTest{
	{`{"repository":{"slug":"test","id":61,"name":"test","scmId":"git","state":"AVAILABLE","statusMessage":"Available","forkable":true,"project":{"key":"~SRT","id":2,"name":"Stefan Reuter","type":"PERSONAL","isPersonal":true,"owner":{"name":"srt","emailAddress":"stefan.reuter@example.com","id":5,"displayName":"Stefan Reuter","active":true,"slug":"srt","type":"NORMAL"}},"public":false},"refChanges":[{"refId":"refs/heads/master","fromHash":"8ba5559a89a4c4ae0402cac4f8deec8b8bcb2ae1","toHash":"fec32d23bb614eabd451d8c6e6f9dc14378dcf7c","type":"UPDATE"}],"changesets":{"size":1,"limit":100,"isLastPage":true,"values":[{"fromCommit":{"id":"8ba5559a89a4c4ae0402cac4f8deec8b8bcb2ae1","displayId":"8ba5559"},"toCommit":{"id":"fec32d23bb614eabd451d8c6e6f9dc14378dcf7c","displayId":"fec32d2","author":{"name":"Stefan Reuter","emailAddress":"stefan.reuter@example.com"},"authorTimestamp":1398238608000,"message":"...","parents":[{"id":"8ba5559a89a4c4ae0402cac4f8deec8b8bcb2ae1","displayId":"8ba5559"}]},"changes":{"size":1,"limit":100,"isLastPage":true,"values":[{"contentId":"03ef7017032edcf73758eae1986fb2b772d0f26d","path":{"components":["README.md"],"parent":"","name":"README.md","extension":"md","toString":"README.md"},"executable":false,"percentUnchanged":-1,"type":"MODIFY","nodeType":"FILE","srcExecutable":false,"link":{"url":"/users/srt/repos/test/commits/fec32d23bb614eabd451d8c6e6f9dc14378dcf7c#README.md","rel":"self"},"links":{"self":[{"href":"https://git.example.com/users/srt/repos/test/commits/fec32d23bb614eabd451d8c6e6f9dc14378dcf7c#README.md"}]}}],"start":0,"filter":null},"link":{"url":"/users/srt/repos/test/commits/fec32d23bb614eabd451d8c6e6f9dc14378dcf7c#README.md","rel":"self"},"links":{"self":[{"href":"https://git.example.com/users/srt/repos/test/commits/fec32d23bb614eabd451d8c6e6f9dc14378dcf7c#README.md"}]}}],"start":0,"filter":null}}`,
		StashNotification{
			Repository: StashRepository{
				Slug:          "test",
				ID:            61,
				Name:          "test",
				ScmID:         "git",
				State:         "AVAILABLE",
				StatusMessage: "Available",
				Forkable:      true,
				Project: StashProject{
					Key:        "~SRT",
					ID:         2,
					Name:       "Stefan Reuter",
					Type:       "PERSONAL",
					IsPersonal: true,
					Owner: StashUser{
						Name:         "srt",
						EmailAddress: "stefan.reuter@example.com",
						ID:           5,
						DisplayName:  "Stefan Reuter",
						Active:       true,
						Slug:         "srt",
						Type:         "NORMAL",
					},
				},
				Public: false,
			},
			RefChanges: []StashRefChange{
				StashRefChange{
					RefID:    "refs/heads/master",
					FromHash: "8ba5559a89a4c4ae0402cac4f8deec8b8bcb2ae1",
					ToHash:   "fec32d23bb614eabd451d8c6e6f9dc14378dcf7c",
					Type:     "UPDATE",
				},
			},
		}},
}

func TestStashParseBytes(t *testing.T) {
	for _, test := range stashParseBytesTests {
		actual, err := stashParseBytes([]byte(test.input))
		if err != nil {
			t.Error(err)
		}
		// Only check Repository and RefChanges as Changesets is unparsed
		if !reflect.DeepEqual(actual.Repository, test.expected.Repository) {
			t.Errorf("got %#v, want %#v", actual.Repository, test.expected.Repository)
		}
		if !reflect.DeepEqual(actual.RefChanges, test.expected.RefChanges) {
			t.Errorf("got %#v, want %#v", actual.RefChanges, test.expected.RefChanges)
		}
	}
}
