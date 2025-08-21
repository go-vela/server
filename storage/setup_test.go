// SPDX-License-Identifier: Apache-2.0

package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/go-vela/server/constants"
)

func TestSetup_Minio(t *testing.T) {
	setup := &Setup{
		Enable:    true,
		Driver:    constants.DriverMinio,
		Endpoint:  "http://minio.example.com",
		AccessKey: "access-key",
		SecretKey: "secret-key",
		Bucket:    "bucket-name",
		Secure:    true,
	}

	storageClient, err := setup.Minio()
	assert.NoError(t, err)
	assert.NotNil(t, storageClient)
}

func TestSetup_Validate(t *testing.T) {
	tests := []struct {
		name    string
		setup   *Setup
		wantErr bool
	}{
		{
			name: "storage disabled",
			setup: &Setup{
				Enable: false,
			},
			wantErr: false,
		},
		{
			name: "valid config",
			setup: &Setup{
				Enable:    true,
				Driver:    constants.DriverMinio,
				Endpoint:  "http://example.com",
				AccessKey: "access-key",
				SecretKey: "secret-key",
				Bucket:    "bucket-name",
			},
			wantErr: false,
		},
		{
			name: "missing bucket",
			setup: &Setup{
				Enable:    true,
				Endpoint:  "http://example.com",
				AccessKey: "access-key",
				SecretKey: "secret-key",
			},
			wantErr: true,
		},
		{
			name: "driver set",
			setup: &Setup{
				Enable:    true,
				Driver:    constants.DriverMinio,
				Endpoint:  "http://example.com",
				AccessKey: "access-key",
				SecretKey: "secret-key",
				Bucket:    "bucket-name",
			},
			wantErr: false,
		},
		{
			name: "missing endpoint",
			setup: &Setup{
				Enable:    true,
				AccessKey: "access-key",
				SecretKey: "secret-key",
				Bucket:    "bucket-name",
			},
			wantErr: true,
		},
		{
			name: "missing credentials",
			setup: &Setup{
				Enable:   true,
				Endpoint: "http://example.com",
				Bucket:   "bucket-name",
			},
			wantErr: true,
		},
		{
			name: "invalid endpoint URL",
			setup: &Setup{
				Enable:    true,
				Endpoint:  "://bad-url",
				AccessKey: "access-key",
				SecretKey: "secret-key",
				Bucket:    "bucket-name",
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.setup.Validate()
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
