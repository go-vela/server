// SPDX-License-Identifier: Apache-2.0

package minio

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	api "github.com/go-vela/server/api/types"
)

func TestMinioClient_CreateBucket(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	ctx, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.PUT("/foo/", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
	})

	fake := httptest.NewServer(engine)
	defer fake.Close()

	b := new(api.Bucket)
	b.BucketName = "foo"

	client, _ := NewTest(fake.URL, "miniokey", "miniosecret", "foo", false)

	// run test
	err := client.CreateBucket(ctx, b)
	if resp.Code != http.StatusOK {
		t.Errorf("CreateBucket returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("CreateBucket returned err: %v", err)
	}
}
