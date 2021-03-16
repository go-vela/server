// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package secret

import (
	"reflect"
	"testing"

	"github.com/go-vela/server/database"
)

func TestSource_New(t *testing.T) {
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
		want    Service
	}{
		{
			failure: true,
			setup: &Setup{
				Driver:   "native",
				Database: _database,
			},
			want: nil,
		},
		{
			failure: true,
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
			want: nil,
		},
		{
			failure: true,
			setup: &Setup{
				Driver:   "native",
				Database: nil,
			},
			want: nil,
		},
		{
			failure: true,
			setup: &Setup{
				Driver:        "kubernetes",
				Address:       "https://kubernetes.example.com",
				AuthMethod:    "aws",
				AwsRole:       "foo",
				Prefix:        "bar",
				Token:         "baz",
				TokenDuration: 0,
				Version:       "1",
			},
			want: nil,
		},
	}

	// run tests
	for _, test := range tests {
		got, err := New(test.setup)

		if test.failure {
			if err == nil {
				t.Errorf("New should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("New returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("New is %v, want %v", got, test.want)
		}
	}
}
