// SPDX-License-Identifier: Apache-2.0

package minio

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
)

//func TestMinioClient_List_Object_Success(t *testing.T) {
//	// setup context
//	gin.SetMode(gin.TestMode)
//
//	resp := httptest.NewRecorder()
//	ctx, engine := gin.CreateTestContext(resp)
//
//	// mock list call
//	engine.GET("/foo?list-type=2", func(c *gin.Context) {
//		c.XML(http.StatusOK, gin.H{
//			"Name":        "foo",
//			"Prefix":      "",
//			"KeyCount":    1,
//			"MaxKeys":     1000,
//			"IsTruncated": false,
//			"Contents": gin.H{
//				"Key":          "test.xml",
//				"LastModified": "2021-07-01T00:00:00Z",
//				"ETag":         "1234567890",
//				"Size":         1234567890,
//				"StorageClass": "STANDARD",
//				"Owner": gin.H{
//					"ID":          "1234567890",
//					"DisplayName": "foo",
//				},
//			},
//		})
//	})
//	fake := httptest.NewServer(engine)
//	defer fake.Close()
//
//	b := new(api.Bucket)
//	b.BucketName = "foo"
//
//	client, err := NewTest(fake.URL, "miniokey", "miniosecret", "foo", false)
//	if err != nil {
//		t.Errorf("Failed to create MinIO client: %v", err)
//	}
//
//	// run test
//	results, err := client.ListObjects(ctx, b)
//	if err != nil {
//		t.Errorf("ListObject returned err: %v", err)
//	}
//
//	// check if file exists in the list
//	expected := "test.xml"
//	found := false
//	for _, result := range results {
//		if result == expected {
//			found = true
//		}
//	}
//	if !found {
//		t.Errorf("Object %v not found in list %v", expected, results)
//	}
//}

func TestMinioClient_List_Object_Success(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	// mock create bucket call
	engine.PUT("/foo", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
	})

	// mock list call
	engine.GET("/foo", func(c *gin.Context) {
		objects := []gin.H{
			{"Key": "test.xml"},
		}

		c.Stream(func(_ io.Writer) bool {
			for _, object := range objects {
				c.SSEvent("object", object)
			}
			return false
		})
	})

	fake := httptest.NewServer(engine)
	defer fake.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize MinIO client with the mock server URL
	client, err := NewTest(fake.URL, "miniokey", "miniosecret", "foo", false)
	if err != nil {
		t.Fatalf("Failed to create MinIO client: %v", err)
	}

	// Create bucket if it doesn't exist
	err = client.client.MakeBucket(ctx, "foo", minio.MakeBucketOptions{})
	if err != nil {
		t.Fatalf("Failed to create bucket: %v", err)
	}

	// List objects in the bucket
	opts := minio.ListObjectsOptions{
		UseV1:     true,
		Recursive: true,
	}
	results := []string{}
	for object := range client.client.ListObjects(ctx, "foo", opts) {
		if object.Err != nil {
			t.Errorf("ListObjects returned err: %v", object.Err)
		}
		results = append(results, object.Key)
	}

	// Check if the file exists in the list
	found := false
	for _, result := range results {
		if result == "test.xml" {
			found = true
		}
	}
	if !found {
		t.Errorf("Object %v not found in list %v", "test.xml", results)
	}
}
