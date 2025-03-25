// SPDX-License-Identifier: Apache-2.0

package minio

import (
	"context"
	"github.com/gin-gonic/gin"
	api "github.com/go-vela/server/api/types"
	"net/http"
	"net/http/httptest"
	"testing"
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
