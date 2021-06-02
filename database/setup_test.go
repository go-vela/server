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
		Address:          "postgres://foo:bar@localhost:5432/vela",
		CompressionLevel: 3,
		ConnectionLife:   10 * time.Second,
		ConnectionIdle:   5,
		ConnectionOpen:   20,
		EncryptionKey:    "A1B2C3D4E5G6H7I8J9K0LMNOPQRSTUVW",
	}

	// setup tests
	tests := []struct {
		failure bool
		setup   *Setup
	}{
		{
			failure: true,
			setup:   _setup,
		},
		{
			failure: true,
			setup:   &Setup{Driver: "postgres"},
		},
	}

	// run tests
	for _, test := range tests {
		_, err := test.setup.Postgres()

		if test.failure {
			if err == nil {
				t.Errorf("Postgres should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("Postgres returned err: %v", err)
		}
	}
}

func TestDatabase_Setup_Sqlite(t *testing.T) {
	// setup types
	_setup := &Setup{
		Driver:           "sqlite3",
		Address:          "file::memory:?cache=shared",
		CompressionLevel: 3,
		ConnectionLife:   10 * time.Second,
		ConnectionIdle:   5,
		ConnectionOpen:   20,
		EncryptionKey:    "A1B2C3D4E5G6H7I8J9K0LMNOPQRSTUVW",
	}

	// setup tests
	tests := []struct {
		failure bool
		setup   *Setup
	}{
		{
			failure: false,
			setup:   _setup,
		},
		{
			failure: true,
			setup:   &Setup{Driver: "sqlite3"},
		},
	}

	// run tests
	for _, test := range tests {
		_, err := test.setup.Sqlite()

		if test.failure {
			if err == nil {
				t.Errorf("Sqlite should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("Sqlite returned err: %v", err)
		}
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
				Address:          "postgres://foo:bar@localhost:5432/vela",
				CompressionLevel: 3,
				ConnectionLife:   10 * time.Second,
				ConnectionIdle:   5,
				ConnectionOpen:   20,
				EncryptionKey:    "A1B2C3D4E5G6H7I8J9K0LMNOPQRSTUVW",
			},
		},
		{
			failure: false,
			setup: &Setup{
				Driver:           "sqlite3",
				Address:          "file::memory:?cache=shared",
				CompressionLevel: 3,
				ConnectionLife:   10 * time.Second,
				ConnectionIdle:   5,
				ConnectionOpen:   20,
				EncryptionKey:    "A1B2C3D4E5G6H7I8J9K0LMNOPQRSTUVW",
			},
		},
		{
			failure: true,
			setup: &Setup{
				Driver:           "postgres",
				Address:          "postgres://foo:bar@localhost:5432/vela/",
				CompressionLevel: 3,
				ConnectionLife:   10 * time.Second,
				ConnectionIdle:   5,
				ConnectionOpen:   20,
				EncryptionKey:    "A1B2C3D4E5G6H7I8J9K0LMNOPQRSTUVW",
			},
		},
		{
			failure: true,
			setup: &Setup{
				Driver:           "",
				Address:          "postgres://foo:bar@localhost:5432/vela",
				CompressionLevel: 3,
				ConnectionLife:   10 * time.Second,
				ConnectionIdle:   5,
				ConnectionOpen:   20,
				EncryptionKey:    "A1B2C3D4E5G6H7I8J9K0LMNOPQRSTUVW",
			},
		},
		{
			failure: true,
			setup: &Setup{
				Driver:           "postgres",
				Address:          "",
				CompressionLevel: 3,
				ConnectionLife:   10 * time.Second,
				ConnectionIdle:   5,
				ConnectionOpen:   20,
				EncryptionKey:    "A1B2C3D4E5G6H7I8J9K0LMNOPQRSTUVW",
			},
		},
		{
			failure: true,
			setup: &Setup{
				Driver:           "postgres",
				Address:          "postgres://foo:bar@localhost:5432/vela",
				CompressionLevel: 3,
				ConnectionLife:   10 * time.Second,
				ConnectionIdle:   5,
				ConnectionOpen:   20,
				EncryptionKey:    "",
			},
		},
		{
			failure: true,
			setup: &Setup{
				Driver:           "postgres",
				Address:          "postgres://foo:bar@localhost:5432/vela",
				CompressionLevel: 3,
				ConnectionLife:   10 * time.Second,
				ConnectionIdle:   5,
				ConnectionOpen:   20,
				EncryptionKey:    "A1B2C3D4E5G6H7I8J9K0",
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
