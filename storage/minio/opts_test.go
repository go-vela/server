package minio

import (
	"reflect"
	"testing"
)

func TestWithAccessKey(t *testing.T) {
	// setup tests
	tests := []struct {
		accessKey string
		want      string
	}{
		{
			accessKey: "validAccessKey",
			want:      "validAccessKey",
		},
		{
			accessKey: "",
			want:      "",
		},
	}

	// run tests
	for _, test := range tests {
		client, err := New("http://localhost:8080", WithAccessKey(test.accessKey))

		if err != nil && test.accessKey != "" {
			t.Errorf("WithAccessKey returned err: %v", err)
		}

		if !reflect.DeepEqual(client.config.AccessKey, test.want) {
			t.Errorf("WithAccessKey is %v, want %v", client.config.AccessKey, test.want)
		}
	}
}

func TestWithSecretKey(t *testing.T) {
	tests := []struct {
		secretKey string
		wantErr   bool
	}{
		{"validSecretKey", false},
		{"", true},
	}

	for _, test := range tests {
		client := &Client{config: &config{}}
		err := WithSecretKey(test.secretKey)(client)
		if (err != nil) != test.wantErr {
			t.Errorf("WithSecretKey() error = %v, wantErr %v", err, test.wantErr)
		}
		if !test.wantErr && client.config.SecretKey != test.secretKey {
			t.Errorf("WithSecretKey() = %v, want %v", client.config.SecretKey, test.secretKey)
		}
	}
}

func TestWithSecure(t *testing.T) {
	tests := []struct {
		secure bool
	}{
		{true},
		{false},
	}

	for _, test := range tests {
		client := &Client{config: &config{}}
		err := WithSecure(test.secure)(client)
		if err != nil {
			t.Errorf("WithSecure() error = %v", err)
		}
		if client.config.Secure != test.secure {
			t.Errorf("WithSecure() = %v, want %v", client.config.Secure, test.secure)
		}
	}
}

func TestWithBucket(t *testing.T) {
	tests := []struct {
		bucket  string
		wantErr bool
	}{
		{"validBucket", false},
		{"", true},
	}

	for _, test := range tests {
		client := &Client{config: &config{}}
		err := WithBucket(test.bucket)(client)
		if (err != nil) != test.wantErr {
			t.Errorf("WithBucket() error = %v, wantErr %v", err, test.wantErr)
		}
		if !test.wantErr && client.config.Bucket != test.bucket {
			t.Errorf("WithBucket() = %v, want %v", client.config.Bucket, test.bucket)
		}
	}
}
