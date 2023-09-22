// SPDX-License-Identifier: Apache-2.0

package github

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/google/go-github/v54/github"
	"golang.org/x/oauth2"
)

func TestGithub_New(t *testing.T) {
	// setup tests
	tests := []struct {
		failure bool
		id      string
	}{
		{
			failure: false,
			id:      "foo",
		},
		{
			failure: true,
			id:      "",
		},
	}

	// run tests
	for _, test := range tests {
		_, err := New(
			WithAddress("https://github.com/"),
			WithClientID(test.id),
			WithClientSecret("bar"),
			WithServerAddress("https://vela-server.example.com"),
			WithStatusContext("continuous-integration/vela"),
			WithWebUIAddress("https://vela.example.com"),
			WithScopes([]string{"repo", "repo:status", "user:email", "read:user", "read:org"}),
		)

		if test.failure {
			if err == nil {
				t.Errorf("New should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("New returned err: %v", err)
		}
	}
}

func TestGithub_newClientToken(t *testing.T) {
	// setup router
	s := httptest.NewServer(http.NotFoundHandler())
	defer s.Close()

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "foobar"},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	want := github.NewClient(tc)
	want.BaseURL, _ = url.Parse(s.URL + "/api/v3/")

	// setup client
	client, _ := NewTest(s.URL)

	// run test
	got := client.newClientToken("foobar")

	//nolint:staticcheck // ignore false positive
	if got == nil {
		t.Errorf("newClientToken is nil, want %v", want)
	}

	//nolint:staticcheck // ignore false positive
	if !reflect.DeepEqual(got.BaseURL, want.BaseURL) {
		t.Errorf("newClientToken BaseURL is %v, want %v", got.BaseURL, want.BaseURL)
	}
}
