package minio

import (
	"context"
	"github.com/gin-gonic/gin"
	api "github.com/go-vela/server/api/types"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMinioClient_Bucket_Delete_Success(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	// mock create bucket call
	engine.PUT("/foo/", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
	})

	// mock delete bucket call
	engine.DELETE("/foo/", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	fake := httptest.NewServer(engine)
	defer fake.Close()
	ctx := context.TODO()
	b := new(api.Bucket)
	b.BucketName = "foo"
	client, _ := NewTest(fake.URL, "miniokey", "miniosecret", "foo", false)

	// create bucket
	err := client.CreateBucket(ctx, b)
	if err != nil {
		t.Errorf("CreateBucket returned err: %v", err)
	}

	// run test
	err = client.DeleteBucket(ctx, b)
	if resp.Code != http.StatusOK {
		t.Errorf("DeleteBucket returned %v, want %v", resp.Code, http.StatusOK)
	}

	// in Minio SDK, removeBucket returns status code 200 OK as error if a bucket is deleted successfully
	if err != nil && err.Error() != "200 OK" {
		t.Errorf("DeleteBucket returned err: %v", err)
	}

}

func TestMinioClient_Bucket_Delete_BucketNotFound(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// mock delete bucket call
	engine.DELETE("/foo/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "The specified bucket does not exist"})
	})

	fake := httptest.NewServer(engine)
	defer fake.Close()
	ctx := context.TODO()
	b := new(api.Bucket)
	b.BucketName = "foo"
	client, _ := NewTest(fake.URL, "miniokey", "miniosecret", "foo", false)

	// run test
	err := client.DeleteBucket(ctx, b)
	if resp.Code != http.StatusOK {
		t.Errorf("DeleteBucket returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err == nil {
		t.Errorf("DeleteBucket expected error, got nil")
	}

}
