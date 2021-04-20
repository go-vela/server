// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package secret

import (
	"reflect"
	"testing"

	"github.com/go-vela/server/database"
)

func TestSecret_Setup_Native(t *testing.T) {
	// setup types
	_database, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create database service: %v", err)
	}
	defer _database.Database.Close()

	_setup := &Setup{
		Driver:   "native",
		Database: _database,
	}

	_native, err := _setup.Native()
	if err != nil {
		t.Errorf("unable to setup secret service: %v", err)
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
			want:    _native,
		},
		{
			failure: true,
			setup:   &Setup{Driver: "native"},
			want:    nil,
		},
	}

	// run tests
	for _, test := range tests {
		got, err := test.setup.Native()

		if test.failure {
			if err == nil {
				t.Errorf("Native should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("Native returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("Native is %v, want %v", got, test.want)
		}
	}
}

func TestSecret_Setup_Vault(t *testing.T) {
	// setup types
	_setup := &Setup{
		Driver:        "vault",
		Address:       "https://vault.example.com",
		AuthMethod:    "",
		AwsRole:       "",
		Prefix:        "bar",
		Token:         "baz",
		TokenDuration: 0,
		Version:       "1",
	}

	_vault, err := _setup.Vault()
	if err != nil {
		t.Errorf("unable to setup secret service: %v", err)
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
			want:    _vault,
		},
		{
			failure: true,
			setup:   &Setup{Driver: "vault"},
			want:    nil,
		},
	}

	// run tests
	for _, test := range tests {
		_, err := test.setup.Vault()

		if test.failure {
			if err == nil {
				t.Errorf("Vault should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("Vault returned err: %v", err)
		}
	}
}

func TestSecret_Setup_Validate(t *testing.T) {
	// setup types
	_database, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create database service: %v", err)
	}
	defer _database.Database.Close()

	// setup tests
	tests := []struct {
		failure bool
		setup   *Setup
	}{
		{
			failure: false,
			setup: &Setup{
				Driver:   "native",
				Database: _database,
			},
		},
		{
			failure: false,
			setup: &Setup{
				Driver:        "vault",
				Address:       "https://vault.example.com",
				AuthMethod:    "aws",
				AwsRole:       "foo",
				Prefix:        "bar",
				Token:         "baz",
				TokenDuration: 0,
				Version:       "1",
			},
		},
		{
			failure: true,
			setup: &Setup{
				Driver: "",
			},
		},
		{
			failure: true,
			setup: &Setup{
				Driver:   "native",
				Database: nil,
			},
		},
		{
			failure: true,
			setup: &Setup{
				Driver:        "vault",
				Address:       "https://vault.example.com/",
				AuthMethod:    "aws",
				AwsRole:       "foo",
				Prefix:        "bar",
				Token:         "baz",
				TokenDuration: 0,
				Version:       "1",
			},
		},
		{
			failure: true,
			setup: &Setup{
				Driver:        "vault",
				Address:       "vault.example.com",
				AuthMethod:    "aws",
				AwsRole:       "foo",
				Prefix:        "bar",
				Token:         "baz",
				TokenDuration: 0,
				Version:       "1",
			},
		},
		{
			failure: true,
			setup: &Setup{
				Driver:        "vault",
				Address:       "",
				AuthMethod:    "aws",
				AwsRole:       "foo",
				Prefix:        "bar",
				Token:         "baz",
				TokenDuration: 0,
				Version:       "1",
			},
		},
		{
			failure: true,
			setup: &Setup{
				Driver:        "vault",
				Address:       "https://vault.example.com",
				AuthMethod:    "",
				AwsRole:       "foo",
				Prefix:        "bar",
				Token:         "",
				TokenDuration: 0,
				Version:       "1",
			},
		},
		{
			failure: true,
			setup: &Setup{
				Driver:        "vault",
				Address:       "https://vault.example.com",
				AuthMethod:    "aws",
				AwsRole:       "",
				Prefix:        "bar",
				Token:         "",
				TokenDuration: 0,
				Version:       "1",
			},
		},
		{
			failure: true,
			setup: &Setup{
				Driver:        "vault",
				Address:       "https://vault.example.com",
				AuthMethod:    "ldap",
				AwsRole:       "foo",
				Prefix:        "bar",
				Token:         "",
				TokenDuration: 0,
				Version:       "1",
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
