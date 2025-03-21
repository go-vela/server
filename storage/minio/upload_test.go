// SPDX-License-Identifier: Apache-2.0

package minio

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"

	api "github.com/go-vela/server/api/types"
)

func TestMinioClient_Upload_Success(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	// mock create bucket call
	engine.PUT("/foo/", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
	})
	// mock bucket exists call
	engine.PUT("/foo/test.xml", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("test_data/test.xml")
	})
	fake := httptest.NewServer(engine)
	defer fake.Close()
	ctx := context.TODO()
	obj := new(api.Object)
	obj.Bucket.BucketName = "foo"
	obj.ObjectName = "test.xml"
	obj.FilePath = "test_data/test.xml"
	client, _ := NewTest(fake.URL, "miniokey", "miniosecret", "foo", false)

	// create bucket
	err := client.CreateBucket(ctx, &api.Bucket{BucketName: "foo"})
	if err != nil {
		t.Errorf("CreateBucket returned err: %v", err)
	}

	// run test
	err = client.Upload(ctx, obj)
	if resp.Code != http.StatusOK {
		t.Errorf("Upload returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("Upload returned err: %v", err)
	}
}

func TestMinioClient_Upload_Failure(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	// mock create bucket call
	engine.PUT("/foo/", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
	})
	// mock bucket exists call
	engine.PUT("/foo/test.xml", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("test_data/test.xml")
	})
	fake := httptest.NewServer(engine)
	defer fake.Close()
	ctx := context.TODO()
	obj := new(api.Object)
	obj.Bucket.BucketName = "foo"
	obj.ObjectName = "test.xml"
	obj.FilePath = "nonexist/test.xml"
	client, _ := NewTest(fake.URL, "miniokey", "miniosecret", "foo", false)

	// create bucket
	err := client.CreateBucket(ctx, &api.Bucket{BucketName: "foo"})
	if err != nil {
		t.Errorf("CreateBucket returned err: %v", err)
	}

	// run test
	err = client.Upload(ctx, obj)
	if resp.Code != http.StatusOK {
		t.Errorf("Upload returned %v, want %v", resp.Code, http.StatusOK)
	}

	if !os.IsNotExist(err) {
		t.Errorf("Upload returned err: %v", err)
	}
}
