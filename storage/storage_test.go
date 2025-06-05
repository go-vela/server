// SPDX-License-Identifier: Apache-2.0

package storage

import (
	"testing"

	"github.com/go-vela/server/constants"
)

func TestStorage_New(t *testing.T) {
	tests := []struct {
		failure bool
		setup   *Setup
	}{
		{
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
			failure: true,
			setup: &Setup{
				Driver:    "invalid-driver",
				Enable:    false,
				Endpoint:  "http://invalid.example.com",
				AccessKey: "access-key",
				SecretKey: "secret-key",
				Bucket:    "bucket-name",
				Secure:    true,
			},
		},
	}

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
