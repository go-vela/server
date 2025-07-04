// SPDX-License-Identifier: Apache-2.0

package github

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-github/v73/github"
	"golang.org/x/oauth2"
)

func TestGithub_New(t *testing.T) {
	// setup router
	s := httptest.NewServer(http.NotFoundHandler())
	defer s.Close()

	gitClient := github.NewClient(nil)

	gitClient.BaseURL, _ = url.Parse(s.URL + "/api/v3/")

	want := &Client{
		githubClient: gitClient,
		URL:          s.URL,
		API:          s.URL + "/api/v3/",
	}

	// run test
	got, err := New(context.Background(), s.URL, "")

	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	if !cmp.Equal(got, want) {
		t.Errorf("New is %v, want %v", got, want)
	}
}

func TestGithub_NewToken(t *testing.T) {
	// setup router
	s := httptest.NewServer(http.NotFoundHandler())
	defer s.Close()

	token := "foobar"
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	gitClient := github.NewClient(tc)
	gitClient.BaseURL, _ = url.Parse(s.URL + "/api/v3/")

	want := &Client{
		githubClient: gitClient,
		URL:          s.URL,
		API:          s.URL + "/api/v3/",
	}

	// run test
	got, err := New(context.Background(), s.URL, token)

	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	if !cmp.Equal(got, want) {
		t.Errorf("New is %v, want %v", got, want)
	}
}

func TestGithub_NewURL(t *testing.T) {
	// setup tests
	tests := []struct {
		address string
		want    Client
	}{
		{
			// address matches default, so no change to default URL or API.
			address: "https://github.com/",
			want: Client{
				URL: "https://github.com/",
				API: "https://api.github.com/",
			},
		},
		{
			// not the default address, but has github.com, so keep default API.
			address: "https://github.com",
			want: Client{
				URL: "https://github.com",
				API: "https://api.github.com/",
			},
		},
		{
			// github-enterprise install with /
			address: "https://git.example.com/",
			want: Client{
				URL: "https://git.example.com",
				API: "https://git.example.com/api/v3/",
			},
		},
		{
			// github-enterprise install without /
			address: "https://git.example.com",
			want: Client{
				URL: "https://git.example.com",
				API: "https://git.example.com/api/v3/",
			},
		},
	}

	// run tests
	for _, test := range tests {
		// run test
		got, err := New(context.Background(), test.address, "foobar")

		if err != nil {
			t.Errorf("New returned err: %v", err)
		}

		if got.URL != test.want.URL {
			t.Errorf("New URL is %v, want %v", got.URL, test.want.URL)
		}

		if got.API != test.want.API {
			t.Errorf("New API is %v, want %v", got.API, test.want.API)
		}
	}
}
