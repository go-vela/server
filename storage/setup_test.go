package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetup_Minio(t *testing.T) {
	setup := &Setup{
		Enable:    true,
		Driver:    "minio",
		Endpoint:  "minio.example.com",
		AccessKey: "access-key",
		SecretKey: "secret-key",
		Bucket:    "bucket-name",
		Secure:    true,
	}

	storage, err := setup.Minio()
	assert.NoError(t, err)
	assert.NotNil(t, storage)
}

func TestSetup_Validate(t *testing.T) {
	tests := []struct {
		failure bool
		setup   *Setup
	}{
		{
			failure: false,
			setup: &Setup{
				Enable:    true,
				Endpoint:  "example.com",
				AccessKey: "access-key",
				SecretKey: "secret-key",
				Bucket:    "bucket-name",
			},
		},
		{
			failure: true,
			setup: &Setup{
				Enable: false,
			},
		},
		{
			failure: true,
			setup: &Setup{
				Enable:    true,
				AccessKey: "access-key",
				SecretKey: "secret-key",
				Bucket:    "bucket-name",
			},
		},
		{
			failure: true,
			setup: &Setup{
				Enable:    true,
				Endpoint:  "example.com",
				SecretKey: "secret-key",
				Bucket:    "bucket-name",
			},
		},
		{
			failure: true,
			setup: &Setup{
				Enable:    true,
				Endpoint:  "example.com",
				AccessKey: "access-key",
				Bucket:    "bucket-name",
			},
		},
		{
			failure: true,
			setup: &Setup{
				Enable:    true,
				Endpoint:  "example.com",
				AccessKey: "access-key",
				SecretKey: "secret-key",
			},
		},
	}

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
