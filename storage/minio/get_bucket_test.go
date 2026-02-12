// SPDX-License-Identifier: Apache-2.0

package minio

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMinioClient_GetBucket_ReturnsConfiguredBucket(t *testing.T) {
	gin.SetMode(gin.TestMode)

	_, engine := gin.CreateTestContext(httptest.NewRecorder())

	fake := httptest.NewServer(engine)
	defer fake.Close()

	client, err := NewTest(fake.URL, "miniokey", "miniosecret", "foo", false)
	if err != nil {
		t.Fatalf("failed to create minio test client: %v", err)
	}

	got := client.GetBucket()
	want := "foo"

	if got != want {
		t.Fatalf("GetBucket() = %q, want %q", got, want)
	}
}

func TestMinioClient_GetBucket_EmptyWhenUnset(t *testing.T) {
	gin.SetMode(gin.TestMode)

	_, engine := gin.CreateTestContext(httptest.NewRecorder())

	fake := httptest.NewServer(engine)
	defer fake.Close()

	client, err := NewTest(fake.URL, "miniokey", "miniosecret", "foo", false)
	if err != nil {
		t.Fatalf("failed to create minio test client: %v", err)
	}

	client.config.Bucket = ""

	got := client.GetBucket()
	if got != "" {
		t.Fatalf("GetBucket() = %q, want empty string", got)
	}
}
