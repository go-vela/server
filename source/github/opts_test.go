// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
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
		want    string
	}{
		{
			address: "https://git.example.com",
			want:    "https://git.example.com",
		},
		{
			address: "",
			want:    defaultURL,
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

		if !reflect.DeepEqual(_service.config.Address, test.want) {
			t.Errorf("WithAddress is %v, want %v", _service.config.Address, test.want)
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
