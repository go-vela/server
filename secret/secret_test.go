// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package secret

import (
	"testing"

	"github.com/go-vela/server/database"
)

func TestSecret_New(t *testing.T) {
	// setup types
	db, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create test database engine: %v", err)
	}
	defer db.Close()

	// setup tests
	tests := []struct {
		failure bool
		setup   *Setup
	}{
		{
			failure: false,
			setup: &Setup{
				Driver:   "native",
				Database: db,
			},
		},
		{
			failure: false,
			setup: &Setup{
				Driver:        "vault",
				Address:       "https://vault.example.com",
				AuthMethod:    "",
				AwsRole:       "",
				Prefix:        "bar",
				Token:         "baz",
				TokenDuration: 0,
				Version:       "1",
			},
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
