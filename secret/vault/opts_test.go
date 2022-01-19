// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package vault

import (
	"reflect"
	"testing"
	"time"
)

func TestVault_ClientOpt_WithAddress(t *testing.T) {
	// setup tests
	tests := []struct {
		failure bool
		address string
		want    string
	}{
		{
			failure: false,
			address: "https://vault.example.com",
			want:    "https://vault.example.com",
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
			WithAddress(test.address),
			WithVersion("1"),
		)

		if test.failure {
			if err == nil {
				t.Errorf("WithAddress should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("WithAddress returned err: %v", err)
		}

		if !reflect.DeepEqual(_service.config.Address, test.want) {
			t.Errorf("WithAddress is %v, want %v", _service.config.Address, test.want)
		}
	}
}

func TestVault_ClientOpt_WithAWSRole(t *testing.T) {
	// setup tests
	tests := []struct {
		role string
		want string
	}{
		{
			role: "foo",
			want: "foo",
		},
		{
			role: "",
			want: "",
		},
	}

	// run tests
	for _, test := range tests {
		_service, err := New(
			WithAddress("https://vault.example.com"),
			WithAWSRole(test.role),
			WithVersion("1"),
		)

		if err != nil {
			t.Errorf("WithAWSRole returned err: %v", err)
		}

		if !reflect.DeepEqual(_service.config.AWSRole, test.want) {
			t.Errorf("WithAWSRole is %v, want %v", _service.config.AWSRole, test.want)
		}
	}
}

func TestVault_ClientOpt_WithPrefix(t *testing.T) {
	// setup tests
	tests := []struct {
		prefix string
		want   string
	}{
		{
			prefix: "foo",
			want:   "secret/foo",
		},
		{
			prefix: "",
			want:   "secret",
		},
	}

	// run tests
	for _, test := range tests {
		_service, err := New(
			WithAddress("https://vault.example.com"),
			WithPrefix(test.prefix),
			WithVersion("1"),
		)

		if err != nil {
			t.Errorf("WithPrefix returned err: %v", err)
		}

		if !reflect.DeepEqual(_service.config.Prefix, test.want) {
			t.Errorf("WithPrefix is %v, want %v", _service.config.Prefix, test.want)
		}
	}
}

func TestVault_ClientOpt_WithToken(t *testing.T) {
	// setup tests
	tests := []struct {
		token string
		want  string
	}{
		{
			token: "foo",
			want:  "foo",
		},
		{
			token: "",
			want:  "",
		},
	}

	// run tests
	for _, test := range tests {
		_service, err := New(
			WithAddress("https://vault.example.com"),
			WithToken(test.token),
			WithVersion("1"),
		)

		if err != nil {
			t.Errorf("WithToken returned err: %v", err)
		}

		if !reflect.DeepEqual(_service.config.Token, test.want) {
			t.Errorf("WithToken is %v, want %v", _service.config.Token, test.want)
		}
	}
}

func TestVault_ClientOpt_WithTokenDuration(t *testing.T) {
	// setup tests
	tests := []struct {
		tokenDuration time.Duration
		want          time.Duration
	}{
		{
			tokenDuration: 5 * time.Minute,
			want:          5 * time.Minute,
		},
		{
			tokenDuration: 0,
			want:          0,
		},
	}

	// run tests
	for _, test := range tests {
		_service, err := New(
			WithAddress("https://vault.example.com"),
			WithTokenDuration(test.tokenDuration),
			WithVersion("1"),
		)

		if err != nil {
			t.Errorf("WithTokenDuration returned err: %v", err)
		}

		if !reflect.DeepEqual(_service.config.TokenDuration, test.want) {
			t.Errorf("WithTokenDuration is %v, want %v", _service.config.TokenDuration, test.want)
		}
	}
}

func TestVault_ClientOpt_WithVersion(t *testing.T) {
	// setup tests
	tests := []struct {
		failure bool
		version string
		want    string
	}{
		{
			failure: false,
			version: "1",
			want:    "1",
		},
		{
			failure: false,
			version: "2",
			want:    "2",
		},
		{
			failure: true,
			version: "3",
			want:    "",
		},
		{
			failure: true,
			version: "",
			want:    "",
		},
	}

	// run tests
	for _, test := range tests {
		_service, err := New(
			WithAddress("https://vault.example.com"),
			WithVersion(test.version),
		)

		if test.failure {
			if err == nil {
				t.Errorf("WithVersion should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("WithVersion returned err: %v", err)
		}

		if !reflect.DeepEqual(_service.config.Version, test.want) {
			t.Errorf("WithVersion is %v, want %v", _service.config.Version, test.want)
		}
	}
}
