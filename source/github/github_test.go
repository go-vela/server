// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
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

	"github.com/google/go-github/v29/github"
	"golang.org/x/oauth2"
)

func TestGithub_New(t *testing.T) {
	// setup tests
	tests := []struct {
		failure bool
		address string
	}{
		{
			failure: false,
			address: "https://github.com",
		},
		{
			failure: true,
			address: "",
		},
	}

	// run tests
	for _, test := range tests {
		_, err := New(
			WithAddress(test.address),
			WithClientID("foo"),
			WithClientSecret("bar"),
			WithServerAddress("https://vela-server.example.com"),
			WithStatusContext("continuous-integration/vela"),
			WithWebUIAddress("https://vela.example.com"),
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

	if got == nil {
		t.Errorf("newClientToken is nil, want %v", want)
	}

	// nolint: staticcheck // ignore false positive
	if !reflect.DeepEqual(got.BaseURL, want.BaseURL) {
		t.Errorf("newClientToken BaseURL is %v, want %v", got.BaseURL, want.BaseURL)
	}
}
