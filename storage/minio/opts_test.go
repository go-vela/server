// SPDX-License-Identifier: Apache-2.0

package minio

import (
	"testing"
)

func TestWithAccessKey(t *testing.T) {
	// setup tests
	tests := []struct {
		failure   bool
		accessKey string
		want      string
	}{
		{
			failure:   false,
			accessKey: "validAccessKey",
			want:      "validAccessKey",
		},
		{
			failure:   true,
			accessKey: "",
			want:      "",
		},
	}

	// run tests
	for _, test := range tests {
		client, err := NewTest("https://minio.example.com",
			test.accessKey,
			"miniosecret",
			"foo",
			false)

		if test.failure {
			if err == nil {
				t.Errorf("WithAddress should have returned err")
			}

			continue
		}

		if err != nil && test.accessKey != "" {
			t.Errorf("WithAccessKey returned err: %v", err)
		}

		if client.config.AccessKey != test.want {
			t.Errorf("WithAccessKey is %v, want %v", client.config.AccessKey, test.want)
		}
	}
}

func TestWithSecretKey(t *testing.T) {
	// setup tests
	tests := []struct {
		failure   bool
		secretKey string
		want      string
	}{
		{
			failure:   false,
			secretKey: "validSecretKey",
			want:      "validSecretKey",
		},
		{
			failure:   true,
			secretKey: "",
			want:      "",
		},
	}

	// run tests
	for _, test := range tests {
		client, err := NewTest("https://minio.example.com",
			"minioaccess",
			test.secretKey,
			"foo",
			false)

		if test.failure {
			if err == nil {
				t.Errorf("WithSecretKey should have returned err")
			}

			continue
		}

		if err != nil && test.secretKey != "" {
			t.Errorf("WithSecretKey returned err: %v", err)
		}

		if client.config.SecretKey != test.want {
			t.Errorf("WithSecretKey is %v, want %v", client.config.SecretKey, test.want)
		}
	}
}

func TestWithSecure(t *testing.T) {
	// setup tests
	tests := []struct {
		failure bool
		secure  bool
		want    bool
	}{
		{
			failure: false,
			secure:  true,
			want:    true,
		},
		{
			failure: false,
			secure:  false,
			want:    false,
		},
	}

	// run tests
	for _, test := range tests {
		client, err := NewTest("https://minio.example.com",
			"minioaccess",
			"miniosecret",
			"foo",
			test.secure)

		if test.failure {
			if err == nil {
				t.Errorf("WithSecure should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("WithSecure returned err: %v", err)
		}

		if client.config.Secure != test.want {
			t.Errorf("WithSecure is %v, want %v", client.config.Secure, test.want)
		}
	}
}

func TestWithBucket(t *testing.T) {
	// setup tests
	tests := []struct {
		failure bool
		bucket  string
		want    string
	}{
		{
			failure: false,
			bucket:  "validBucket",
			want:    "validBucket",
		},
		{
			failure: true,
			bucket:  "",
			want:    "",
		},
	}

	// run tests
	for _, test := range tests {
		client, err := NewTest("https://minio.example.com",
			"minioaccess",
			"miniosecret",
			test.bucket,
			false)

		if test.failure {
			if err == nil {
				t.Errorf("WithBucket should have returned err")
			}

			continue
		}

		if err != nil && test.bucket != "" {
			t.Errorf("WithBucket returned err: %v", err)
		}

		if client.config.Bucket != test.want {
			t.Errorf("WithBucket is %v, want %v", client.config.Bucket, test.want)
		}
	}
}
