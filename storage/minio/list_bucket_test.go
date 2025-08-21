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
	ctx, engine := gin.CreateTestContext(resp)

	// mock list buckets call
	// Match root with or without query params
	engine.GET("/", func(c *gin.Context) {
		// Raw S3-compatible XML response for ListBuckets
		xmlResponse := `<?xml version="1.0" encoding="UTF-8"?>
<ListAllMyBucketsResult xmlns="http://s3.amazonaws.com/doc/2006-03-01/">
  <Owner>
    <ID>minio</ID>
    <DisplayName>minio</DisplayName>
  </Owner>
  <Buckets>
    <Bucket>
      <Name>foo</Name>
      <CreationDate>2025-03-20T19:01:40.968Z</CreationDate>
    </Bucket>
  </Buckets>
</ListAllMyBucketsResult>`
		c.Data(http.StatusOK, "application/xml", []byte(xmlResponse))
	})

	fake := httptest.NewServer(engine)
	defer fake.Close()

	client, _ := NewTest(fake.URL, "miniokey", "miniosecret", "foo", false)
	b := new(api.Bucket)
	b.BucketName = "foo"

	buckets, err := client.ListBuckets(ctx)
	if err != nil {
		t.Errorf("ListBuckets returned err: %v", err)
	}

	// Ignore for now as xmlDecoder from minio-go is does not parse correctly with sample data
	// check if buckets are correct
	expectedBuckets := []string{"foo"}
	if len(buckets) != len(expectedBuckets) {
		t.Errorf("Expected %d buckets, got %d", len(expectedBuckets), len(buckets))
	}
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
