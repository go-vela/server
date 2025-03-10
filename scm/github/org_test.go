// SPDX-License-Identifier: Apache-2.0

package github

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	api "github.com/go-vela/server/api/types"
)

func TestGithub_GetOrgIdentifiers(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/api/v3/orgs/:org", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/get_org.json")
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	u := new(api.User)
	u.SetName("foo")
	u.SetToken("bar")

	wantName := "github"
	wantID := int64(1)

	client, _ := NewTest(s.URL)

	// run test
	gotName, gotID, err := client.GetOrgIdentifiers(context.TODO(), u, "github")

	if resp.Code != http.StatusOK {
		t.Errorf("GetOrgName returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("GetOrgName returned err: %v", err)
	}

	if gotName != wantName {
		t.Errorf("GetOrgIdentifiers name is %v, want %v", gotName, wantName)
	}

	if gotID != wantID {
		t.Errorf("GetOrgIdentifiers id is %v, want %v", gotID, wantID)
	}
}

func TestGithub_GetOrgName_Personal(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/api/v3/user", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/user.json")
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	u := new(api.User)
	u.SetName("foo")
	u.SetToken("bar")

	wantName := "octocat"
	wantID := int64(1)

	client, _ := NewTest(s.URL)

	// run test
	gotName, gotID, err := client.GetOrgIdentifiers(context.TODO(), u, "octocat")

	if resp.Code != http.StatusOK {
		t.Errorf("GetOrgName returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("GetOrgName returned err: %v", err)
	}

	if gotName != wantName {
		t.Errorf("GetOrgIdentifiers name is %v, want %v", gotName, wantName)
	}

	if gotID != wantID {
		t.Errorf("GetOrgIdentifiers id is %v, want %v", gotID, wantID)
	}
}

func TestGithub_GetOrgName_Fail(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/api/v3/orgs/:org", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusNotFound)
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	u := new(api.User)
	u.SetName("foo")
	u.SetToken("bar")

	client, _ := NewTest(s.URL)

	// run test
	_, _, err := client.GetOrgIdentifiers(context.TODO(), u, "octocat")

	if err == nil {
		t.Error("GetOrgName should return error")
	}
}
