package minio

import (
	"context"
	"github.com/gin-gonic/gin"
	api "github.com/go-vela/server/api/types"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMinioClient_CreateBucket(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.PUT("/foo/", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
	})

	fake := httptest.NewServer(engine)
	defer fake.Close()

	b := new(api.Bucket)
	b.BucketName = "foo"

	client, _ := NewTest(fake.URL, "miniokey", "miniosecret", false)

	// run test
	err := client.CreateBucket(context.TODO(), b)
	if resp.Code != http.StatusOK {
		t.Errorf("CreateBucket returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("CreateBucket returned err: %v", err)
	}
}
