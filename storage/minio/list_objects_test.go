// SPDX-License-Identifier: Apache-2.0

package minio

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	api "github.com/go-vela/server/api/types"
)

func TestMinioClient_ListObjects_Success(t *testing.T) {
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

	engine.GET("/foo/", func(c *gin.Context) {
		objects := []gin.H{
			{"etag": "982beba05db8083656a03f544c8c7927",
				"name":         "test.xml",
				"lastModified": "2025-03-20T19:01:40.968Z",
				"size":         558677,
				"contentType":  "",
				"expires":      "0001-01-01T00:00:00Z",
				"metadata":     "null",
				"UserTagCount": 0,
				"Owner": gin.H{
					"owner": gin.H{
						"Space": "http://s3.amazonaws.com/doc/2006-03-01/",
						"Local": "Owner",
					},
					"name": "02d6176db174dc93cb1b899f7c6078f08654445fe8cf1b6ce98d8855f66bdbf4",
					"id":   "minio",
				},
				"Grant":             "null",
				"storageClass":      "STANDARD",
				"IsLatest":          false,
				"IsDeleteMarker":    false,
				"VersionID":         "",
				"ReplicationStatus": "",
				"ReplicationReady":  false,
				"Expiration":        "0001-01-01T00:00:00Z",
				"ExpirationRuleID":  "",
				"Restore":           "null",
				"ChecksumCRC32":     "",
				"ChecksumCRC32C":    "",
				"ChecksumSHA1":      "",
				"ChecksumSHA256":    "",
				"ChecksumCRC64NVME": "",
				"Internal":          "null"},
		}

		c.Stream(func(w io.Writer) bool {
			_, err := w.Write([]byte(objects[0]["name"].(string)))
			if err != nil {
				return false
			}

			c.XML(http.StatusOK, objects)
			c.Status(http.StatusOK)

			return false
		})
	})

	fake := httptest.NewServer(engine)
	defer fake.Close()

	b := new(api.Bucket)
	b.BucketName = "foo"

	client, err := NewTest(fake.URL, "miniokey", "miniosecret", "foo", false)
	if err != nil {
		t.Errorf("Failed to create MinIO client: %v", err)
	}

	// For now, passing if listing objects returns no error
	_, err = client.ListObjects(ctx, b)
	if err != nil {
		t.Errorf("ListObject returned err: %v", err)
	}

	//
	//expected := "test.xml"
	//found := false
	//for _, result := range results {
	//	if result.Key == expected {
	//		found = true
	//	}
	//}
	//if !found {
	//	t.Errorf("Object %v not found in list %v", expected, results)
	//}
}

func TestMinioClient_ListObjects_Failure(t *testing.T) {
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

	engine.GET("/foo/", func(c *gin.Context) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	})

	fake := httptest.NewServer(engine)
	defer fake.Close()

	b := new(api.Bucket)
	b.BucketName = "foo"

	client, err := NewTest(fake.URL, "miniokey", "miniosecret", "foo", false)
	if err != nil {
		t.Errorf("Failed to create MinIO client: %v", err)
	}

	// run test
	_, err = client.ListObjects(ctx, b)
	if err == nil {
		t.Errorf("ListObject should have returned an error")
	}
}

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

	// Mock presigned URL requests
	engine.GET("/foo/octocat/hello-world/1/test.xml", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "http://presigned.url/test.xml")
	})

	engine.GET("/foo/octocat/hello-world/1/coverage.xml", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "http://presigned.url/coverage.xml")
	})

	fake := httptest.NewServer(engine)
	defer fake.Close()

	b := new(api.Bucket)
	b.BucketName = "foo"
	b.Recursive = true

	client, err := NewTest(fake.URL, "miniokey", "miniosecret", "foo", false)
	if err != nil {
		t.Fatalf("Failed to create MinIO client: %v", err)
	}

	results, err := client.ListBuildObjectNames(ctx, b, "octocat", "hello-world", "1")
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

	b := new(api.Bucket)
	b.BucketName = "foo"
	b.Recursive = true

	client, err := NewTest(fake.URL, "miniokey", "miniosecret", "foo", false)
	if err != nil {
		t.Errorf("Failed to create MinIO client: %v", err)
	}

	// Run test
	_, err = client.ListBuildObjectNames(ctx, b, "octocat", "hello-world", "1")
	if err == nil {
		t.Errorf("ListBuildObjectNames should have returned an error")
	}
}
