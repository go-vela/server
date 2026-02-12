// SPDX-License-Identifier: Apache-2.0

package minio

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/assert"

	api "github.com/go-vela/server/api/types"
)

func Test_StatObject_Success(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	ctx, engine := gin.CreateTestContext(resp)

	// mock create bucket call
	engine.GET("/foo/", func(c *gin.Context) {
		c.Header("Content-Type", "application/xml")
		c.XML(200, gin.H{
			"Buckets": []minio.BucketInfo{
				{
					Name: "foo",
				},
			},
		})
	})
	// mock stat object call
	engine.HEAD("/foo/test.xml", func(c *gin.Context) {
		c.Header("Content-Type", "application/xml")
		c.Header("Last-Modified", "Mon, 2 Jan 2006 15:04:05 GMT")
		c.XML(200, gin.H{
			"etag":         "982beba05db8083656a03f544c8c7927",
			"name":         "test.xml",
			"lastModified": "2025-03-20T19:01:40.968Z",
			"size":         558677,
			"contentType":  "",
			"expires":      time.Now(),
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
			"Expiration":        time.Now(),
			"ExpirationRuleID":  "",
			"Restore":           "null",
			"ChecksumCRC32":     "",
			"ChecksumCRC32C":    "",
			"ChecksumSHA1":      "",
			"ChecksumSHA256":    "",
			"ChecksumCRC64NVME": "",
			"Internal":          "null",
		})
		c.Status(http.StatusOK)
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
	result, err := client.StatObject(ctx, object)
	assert.NoError(t, err)
	assert.Equal(t, "test.xml", result.ObjectName)
}

func Test_StatObject_Failure(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	ctx, engine := gin.CreateTestContext(resp)

	// mock stat object call
	engine.HEAD("/foo/test.xml", func(c *gin.Context) {
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
	result, err := client.StatObject(ctx, object)
	assert.Error(t, err)
	assert.Nil(t, result)
}
