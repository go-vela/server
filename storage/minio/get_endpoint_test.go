// SPDX-License-Identifier: Apache-2.0

package minio

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMinioClient_GetEndpoint_ReturnsConfiguredBucket(t *testing.T) {
	gin.SetMode(gin.TestMode)

	_, engine := gin.CreateTestContext(httptest.NewRecorder())

	fake := httptest.NewServer(engine)
	defer fake.Close()

	client, err := NewTest(fake.URL, "miniokey", "miniosecret", "foo", false)
	if err != nil {
		t.Fatalf("failed to create minio test client: %v", err)
	}

	got := client.GetEndpoint()
	want := fake.URL

	if got != want {
		t.Fatalf("GetAddress() = %q, want %q", got, want)
	}
}

func TestMinioClient_GetEndpoint_EmptyWhenUnset(t *testing.T) {
	gin.SetMode(gin.TestMode)

	_, engine := gin.CreateTestContext(httptest.NewRecorder())

	fake := httptest.NewServer(engine)
	defer fake.Close()

	client, err := NewTest(fake.URL, "miniokey", "miniosecret", "foo", false)
	if err != nil {
		t.Fatalf("failed to create minio test client: %v", err)
	}

	client.config.Endpoint = ""

	got := client.GetEndpoint()
	if got != "" {
		t.Fatalf("GetAddress() = %q, want empty string", got)
	}
}
