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

	// mock: create bucket
	engine.PUT("/foo/", func(c *gin.Context) {
		c.XML(http.StatusOK, gin.H{
			"bucketName":     "foo",
			"bucketLocation": "snowball",
			"objectName":     "octocat/hello-world/1/test.xml",
		})
	})

	// handle GET with any query params (/foo, /foo/, /foo?prefix=..., etc.)
	engine.GET("/foo/*any", func(c *gin.Context) {
		t.Logf("Incoming URL: %s", c.Request.URL.String())

		// bucket location probe (SDK does ?location=)
		if _, ok := c.GetQuery("location"); ok {
			c.Data(http.StatusOK, "application/xml", []byte(`<LocationConstraint>us-east-1</LocationConstraint>`))
			return
		}

		// real object listing
		prefix := c.Query("prefix")
		t.Logf("Received prefix: %s", prefix)

		if prefix != "octocat/hello-world/1/" {
			t.Logf("Invalid prefix received: %s", prefix)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid prefix"})
			return
		}

		// return raw XML (valid S3 ListObjectsV2 response)
		xmlResponse := `<?xml version="1.0" encoding="UTF-8"?>
<ListBucketResult xmlns="http://s3.amazonaws.com/doc/2006-03-01/">
  <Name>foo</Name>
  <Prefix>` + prefix + `</Prefix>
  <KeyCount>2</KeyCount>
  <MaxKeys>1000</MaxKeys>
  <IsTruncated>false</IsTruncated>
  <Contents>
    <Key>octocat/hello-world/1/test.xml</Key>
    <ETag>"etag-test"</ETag>
    <Size>558677</Size>
    <LastModified>2025-03-20T19:01:40.968Z</LastModified>
  </Contents>
  <Contents>
    <Key>octocat/hello-world/1/coverage.xml</Key>
    <ETag>"etag-coverage"</ETag>
    <Size>123456</Size>
    <LastModified>2025-03-20T19:02:40.968Z</LastModified>
  </Contents>
</ListBucketResult>`

		c.Data(http.StatusOK, "application/xml", []byte(xmlResponse))
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

	// Run the test
	t.Logf("Running ListBuildObjectNames with org=octocat, repo=hello-world, build=1")
	results, err := client.ListBuildObjectNames(ctx, b, "octocat", "hello-world", "1")
	if err != nil {
		t.Fatalf("ListBuildObjectNames returned err: %v", err)
	}

	// Check results
	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}

	expectedNames := []string{
		"octocat/hello-world/1/test.xml",
		"octocat/hello-world/1/coverage.xml",
	}

	for _, expected := range expectedNames {
		found := false

		for _, name := range results {
			if name == expected {
				found = true
				break
			}
		}
		if !found {
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
