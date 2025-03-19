package minio

import (
	"context"
	"github.com/gin-gonic/gin"
	api "github.com/go-vela/server/api/types"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMinioClient_List_Object_Success(t *testing.T) {
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

	// mock upload call
	engine.PUT("/foo/test.xml", func(c *gin.Context) {
		c.Header("Content-Type", "text/xml")
		c.Status(http.StatusOK)
		c.File("test_data/test.xml")
	})

	// mock list call
	engine.GET("/foo?", func(c *gin.Context) {
		listType := c.Query("list-type")
		if listType != "2" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "list-type must be 2",
			})
			return
		}

		c.XML(http.StatusOK, gin.H{
			"Contents": []gin.H{
				{"Key": "test.xml"},
			},
		})
	})
	fake := httptest.NewServer(engine)
	defer fake.Close()
	ctx := context.TODO()

	b := new(api.Bucket)
	b.BucketName = "foo"

	obj := new(api.Object)
	obj.Bucket.BucketName = "foo"
	obj.ObjectName = "test.xml"
	obj.FilePath = "test_data/test.xml"
	client, _ := NewTest(fake.URL, "miniokey", "miniosecret", "foo", false)

	// create bucket
	err := client.CreateBucket(ctx, b)
	if err != nil {
		t.Errorf("CreateBucket returned err: %v", err)
	}

	// upload test file
	err = client.Upload(ctx, obj)
	if resp.Code != http.StatusOK {
		t.Errorf("Upload returned %v, want %v", resp.Code, http.StatusOK)
	}
	// run test
	results, err := client.ListObjects(ctx, b)
	if err != nil {
		t.Errorf("ListObject returned err: %v", err)
	}

	// check if file exists in the list
	found := false
	for _, result := range results {
		if result == obj.ObjectName {
			found = true
		}
	}
	if !found {
		t.Errorf("Object %v not found in list %v", obj.ObjectName, results)
	}
}
