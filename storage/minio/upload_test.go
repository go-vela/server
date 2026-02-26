// SPDX-License-Identifier: Apache-2.0

package minio

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	api "github.com/go-vela/server/api/types"
)

func TestMinioClient_UploadObject_Success(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// mock bucket location check (required by minio-go before upload)
	engine.GET("/foo/", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/xml", []byte(`<LocationConstraint>us-east-1</LocationConstraint>`))
	})

	// mock upload (PutObject) call
	engine.PUT("/foo/test.xml", func(c *gin.Context) {
		c.Header("ETag", "\"abc123\"")
		c.Status(http.StatusOK)
	})

	fake := httptest.NewServer(engine)
	defer fake.Close()

	ctx := context.TODO()
	obj := &api.Object{
		ObjectName: "test.xml",
		FilePath:   "test_data/test.xml",
		Bucket: api.Bucket{
			BucketName: "foo",
		},
	}

	content := []byte("<note><body>test</body></note>")
	reader := bytes.NewReader(content)

	client, _ := NewTest(fake.URL, "miniokey", "miniosecret", "foo", false)

	// run test
	err := client.UploadObject(ctx, obj, reader, int64(len(content)))
	if err != nil {
		t.Errorf("UploadObject returned err: %v", err)
	}
}

func TestMinioClient_UploadObject_Failure(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	ctx, engine := gin.CreateTestContext(resp)

	// mock upload (PutObject) call - return internal server error
	engine.PUT("/foo/test.xml", func(c *gin.Context) {
		c.Status(http.StatusInternalServerError)
	})

	fake := httptest.NewServer(engine)
	defer fake.Close()

	obj := &api.Object{
		ObjectName: "test.xml",
		FilePath:   "test_data/test.xml",
		Bucket: api.Bucket{
			BucketName: "foo",
		},
	}

	content := []byte("<note><body>test</body></note>")
	reader := bytes.NewReader(content)

	client, _ := NewTest(fake.URL, "miniokey", "miniosecret", "foo", false)

	// run test - expect error due to server returning 500
	err := client.UploadObject(ctx, obj, reader, int64(len(content)))
	if err == nil {
		t.Errorf("UploadObject should have returned err")
	}
}
