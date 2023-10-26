// SPDX-License-Identifier: Apache-2.0

package database

import (
	"testing"
	"time"
)

func TestDatabase_Config_Validate(t *testing.T) {
	// setup tests
	tests := []struct {
		failure bool
		name    string
		config  *config
	}{
		{
			name:    "success with postgres",
			failure: false,
			config: &config{
				Driver:           "postgres",
				Address:          "postgres://foo:bar@localhost:5432/vela",
				CompressionLevel: 3,
				ConnectionLife:   10 * time.Second,
				ConnectionIdle:   5,
				ConnectionOpen:   20,
				EncryptionKey:    "A1B2C3D4E5G6H7I8J9K0LMNOPQRSTUVW",
				SkipCreation:     false,
			},
		},
		{
			name:    "success with sqlite3",
			failure: false,
			config: &config{
				Driver:           "sqlite3",
				Address:          "file::memory:?cache=shared",
				CompressionLevel: 3,
				ConnectionLife:   10 * time.Second,
				ConnectionIdle:   5,
				ConnectionOpen:   20,
				EncryptionKey:    "A1B2C3D4E5G6H7I8J9K0LMNOPQRSTUVW",
				SkipCreation:     false,
			},
		},
		{
			name:    "success with negative compression level",
			failure: false,
			config: &config{
				Driver:           "postgres",
				Address:          "postgres://foo:bar@localhost:5432/vela",
				CompressionLevel: -1,
				ConnectionLife:   10 * time.Second,
				ConnectionIdle:   5,
				ConnectionOpen:   20,
				EncryptionKey:    "A1B2C3D4E5G6H7I8J9K0LMNOPQRSTUVW",
				SkipCreation:     false,
			},
		},
		{
			name:    "failure with empty driver",
			failure: true,
			config: &config{
				Driver:           "",
				Address:          "postgres://foo:bar@localhost:5432/vela",
				CompressionLevel: 3,
				ConnectionLife:   10 * time.Second,
				ConnectionIdle:   5,
				ConnectionOpen:   20,
				EncryptionKey:    "A1B2C3D4E5G6H7I8J9K0LMNOPQRSTUVW",
				SkipCreation:     false,
			},
		},
		{
			name:    "failure with empty address",
			failure: true,
			config: &config{
				Driver:           "postgres",
				Address:          "",
				CompressionLevel: 3,
				ConnectionLife:   10 * time.Second,
				ConnectionIdle:   5,
				ConnectionOpen:   20,
				EncryptionKey:    "A1B2C3D4E5G6H7I8J9K0LMNOPQRSTUVW",
				SkipCreation:     false,
			},
		},
		{
			name:    "failure with invalid address",
			failure: true,
			config: &config{
				Driver:           "postgres",
				Address:          "postgres://foo:bar@localhost:5432/vela/",
				CompressionLevel: 3,
				ConnectionLife:   10 * time.Second,
				ConnectionIdle:   5,
				ConnectionOpen:   20,
				EncryptionKey:    "A1B2C3D4E5G6H7I8J9K0LMNOPQRSTUVW",
				SkipCreation:     false,
			},
		},
		{
			name:    "failure with invalid compression level",
			failure: true,
			config: &config{
				Driver:           "postgres",
				Address:          "postgres://foo:bar@localhost:5432/vela",
				CompressionLevel: 10,
				ConnectionLife:   10 * time.Second,
				ConnectionIdle:   5,
				ConnectionOpen:   20,
				EncryptionKey:    "A1B2C3D4E5G6H7I8J9K0LMNOPQRSTUVW",
				SkipCreation:     false,
			},
		},
		{
			name:    "failure with empty encryption key",
			failure: true,
			config: &config{
				Driver:           "postgres",
				Address:          "postgres://foo:bar@localhost:5432/vela",
				CompressionLevel: 3,
				ConnectionLife:   10 * time.Second,
				ConnectionIdle:   5,
				ConnectionOpen:   20,
				EncryptionKey:    "",
				SkipCreation:     false,
			},
		},
		{
			name:    "failure with invalid encryption key",
			failure: true,
			config: &config{
				Driver:           "postgres",
				Address:          "postgres://foo:bar@localhost:5432/vela",
				CompressionLevel: 3,
				ConnectionLife:   10 * time.Second,
				ConnectionIdle:   5,
				ConnectionOpen:   20,
				EncryptionKey:    "A1B2C3D4E5G6H7I8J9K0",
				SkipCreation:     false,
			},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.config.Validate()

			if test.failure {
				if err == nil {
					t.Errorf("Validate for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("Validate for %s returned err: %v", test.name, err)
			}
		})
	}
}
