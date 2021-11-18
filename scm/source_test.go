// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package scm

import (
	"testing"
)

func TestSource_New(t *testing.T) {
	// setup tests
	tests := []struct {
		failure bool
		setup   *Setup
	}{
		{
			failure: false,
			setup: &Setup{
				Driver:               "github",
				Address:              "https://github.com",
				ClientID:             "foo",
				ClientSecret:         "bar",
				ServerAddress:        "https://vela-server.example.com",
				ServerWebhookAddress: "",
				StatusContext:        "continuous-integration/vela",
				WebUIAddress:         "https://vela.example.com",
				Scopes:               []string{"repo", "repo:status", "user:email", "read:user", "read:org"},
			},
		},
		{
			failure: true,
			setup: &Setup{
				Driver:               "gitlab",
				Address:              "https://gitlab.com",
				ClientID:             "foo",
				ClientSecret:         "bar",
				ServerAddress:        "https://vela-server.example.com",
				ServerWebhookAddress: "",
				StatusContext:        "continuous-integration/vela",
				WebUIAddress:         "https://vela.example.com",
				Scopes:               []string{"repo", "repo:status", "user:email", "read:user", "read:org"},
			},
		},
		{
			failure: true,
			setup: &Setup{
				Driver:               "bitbucket",
				Address:              "https://bitbucket.org",
				ClientID:             "foo",
				ClientSecret:         "bar",
				ServerAddress:        "https://vela-server.example.com",
				ServerWebhookAddress: "",
				StatusContext:        "continuous-integration/vela",
				WebUIAddress:         "https://vela.example.com",
				Scopes:               []string{"repo", "repo:status", "user:email", "read:user", "read:org"},
			},
		},
		{
			failure: true,
			setup: &Setup{
				Driver:               "github",
				Address:              "",
				ClientID:             "foo",
				ClientSecret:         "bar",
				ServerAddress:        "https://vela-server.example.com",
				ServerWebhookAddress: "",
				StatusContext:        "continuous-integration/vela",
				WebUIAddress:         "https://vela.example.com",
				Scopes:               []string{"repo", "repo:status", "user:email", "read:user", "read:org"},
			},
		},
	}

	// run tests
	for _, test := range tests {
		_, err := New(test.setup)

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
