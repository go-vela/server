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

func Test_PresignedGetObject_Success(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)
	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// mock presigned get object call
	engine.GET("/foo/", func(c *gin.Context) {
		c.Header("Content-Type", "application/xml")
		c.XML(200, gin.H{
			"bucketName": "foo",
		})
		c.Status(http.StatusOK)
	})

	fake := httptest.NewServer(engine)
	defer fake.Close()
	ctx := context.TODO()
	client, _ := NewTest(fake.URL, "miniokey", "miniosecret", "foo", false)

	object := &api.Object{
		ObjectName: "test.xml",
		Bucket: api.Bucket{
			BucketName: "foo",
		},
	}

	// run test
	url, err := client.PresignedGetObject(ctx, object)
	if err != nil {
		t.Errorf("PresignedGetObject returned err: %v", err)
	}

	// check if URL is valid
	if url == "" {
		t.Errorf("PresignedGetObject returned empty URL")
	}
}

func Test_PresignedGetObject_Failure(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)
	resp := httptest.NewRecorder()
	ctx, engine := gin.CreateTestContext(resp)

	// mock presigned get object call
	engine.GET("/foo/", func(c *gin.Context) {
		c.Header("Content-Type", "application/xml")
		c.XML(500, gin.H{
			"error": "Internal Server Error",
		})
		c.Status(http.StatusInternalServerError)
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
	url, err := client.PresignedGetObject(ctx, object)
	if err == nil {
		t.Errorf("PresignedGetObject expected error but got none")
	}

	if url != "" {
		t.Errorf("PresignedGetObject returned URL when it should have failed")
	}
}
