// SPDX-License-Identifier: Apache-2.0

package scm

import (
	"context"
	"reflect"
	"testing"

	"github.com/go-vela/server/constants"
)

func TestSCM_Setup_Github(t *testing.T) {
	// setup types
	_setup := &Setup{
		Driver:               "github",
		Address:              "https://github.com",
		ClientID:             "foo",
		ClientSecret:         "bar",
		ServerAddress:        "https://vela-server.example.com",
		ServerWebhookAddress: "",
		StatusContext:        "continuous-integration/vela",
		WebUIAddress:         "https://vela.example.com",
		OAuthScopes:          []string{"repo", "repo:status", "user:email", "read:user", "read:org"},
		RepoRoleMap:          map[string]string{"read": constants.PermissionRead, "write": constants.PermissionWrite, "admin": constants.PermissionAdmin},
		OrgRoleMap:           map[string]string{"read": constants.PermissionRead, "admin": constants.PermissionAdmin},
		TeamRoleMap:          map[string]string{"admin": constants.PermissionAdmin},
	}

	_github, err := _setup.Github(context.Background())
	if err != nil {
		t.Errorf("unable to setup scm: %v", err)
	}

	// setup tests
	tests := []struct {
		failure bool
		setup   *Setup
		want    Service
	}{
		{
			failure: false,
			setup:   _setup,
			want:    _github,
		},
		{
			failure: true,
			setup:   &Setup{Driver: "github"},
			want:    nil,
		},
	}

	// run tests
	for _, test := range tests {
		got, err := test.setup.Github(context.Background())

		if test.failure {
			if err == nil {
				t.Errorf("Github should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("Github returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("Github is %v, want %v", got, test.want)
		}
	}
}

func TestSCM_Setup_Gitlab(t *testing.T) {
	// setup types
	_setup := &Setup{
		Driver:               "gitlab",
		Address:              "https://gitlab.com",
		ClientID:             "foo",
		ClientSecret:         "bar",
		ServerAddress:        "https://vela-server.example.com",
		ServerWebhookAddress: "",
		StatusContext:        "continuous-integration/vela",
		WebUIAddress:         "https://vela.example.com",
	}

	// run test
	got, err := _setup.Gitlab(context.Background())
	if err == nil {
		t.Errorf("Gitlab should have returned err")
	}

	if got != nil {
		t.Errorf("Gitlab is %v, want nil", got)
	}
}

func TestSCM_Setup_Validate(t *testing.T) {
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
				OAuthScopes:          []string{"repo", "repo:status", "user:email", "read:user", "read:org"},
			},
		},
		{
			failure: false,
			setup: &Setup{
				Driver:               "gitlab",
				Address:              "https://gitlab.com",
				ClientID:             "foo",
				ClientSecret:         "bar",
				ServerAddress:        "https://vela-server.example.com",
				ServerWebhookAddress: "",
				StatusContext:        "continuous-integration/vela",
				WebUIAddress:         "https://vela.example.com",
				OAuthScopes:          []string{"repo", "repo:status", "user:email", "read:user", "read:org"},
			},
		},
		{
			failure: true,
			setup: &Setup{
				Driver:               "github",
				Address:              "https://github.com/",
				ClientID:             "foo",
				ClientSecret:         "bar",
				ServerAddress:        "https://vela-server.example.com",
				ServerWebhookAddress: "",
				StatusContext:        "continuous-integration/vela",
				WebUIAddress:         "https://vela.example.com",
				OAuthScopes:          []string{"repo", "repo:status", "user:email", "read:user", "read:org"},
			},
		},
		{
			failure: true,
			setup: &Setup{
				Driver:               "github",
				Address:              "github.com",
				ClientID:             "foo",
				ClientSecret:         "bar",
				ServerAddress:        "https://vela-server.example.com",
				ServerWebhookAddress: "",
				StatusContext:        "continuous-integration/vela",
				WebUIAddress:         "https://vela.example.com",
				OAuthScopes:          []string{"repo", "repo:status", "user:email", "read:user", "read:org"},
			},
		},
		{
			failure: true,
			setup: &Setup{
				Driver:               "",
				Address:              "https://github.com",
				ClientID:             "foo",
				ClientSecret:         "bar",
				ServerAddress:        "https://vela-server.example.com",
				ServerWebhookAddress: "",
				StatusContext:        "continuous-integration/vela",
				WebUIAddress:         "https://vela.example.com",
				OAuthScopes:          []string{"repo", "repo:status", "user:email", "read:user", "read:org"},
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
				OAuthScopes:          []string{"repo", "repo:status", "user:email", "read:user", "read:org"},
			},
		},
		{
			failure: true,
			setup: &Setup{
				Driver:               "github",
				Address:              "https://github.com",
				ClientID:             "",
				ClientSecret:         "bar",
				ServerAddress:        "https://vela-server.example.com",
				ServerWebhookAddress: "",
				StatusContext:        "continuous-integration/vela",
				WebUIAddress:         "https://vela.example.com",
				OAuthScopes:          []string{"repo", "repo:status", "user:email", "read:user", "read:org"},
			},
		},
		{
			failure: true,
			setup: &Setup{
				Driver:               "github",
				Address:              "https://github.com",
				ClientID:             "foo",
				ClientSecret:         "",
				ServerAddress:        "https://vela-server.example.com",
				ServerWebhookAddress: "",
				StatusContext:        "continuous-integration/vela",
				WebUIAddress:         "https://vela.example.com",
				OAuthScopes:          []string{"repo", "repo:status", "user:email", "read:user", "read:org"},
			},
		},
		{
			failure: true,
			setup: &Setup{
				Driver:               "github",
				Address:              "https://github.com",
				ClientID:             "foo",
				ClientSecret:         "bar",
				ServerAddress:        "https://vela-server.example.com",
				ServerWebhookAddress: "",
				StatusContext:        "",
				WebUIAddress:         "https://vela.example.com",
				OAuthScopes:          []string{"repo", "repo:status", "user:email", "read:user", "read:org"},
			},
		},
		{
			failure: true,
			setup: &Setup{
				Driver:               "github",
				Address:              "https://github.com",
				ClientID:             "foo",
				ClientSecret:         "bar",
				ServerAddress:        "https://vela-server.example.com",
				ServerWebhookAddress: "",
				StatusContext:        "continuous-integration/vela",
				WebUIAddress:         "https://vela.example.com",
				OAuthScopes:          []string{},
			},
		},
	}

	// run tests
	for _, test := range tests {
		err := test.setup.Validate()

		if test.failure {
			if err == nil {
				t.Errorf("Validate should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("Validate returned err: %v", err)
		}
	}
}
