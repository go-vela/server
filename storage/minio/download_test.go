// SPDX-License-Identifier: Apache-2.0

package minio

import (
	"context"
	"github.com/gin-gonic/gin"
	api "github.com/go-vela/server/api/types"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMinioClient_Download_Success(t *testing.T) {
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

	// mock upload call
	engine.PUT("/foo/test.xml", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("test_data/test.xml")
	})

	// mock download call
	engine.GET("/foo/test.xml", func(c *gin.Context) {
		c.Header("Content-Type", "text/xml")
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
	// upload test file
	err = client.Upload(ctx, obj)
	if resp.Code != http.StatusOK {
		t.Errorf("Upload returned %v, want %v", resp.Code, http.StatusOK)
	}
	// run test
	err = client.Download(ctx, obj)
	if err != nil {
		t.Errorf("Download returned err: %v", err)
	}

	// check if file exists
	if _, err := os.Stat(obj.FilePath); os.IsNotExist(err) {
		t.Errorf("Downloaded file does not exist: %v", obj.FilePath)
	}
}

func TestMinioClient_Download_Failure(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)
	// mock create bucket call
	engine.PUT("/foo/", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// setup mock server
	engine.GET("/foo/test.xml", func(c *gin.Context) {
		c.Header("Content-Type", "application/octet-stream")
		c.Status(http.StatusNotFound)
	})
	fake := httptest.NewServer(engine)
	defer fake.Close()
	ctx := context.TODO()
	obj := new(api.Object)
	obj.Bucket.BucketName = "foo"
	obj.ObjectName = "test.xml"
	obj.FilePath = "testdata/tests.xml"
	client, _ := NewTest(fake.URL, "miniokey", "miniosecret", "foo", false)

	// create bucket
	//err := client.CreateBucket(ctx, &api.Bucket{BucketName: "foo"})
	//if err != nil {
	//	t.Errorf("CreateBucket returned err: %v", err)
	//}

	// run test
	err := client.Download(ctx, obj)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}
