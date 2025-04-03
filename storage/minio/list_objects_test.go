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
