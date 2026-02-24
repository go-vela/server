// SPDX-License-Identifier: Apache-2.0

package minio

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMinioClient_GetAddress_ReturnsConfiguredBucket(t *testing.T) {
	gin.SetMode(gin.TestMode)

	_, engine := gin.CreateTestContext(httptest.NewRecorder())

	fake := httptest.NewServer(engine)
	defer fake.Close()

	client, err := NewTest(fake.URL, "miniokey", "miniosecret", "foo", false)
	if err != nil {
		t.Fatalf("failed to create minio test client: %v", err)
	}

	got := client.GetAddress()
	want := fake.Listener.Addr().String()

	if got != want {
		t.Fatalf("GetAddress() = %q, want %q", got, want)
	}
}

func TestMinioClient_GetAddress_EmptyWhenUnset(t *testing.T) {
	gin.SetMode(gin.TestMode)

	_, engine := gin.CreateTestContext(httptest.NewRecorder())

	fake := httptest.NewServer(engine)
	defer fake.Close()

	client, err := NewTest(fake.URL, "miniokey", "miniosecret", "foo", false)
	if err != nil {
		t.Fatalf("failed to create minio test client: %v", err)
	}

	client.config.Endpoint = ""

	got := client.GetAddress()
	if got != "" {
		t.Fatalf("GetAddress() = %q, want empty string", got)
	}
}
