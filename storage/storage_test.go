// SPDX-License-Identifier: Apache-2.0

package storage

import (
	"testing"

	"github.com/go-vela/server/constants"
)

func TestStorage_New(t *testing.T) {
	tests := []struct {
		name    string
		failure bool
		setup   *Setup
	}{
		{
			name:    "valid-minio-config",
			failure: false,
			setup: &Setup{
				Driver:    constants.DriverMinio,
				Enable:    true,
				Endpoint:  "http://minio.example.com",
				AccessKey: "access-key",
				SecretKey: "secret-key",
				Bucket:    "bucket-name",
				Secure:    true,
			},
		},
		{
			name:    "invalid-driver",
			failure: true,
			setup: &Setup{
				Driver:    "invalid-driver",
				Enable:    true,
				Endpoint:  "http://invalid.example.com",
				AccessKey: "access-key",
				SecretKey: "secret-key",
				Bucket:    "bucket-name",
				Secure:    true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.setup)

			if tt.failure {
				if err == nil {
					t.Errorf("New() expected error, got nil")
				}

				return
			}

			// success case
			if err != nil {
				t.Errorf("New() unexpected error: %v", err)
			}
		})
	}
}
