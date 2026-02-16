// SPDX-License-Identifier: Apache-2.0

package minio

import (
	"testing"
)

var (
	endpoint   = "http://localhost:9000"
	_accessKey = "minio_access_user"
	_secretKey = "minio_secret_key"
	_bucket    = "minio_bucket"
	_useSSL    = false
)

func TestMinio_New(t *testing.T) {
	tests := []struct {
		failure  bool
		endpoint string
	}{
		{
			failure:  false,
			endpoint: endpoint,
		},
		{
			failure:  true,
			endpoint: "",
		},
	}

	// run tests
	for _, test := range tests {
		_, err := New(
			test.endpoint,
			WithOptions(true, _useSSL,
				test.endpoint, _accessKey, _secretKey, _bucket, ""),
		)

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
