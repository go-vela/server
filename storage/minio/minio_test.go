package minio

import (
	"testing"
)

var (
	endpoint   = "http://localhost:9000"
	_accessKey = "minio_access_user"
	_secretKey = "minio_secret_key"
	_useSSL    = false
)

func TestMinio_New(t *testing.T) {
	// setup types
	// create a local fake MinIO instance
	//
	// https://pkg.go.dev/github.com/minio/minio-go/v7#New
	// setup tests
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
			WithAccessKey(_accessKey),
			WithSecretKey(_secretKey),
			WithSecure(_useSSL),
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
