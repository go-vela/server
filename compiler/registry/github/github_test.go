// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package github

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/google/go-github/v51/github"
	"golang.org/x/oauth2"
)

func TestGithub_New(t *testing.T) {
	// setup router
	s := httptest.NewServer(http.NotFoundHandler())
	defer s.Close()

	gitClient := github.NewClient(nil)
	gitClient.BaseURL, _ = url.Parse(s.URL + "/api/v3/")

	want := &client{
		Github: gitClient,
		URL:    s.URL,
		API:    s.URL + "/api/v3/",
	}

	// run test
	got, err := New(s.URL, "")

	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
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

	want := &client{
		Github: gitClient,
		URL:    s.URL,
		API:    s.URL + "/api/v3/",
	}

	// run test
	got, err := New(s.URL, token)

	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("New is %+v, want %+v", got.Github, want.Github)
	}
}

func TestGithub_NewURL(t *testing.T) {
	// setup tests
	tests := []struct {
		address string
		want    client
	}{
		{
			// address matches default, so no change to default URL or API.
			address: "https://github.com/",
			want: client{
				URL: "https://github.com/",
				API: "https://api.github.com/",
			},
		},
		{
			// not the default address, but has github.com, so keep default API.
			address: "https://github.com",
			want: client{
				URL: "https://github.com",
				API: "https://api.github.com/",
			},
		},
		{
			// github-enterprise install with /
			address: "https://git.example.com/",
			want: client{
				URL: "https://git.example.com",
				API: "https://git.example.com/api/v3/",
			},
		},
		{
			// github-enterprise install without /
			address: "https://git.example.com",
			want: client{
				URL: "https://git.example.com",
				API: "https://git.example.com/api/v3/",
			},
		},
	}

	// run tests
	for _, test := range tests {
		// run test
		got, err := New(test.address, "foobar")

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
