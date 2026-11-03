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

func TestWithUseIAM(t *testing.T) {
	// setup tests
	tests := []struct {
		name      string
		useIAM    bool
		accessKey string
		secretKey string
		failure   bool
	}{
		{
			name:      "iam enabled tolerates empty static keys",
			useIAM:    true,
			accessKey: "",
			secretKey: "",
			failure:   false,
		},
		{
			name:      "static mode requires keys",
			useIAM:    false,
			accessKey: "",
			secretKey: "",
			failure:   true,
		},
		{
			name:      "iam enabled with keys still succeeds",
			useIAM:    true,
			accessKey: "minioaccess",
			secretKey: "miniosecret",
			failure:   false,
		},
	}

	// run tests
	for _, test := range tests {
		client, err := New("https://minio.example.com",
			WithOptions(true, false, test.useIAM,
				"https://minio.example.com", test.accessKey, test.secretKey, "foo", "", "minio"))

		if test.failure {
			if err == nil {
				t.Errorf("%s: WithOptions should have returned err", test.name)
			}

			continue
		}

		if err != nil {
			t.Errorf("%s: WithOptions returned err: %v", test.name, err)

			continue
		}

		if client.config.UseIAM != test.useIAM {
			t.Errorf("%s: UseIAM is %v, want %v", test.name, client.config.UseIAM, test.useIAM)
		}
	}
}
