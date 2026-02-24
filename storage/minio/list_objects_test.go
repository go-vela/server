// SPDX-License-Identifier: Apache-2.0

package minio

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMinioClient_ListBuildObjectNames_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	ctx, engine := gin.CreateTestContext(resp)

	// mock create bucket call
	engine.PUT("/foo/", func(c *gin.Context) {
		c.XML(http.StatusOK, gin.H{
			"bucketName":     "foo",
			"bucketLocation": "snowball",
			"objectName":     "test.xml",
		})
	})

	// Mock bucket location check
	engine.GET("/foo/", func(c *gin.Context) {
		if _, ok := c.GetQuery("location"); ok {
			c.Data(http.StatusOK, "application/xml", []byte(`<LocationConstraint>us-east-1</LocationConstraint>`))
			return
		}

		// Handle list objects request
		prefix := c.Query("prefix")
		t.Logf("ListObjects called with prefix: %s", prefix)

		xmlResponse := `<?xml version="1.0" encoding="UTF-8"?>
<ListBucketResult xmlns="http://s3.amazonaws.com/doc/2006-03-01/">
  <Name>foo</Name>
  <Prefix>` + prefix + `</Prefix>
  <KeyCount>2</KeyCount>
  <MaxKeys>1000</MaxKeys>
  <IsTruncated>false</IsTruncated>
  <Contents>
    <Key>octocat/hello-world/1/test.xml</Key>
    <LastModified>2025-03-20T19:01:40.968Z</LastModified>
    <Size>558677</Size>
  </Contents><Contents>
    <Key>octocat/hello-world/1/coverage.xml</Key>
    <LastModified>2025-03-20T19:02:40.968Z</LastModified>
    <Size>123456</Size>
  </Contents>
</ListBucketResult>`

		c.Data(http.StatusOK, "application/xml", []byte(xmlResponse))
	})

	// mock stat object call
	engine.HEAD("/foo/octocat/hello-world/1/test.xml", func(c *gin.Context) {
		c.Header("Content-Type", "application/xml")
		c.Header("Last-Modified", "Mon, 2 Jan 2006 15:04:05 GMT")
		c.XML(200, gin.H{
			"name": "test.xml",
		})
	})
	// Mock presigned URL requests
	engine.GET("/foo/octocat/hello-world/1/test.xml", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "http://presigned.url/test.xml")
	})

	// mock stat object call
	engine.HEAD("/foo/octocat/hello-world/1/coverage.xml", func(c *gin.Context) {
		c.Header("Content-Type", "application/xml")
		c.Header("Last-Modified", "Mon, 2 Jan 2006 15:04:05 GMT")
		c.XML(200, gin.H{
			"name": "test.xml",
		})
	})
	engine.GET("/foo/octocat/hello-world/1/coverage.xml", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "http://presigned.url/coverage.xml")
	})

	fake := httptest.NewServer(engine)
	defer fake.Close()

	client, err := NewTest(fake.URL, "miniokey", "miniosecret", "foo", false)
	if err != nil {
		t.Fatalf("Failed to create MinIO client: %v", err)
	}

	results, err := client.ListBuildObjectNames(ctx, "octocat", "hello-world", "1")
	if err != nil {
		t.Fatalf("ListBuildObjectNames returned err: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}

	expectedNames := []string{
		"octocat/hello-world/1/test.xml",
		"octocat/hello-world/1/coverage.xml",
	}

	for _, expected := range expectedNames {
		if _, found := results[expected]; !found {
			t.Errorf("Expected object name %q not found in results", expected)
		}
	}
}

func TestMinioClient_ListBuildObjectNames_Failure(t *testing.T) {
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	ctx, engine := gin.CreateTestContext(resp)

	// mock bucket endpoint
	engine.PUT("/foo/", func(c *gin.Context) {
		c.XML(http.StatusOK, gin.H{
			"bucketName":     "foo",
			"bucketLocation": "snowball",
			"objectName":     "test.xml",
		})
	})

	// Return error for GET request
	engine.GET("/foo/", func(c *gin.Context) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	})

	fake := httptest.NewServer(engine)
	defer fake.Close()

	client, err := NewTest(fake.URL, "miniokey", "miniosecret", "foo", false)
	if err != nil {
		t.Errorf("Failed to create MinIO client: %v", err)
	}

	// Run test
	_, err = client.ListBuildObjectNames(ctx, "octocat", "hello-world", "1")
	if err == nil {
		t.Errorf("ListBuildObjectNames should have returned an error")
	}
}
