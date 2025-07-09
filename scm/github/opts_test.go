// SPDX-License-Identifier: Apache-2.0

package github

import (
	"context"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/go-vela/server/tracing"
)

func TestGithub_ClientOpt_WithAddress(t *testing.T) {
	// setup tests
	tests := []struct {
		address string
		want    config
	}{
		{
			address: "https://git.example.com",
			want: config{
				Address: "https://git.example.com",
				API:     "https://git.example.com/api/v3/",
			},
		},
		{
			address: "",
			want: config{
				Address: defaultURL,
				API:     defaultAPI,
			},
		},
	}

	// run tests
	for _, test := range tests {
		_service, err := New(context.Background(),
			WithAddress(test.address),
		)

		if err != nil {
			t.Errorf("WithAddress returned err: %v", err)
		}

		if !reflect.DeepEqual(_service.config.Address, test.want.Address) {
			t.Errorf("WithAddress is %v, want %v", _service.config.Address, test.want.Address)
		}

		if !reflect.DeepEqual(_service.config.API, test.want.API) {
			t.Errorf("WithAddress API is %v, want %v", _service.config.API, test.want.API)
		}
	}
}

func TestGithub_ClientOpt_WithClientID(t *testing.T) {
	// setup tests
	tests := []struct {
		failure bool
		id      string
		want    string
	}{
		{
			failure: false,
			id:      "superSecretClientID",
			want:    "superSecretClientID",
		},
		{
			failure: true,
			id:      "",
			want:    "",
		},
	}

	// run tests
	for _, test := range tests {
		_service, err := New(context.Background(),
			WithClientID(test.id),
		)

		if test.failure {
			if err == nil {
				t.Errorf("WithClientID should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("WithClientID returned err: %v", err)
		}

		if !reflect.DeepEqual(_service.config.ClientID, test.want) {
			t.Errorf("WithClientID is %v, want %v", _service.config.ClientID, test.want)
		}
	}
}

func TestGithub_ClientOpt_WithClientSecret(t *testing.T) {
	// setup tests
	tests := []struct {
		failure bool
		secret  string
		want    string
	}{
		{
			failure: false,
			secret:  "superSecretClientSecret",
			want:    "superSecretClientSecret",
		},
		{
			failure: true,
			secret:  "",
			want:    "",
		},
	}

	// run tests
	for _, test := range tests {
		_service, err := New(context.Background(),
			WithClientSecret(test.secret),
		)

		if test.failure {
			if err == nil {
				t.Errorf("WithClientSecret should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("WithClientSecret returned err: %v", err)
		}

		if !reflect.DeepEqual(_service.config.ClientSecret, test.want) {
			t.Errorf("WithClientSecret is %v, want %v", _service.config.ClientSecret, test.want)
		}
	}
}

func TestGithub_ClientOpt_WithServerAddress(t *testing.T) {
	// setup tests
	tests := []struct {
		failure bool
		address string
		want    string
	}{
		{
			failure: false,
			address: "https://vela.example.com",
			want:    "https://vela.example.com",
		},
		{
			failure: true,
			address: "",
			want:    "",
		},
	}

	// run tests
	for _, test := range tests {
		_service, err := New(context.Background(),
			WithServerAddress(test.address),
		)

		if test.failure {
			if err == nil {
				t.Errorf("WithServerAddress should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("WithServerAddress returned err: %v", err)
		}

		if !reflect.DeepEqual(_service.config.ServerAddress, test.want) {
			t.Errorf("WithServerAddress is %v, want %v", _service.config.ServerAddress, test.want)
		}
	}
}

func TestGithub_ClientOpt_WithServerWebhookAddress(t *testing.T) {
	// setup tests
	tests := []struct {
		failure        bool
		address        string
		webhookAddress string
		want           string
	}{
		{
			failure:        false,
			address:        "https://vela.example.com",
			webhookAddress: "",
			want:           "https://vela.example.com/webhook",
		},
		{
			failure:        false,
			address:        "https://vela.example.com",
			webhookAddress: "https://vela.example.com",
			want:           "https://vela.example.com/webhook",
		},
		{
			failure:        false,
			address:        "https://vela.example.com",
			webhookAddress: "https://vela-alternative.example.com",
			want:           "https://vela-alternative.example.com",
		},
	}

	// run tests
	for _, test := range tests {
		_service, err := New(context.Background(),
			WithServerAddress(test.address),
			WithServerWebhookAddress(test.webhookAddress),
		)

		if test.failure {
			if err == nil {
				t.Errorf("WithServerWebhookAddress should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("WithServerWebhookAddress returned err: %v", err)
		}

		if !reflect.DeepEqual(_service.config.ServerWebhookAddress, test.want) {
			t.Errorf("WithServerWebhookAddress is %v, want %v", _service.config.ServerWebhookAddress, test.want)
		}
	}
}

func TestGithub_ClientOpt_WithStatusContext(t *testing.T) {
	// setup tests
	tests := []struct {
		failure bool
		context string
		want    string
	}{
		{
			failure: false,
			context: "continuous-integration/vela",
			want:    "continuous-integration/vela",
		},
		{
			failure: true,
			context: "",
			want:    "",
		},
	}

	// run tests
	for _, test := range tests {
		_service, err := New(context.Background(),
			WithStatusContext(test.context),
		)

		if test.failure {
			if err == nil {
				t.Errorf("WithStatusContext should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("WithStatusContext returned err: %v", err)
		}

		if !reflect.DeepEqual(_service.config.StatusContext, test.want) {
			t.Errorf("WithStatusContext is %v, want %v", _service.config.StatusContext, test.want)
		}
	}
}

func TestGithub_ClientOpt_WithWebUIAddress(t *testing.T) {
	// setup tests
	tests := []struct {
		address string
		want    string
	}{
		{
			address: "https://vela.example.com",
			want:    "https://vela.example.com",
		},
		{
			address: "",
			want:    "",
		},
	}

	// run tests
	for _, test := range tests {
		_service, err := New(context.Background(),
			WithWebUIAddress(test.address),
		)

		if err != nil {
			t.Errorf("WithWebUIAddress returned err: %v", err)
		}

		if !reflect.DeepEqual(_service.config.WebUIAddress, test.want) {
			t.Errorf("WithWebUIAddress is %v, want %v", _service.config.WebUIAddress, test.want)
		}
	}
}

func TestGithub_ClientOpt_WithOAuthScopes(t *testing.T) {
	// setup tests
	tests := []struct {
		failure bool
		scopes  []string
		want    []string
	}{
		{
			failure: false,
			scopes:  []string{"repo", "repo:status", "user:email", "read:user", "read:org"},
			want:    []string{"repo", "repo:status", "user:email", "read:user", "read:org"},
		},
		{
			failure: true,
			scopes:  []string{},
			want:    []string{},
		},
	}

	// run tests
	for _, test := range tests {
		_service, err := New(context.Background(),
			WithOAuthScopes(test.scopes),
		)

		if test.failure {
			if err == nil {
				t.Errorf("WithOAuthScopes should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("WithOAuthScopes returned err: %v", err)
		}

		if !reflect.DeepEqual(_service.config.OAuthScopes, test.want) {
			t.Errorf("WithOAuthScopes is %v, want %v", _service.config.OAuthScopes, test.want)
		}
	}
}

func TestGithub_ClientOpt_WithTracing(t *testing.T) {
	// setup tests
	tests := []struct {
		failure bool
		tracing *tracing.Client
		want    *tracing.Client
	}{
		{
			failure: false,
			tracing: &tracing.Client{},
			want:    &tracing.Client{},
		},
	}

	// run tests
	for _, test := range tests {
		_service, err := New(context.Background(),
			WithTracing(test.tracing),
		)

		if test.failure {
			if err == nil {
				t.Errorf("WithTracing should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("WithTracing returned err: %v", err)
		}

		if !reflect.DeepEqual(_service.Tracing, test.want) {
			t.Errorf("WithTracing is %v, want %v", _service.Tracing, test.want)
		}
	}
}

func TestGithub_ClientOpt_WithGitHubAppPermissions(t *testing.T) {
	// setup tests
	tests := []struct {
		failure     bool
		permissions []string
		want        []string
	}{
		{
			failure:     false,
			permissions: []string{"contents:read"},
			want:        []string{"contents:read"},
		},
		{
			failure:     false,
			permissions: []string{},
			want:        []string{},
		},
	}

	// run tests
	for _, test := range tests {
		_service, err := New(context.Background(),
			WithGitHubAppPermissions(test.permissions),
		)

		if test.failure {
			if err == nil {
				t.Errorf("WithGitHubAppPermissions should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("WithGitHubAppPermissions returned err: %v", err)
		}

		if diff := cmp.Diff(test.want, _service.config.AppPermissions); diff != "" {
			t.Errorf("WithGitHubAppPermissions mismatch (-want +got):\n%s", diff)
		}
	}
}

func TestGithub_ClientOpt_WithRepoRoleMap(t *testing.T) {
	// setup tests
	tests := []struct {
		failure bool
		mapping map[string]string
		want    map[string]string
	}{
		{
			failure: false,
			mapping: map[string]string{"vela": "vela"},
			want:    map[string]string{"vela": "vela"},
		},
		{
			failure: false,
			mapping: map[string]string{},
			want:    map[string]string{},
		},
	}

	// run tests
	for _, test := range tests {
		_service, err := New(context.Background(),
			WithRepoRoleMap(test.mapping),
		)

		if test.failure {
			if err == nil {
				t.Errorf("WithRepoRoleMap should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("WithRepoRoleMap returned err: %v", err)
		}

		if diff := cmp.Diff(test.want, _service.GetRepoRoleMap()); diff != "" {
			t.Errorf("WithRepoRoleMap mismatch (-want +got):\n%s", diff)
		}
	}
}

func TestGithub_ClientOpt_WithOrgRoleMap(t *testing.T) {
	// setup tests
	tests := []struct {
		failure bool
		mapping map[string]string
		want    map[string]string
	}{
		{
			failure: false,
			mapping: map[string]string{"vela": "vela"},
			want:    map[string]string{"vela": "vela"},
		},
		{
			failure: false,
			mapping: map[string]string{},
			want:    map[string]string{},
		},
	}

	// run tests
	for _, test := range tests {
		_service, err := New(context.Background(),
			WithOrgRoleMap(test.mapping),
		)

		if test.failure {
			if err == nil {
				t.Errorf("WithOrgRoleMap should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("WithOrgRoleMap returned err: %v", err)
		}

		if diff := cmp.Diff(test.want, _service.GetOrgRoleMap()); diff != "" {
			t.Errorf("WithOrgRoleMap mismatch (-want +got):\n%s", diff)
		}
	}
}

func TestGithub_ClientOpt_WithTeamRoleMap(t *testing.T) {
	// setup tests
	tests := []struct {
		failure bool
		mapping map[string]string
		want    map[string]string
	}{
		{
			failure: false,
			mapping: map[string]string{"vela": "vela"},
			want:    map[string]string{"vela": "vela"},
		},
		{
			failure: false,
			mapping: map[string]string{},
			want:    map[string]string{},
		},
	}

	// run tests
	for _, test := range tests {
		_service, err := New(context.Background(),
			WithTeamRoleMap(test.mapping),
		)

		if test.failure {
			if err == nil {
				t.Errorf("WithTeamRoleMap should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("WithTeamRoleMap returned err: %v", err)
		}

		if diff := cmp.Diff(test.want, _service.GetTeamRoleMap()); diff != "" {
			t.Errorf("WithTeamRoleMap mismatch (-want +got):\n%s", diff)
		}
	}
}
