// SPDX-License-Identifier: Apache-2.0

package minio

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	api "github.com/go-vela/server/api/types"
)

func TestMinioClient_ListBuckets_Success(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)
	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	engine.PUT("/foo/", func(c *gin.Context) {
		c.Header("Content-Type", "application/xml;charset=utf-8")
		c.Status(http.StatusOK)
	})
	// mock list buckets call
	engine.GET("/", func(c *gin.Context) {
		c.Header("Content-Type", "application/xml;charset=utf-8")
		c.XML(http.StatusOK, gin.H{
			"Name": "foo",
		})
	})
	fake := httptest.NewServer(engine)
	defer fake.Close()
	ctx := context.TODO()
	client, _ := NewTest(fake.URL, "miniokey", "miniosecret", "foo", false)
	b := new(api.Bucket)
	b.BucketName = "foo"

	// run test
	err := client.CreateBucket(ctx, b)
	if resp.Code != http.StatusOK {
		t.Errorf("CreateBucket returned %v, want %v", resp.Code, http.StatusOK)
	}
	if err != nil {
		t.Errorf("CreateBucket returned err: %v", err)
	}

	buckets, err := client.ListBuckets(ctx)
	if err != nil {
		t.Errorf("ListBuckets returned err: %v", err)
	}

	// check if buckets are correct
	expectedBuckets := []string{"foo"}
	for i, bucket := range buckets {
		if bucket != expectedBuckets[i] {
			t.Errorf("Expected bucket %v, got %v", expectedBuckets[i], bucket)
		}
	}
}

func TestMinioClient_ListBuckets_Failure(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// mock list buckets call
	engine.GET("/minio/buckets", func(c *gin.Context) {
		c.Status(http.StatusInternalServerError)
	})
	fake := httptest.NewServer(engine)
	defer fake.Close()
	ctx := context.TODO()
	client, _ := NewTest(fake.URL, "miniokey", "miniosecret", "foo", false)

	// run test
	_, err := client.ListBuckets(ctx)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}
