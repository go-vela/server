// SPDX-License-Identifier: Apache-2.0

package minio

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	api "github.com/go-vela/server/api/types"
)

func Test_PresignedPutObject_Success(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	_, engine := gin.CreateTestContext(httptest.NewRecorder())

	// mock bucket location check (required by minio-go before generating presigned URL)
	engine.GET("/foo/", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/xml", []byte(`<LocationConstraint>us-east-1</LocationConstraint>`))
	})

	fake := httptest.NewServer(engine)
	defer fake.Close()

	client, _ := NewTest(fake.URL, "miniokey", "miniosecret", "foo", false)

	object := &api.Object{
		ObjectName: "test.xml",
		Bucket: api.Bucket{
			BucketName: "foo",
		},
	}

	// run test
	url, err := client.PresignedPutObject(t.Context(), object.ObjectName, 1*time.Minute)
	if err != nil {
		t.Errorf("PresignedPutObject returned err: %v", err)
	}

	// check if URL is valid
	if url == "" {
		t.Errorf("PresignedPutObject returned empty URL")
	}
}

func Test_PresignedPutObject_Failure(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	_, engine := gin.CreateTestContext(httptest.NewRecorder())

	fake := httptest.NewServer(engine)
	defer fake.Close()

	client, _ := NewTest(fake.URL, "miniokey", "miniosecret", "foo", false)

	object := &api.Object{
		ObjectName: "test.xml",
		Bucket: api.Bucket{
			BucketName: "foo",
		},
	}

	// run test - pass a negative duration to trigger a validation error
	url, err := client.PresignedPutObject(t.Context(), object.ObjectName, -1*time.Second)
	if err == nil {
		t.Error("PresignedPutObject should have returned error")
	}

	// on error, PresignedPutObject returns a non-empty error message string as the URL
	if url == "" {
		t.Error("PresignedPutObject should return error message as URL on failure")
	}
}
