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

func TestMinioClient_BucketExists(t *testing.T) {
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
	engine.HEAD("/foo/", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	fake := httptest.NewServer(engine)
	defer fake.Close()

	ctx := context.TODO()

	client, _ := NewTest(fake.URL, "miniokey", "miniosecret", "foo", false)

	// create bucket
	err := client.CreateBucket(ctx, &api.Bucket{BucketName: "foo"})
	if err != nil {
		t.Errorf("CreateBucket returned err: %v", err)
	}

	// run test
	exists, err := client.BucketExists(ctx, &api.Bucket{BucketName: "foo"})

	if resp.Code != http.StatusOK {
		t.Errorf("BucketExists returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("BucketExists returned err: %v", err)
	}

	if !exists {
		t.Errorf("BucketExists returned %v, want %v", exists, true)
	}
}

func TestMinioClient_BucketExists_Failure(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.HEAD("/foo/", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	fake := httptest.NewServer(engine)
	defer fake.Close()

	ctx := context.TODO()

	client, _ := NewTest(fake.URL, "miniokey", "miniosecret", "foo", false)

	// run test
	exists, err := client.BucketExists(ctx, &api.Bucket{BucketName: "bar"})

	if resp.Code != http.StatusOK {
		t.Errorf("BucketExists returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("BucketExists returned err: %v", err)
	}

	if exists {
		t.Errorf("BucketExists returned %v, want %v", exists, false)
	}
}
