// SPDX-License-Identifier: Apache-2.0

package github

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGithub_New(t *testing.T) {
	// setup router
	s := httptest.NewServer(http.NotFoundHandler())
	defer s.Close()

	wantURL := s.URL
	wantAPI := s.URL + "/api/v3/"

	// run test
	got, err := New(context.Background(), s.URL, "", nil)
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	if got.URL != wantURL {
		t.Errorf("New URL is %v, want %v", got.URL, wantURL)
	}

	if got.API != wantAPI {
		t.Errorf("New API is %v, want %v", got.API, wantAPI)
	}
}

func TestGithub_NewToken(t *testing.T) {
	// setup router
	s := httptest.NewServer(http.NotFoundHandler())
	defer s.Close()

	token := "foobar"
	wantURL := s.URL
	wantAPI := s.URL + "/api/v3/"

	// run test
	got, err := New(context.Background(), s.URL, token, nil)
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	if got.URL != wantURL {
		t.Errorf("New URL is %v, want %v", got.URL, wantURL)
	}

	if got.API != wantAPI {
		t.Errorf("New API is %v, want %v", got.API, wantAPI)
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
		got, err := New(context.Background(), test.address, "foobar", nil)
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
