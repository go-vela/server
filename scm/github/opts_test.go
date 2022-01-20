// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package github

import (
	"reflect"
	"testing"
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
		_service, err := New(
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
		_service, err := New(
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
		_service, err := New(
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
		_service, err := New(
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
			want:           "https://vela.example.com",
		},
		{
			failure:        false,
			address:        "https://vela.example.com",
			webhookAddress: "https://vela.example.com",
			want:           "https://vela.example.com",
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
		_service, err := New(
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
		_service, err := New(
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
		_service, err := New(
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

func TestGithub_ClientOpt_WithScopes(t *testing.T) {
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
		_service, err := New(
			WithScopes(test.scopes),
		)

		if test.failure {
			if err == nil {
				t.Errorf("WithScopes should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("WithScopes returned err: %v", err)
		}

		if !reflect.DeepEqual(_service.config.Scopes, test.want) {
			t.Errorf("WithScopes is %v, want %v", _service.config.Scopes, test.want)
		}
	}
}
