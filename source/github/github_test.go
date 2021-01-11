// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package github

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/google/go-github/v29/github"
	"golang.org/x/oauth2"
)

func TestGithub_New(t *testing.T) {
	// setup router
	s := httptest.NewServer(http.NotFoundHandler())
	defer s.Close()

	id := "foo"
	secret := "bar"

	want := &client{
		URL:           s.URL,
		API:           s.URL + "/api/v3/",
		LocalHost:     s.URL,
		WebUIHost:     s.URL,
		StatusContext: "continuous-integration/vela",
		OConfig: &oauth2.Config{
			ClientID:     id,
			ClientSecret: secret,
			Scopes:       []string{"repo", "repo:status", "user:email", "read:user", "read:org"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  fmt.Sprintf("%s/login/oauth/authorize", s.URL),
				TokenURL: fmt.Sprintf("%s/login/oauth/access_token", s.URL),
			},
		},
		AuthReq: &github.AuthorizationRequest{
			ClientID:     &id,
			ClientSecret: &secret,
			Scopes:       []github.Scope{"repo", "repo:status", "user:email", "read:user", "read:org"},
		},
	}

	// run test
	got, err := NewTest(s.URL)

	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("New is %v, want %v", got, want)
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

	if got == nil {
		t.Errorf("newClientToken is nil, want %v", want)
	}

	if !reflect.DeepEqual(got.BaseURL, want.BaseURL) {
		t.Errorf("newClientToken BaseURL is %v, want %v", got.BaseURL, want.BaseURL)
	}
}
