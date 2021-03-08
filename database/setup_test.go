// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package database

import (
	"testing"
	"time"
)

func TestDatabase_Setup_Postgres(t *testing.T) {
	// setup types
	_setup := &Setup{
		Driver:           "postgres",
		Address:          "postgres://postgres:5432/vela?sslmode=disable",
		CompressionLevel: 3,
		ConnectionIdle:   2,
		ConnectionLife:   30 * time.Minute,
		ConnectionOpen:   0,
		EncryptionKey:    "C639A572E14D5075C526FDDD43E4ECF6",
	}

	got, err := _setup.Postgres()
	if err == nil {
		t.Errorf("Postgres should have returned err")
	}

	if got != nil {
		t.Errorf("Postgres is %v, want nil", got)
	}
}

func TestDatabase_Setup_Sqlite(t *testing.T) {
	// setup types
	_setup := &Setup{
		Driver:           "sqlite",
		Address:          ":memory:",
		CompressionLevel: 3,
		ConnectionIdle:   2,
		ConnectionLife:   30 * time.Minute,
		ConnectionOpen:   0,
		EncryptionKey:    "C639A572E14D5075C526FDDD43E4ECF6",
	}

	got, err := _setup.Sqlite()
	if err == nil {
		t.Errorf("Sqlite should have returned err")
	}

	if got != nil {
		t.Errorf("Sqlite is %v, want nil", got)
	}
}

func TestDatabase_Setup_Validate(t *testing.T) {
	// setup tests
	tests := []struct {
		failure bool
		setup   *Setup
	}{
		{
			failure: false,
			setup: &Setup{
				Driver:           "postgres",
				Address:          "postgres://postgres:5432/vela?sslmode=disable",
				CompressionLevel: 3,
				ConnectionIdle:   2,
				ConnectionLife:   30 * time.Minute,
				ConnectionOpen:   0,
				EncryptionKey:    "C639A572E14D5075C526FDDD43E4ECF6",
			},
		},
		{
			failure: false,
			setup: &Setup{
				Driver:           "sqlite",
				Address:          ":memory:",
				CompressionLevel: 3,
				ConnectionIdle:   2,
				ConnectionLife:   30 * time.Minute,
				ConnectionOpen:   0,
				EncryptionKey:    "C639A572E14D5075C526FDDD43E4ECF6",
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
