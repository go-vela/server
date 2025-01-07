//package minio
//
//import (
//	"context"
//	"net/http"
//	"net/http/httptest"
//	"testing"
//
//	api "github.com/go-vela/server/api/types"
//)
//

package minio

import (
	"context"
	"github.com/gin-gonic/gin"
	api "github.com/go-vela/server/api/types"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestMinioClient_BucketExists(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.HEAD("/foo/", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	fake := httptest.NewServer(engine)
	defer fake.Close()
	ctx := context.TODO()
	b := new(api.Bucket)
	b.BucketName = "foo"

	client, _ := NewTest(fake.URL, "miniokey", "miniosecret", false)

	// run test
	exists, err := client.BucketExists(ctx, b)
	if resp.Code != http.StatusOK {
		t.Errorf("BucketExists returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("BucketExists returned err: %v", err)
	}

	if !exists {
		t.Errorf("BucketExists returned %v, want %v", exists, true)
	}
}

func TestMinioClient_BucketExists_Success(t *testing.T) {
	// setup context
	//gin.SetMode(gin.TestMode)
	//
	//resp := httptest.NewRecorder()
	//_, engine := gin.CreateTestContext(resp)
	//
	//// setup mock server
	//engine.GET("/api/v3/orgs/:org", func(c *gin.Context) {
	//	c.Header("Content-Type", "application/json")
	//	c.Status(http.StatusOK)
	//	c.File("testdata/get_org.json")
	//})
	//
	//s := httptest.NewServer(engine)
	//defer s.Close()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Last-Modified", time.DateTime)
		w.Header().Set("Content-Length", "5")

		// Write less bytes than the content length.
		w.Write([]byte("12345"))
	}))
	defer srv.Close()

	//// New - instantiate minio client with options
	//clnt, err := New(srv.Listener.Addr().String(), &Options{
	//	Region: "us-east-1",
	//})
	//if err != nil {
	//	t.Fatal(err)
	//}
	// setup types
	u := new(api.Bucket)
	u.BucketName = "foo"

	//want := "minio"

	client, err := New(srv.URL, WithAccessKey("accessKey"), WithSecretKey("secretKey"), WithSecure(false))
	if err != nil {
		t.Fatal(err)
	}
	// run test
	got, err := client.BucketExists(context.TODO(), u)
	t.Logf("got: %v", got)
	// We expect an error when reading back.
	if got {
		t.Errorf("BucketExists returned %v, want %v", got, false)
		t.Errorf("BucketExists returned err: %v", err)
	}

	//
	//if err != nil {
	//	t.Errorf("GetOrgName returned err: %v", err)
	//}
	//
	//if !reflect.DeepEqual(got, want) {
	//	t.Errorf("GetOrgName is %v, want %v", got, want)
	//}

}
