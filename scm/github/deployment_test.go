// SPDX-License-Identifier: Apache-2.0

package github

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/go-vela/types/library"
)

func TestGithub_CreateDeployment(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.POST("/api/v3/repos/:org/:repo/deployments", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/deployment.json")
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	u := new(library.User)
	u.SetName("foo")
	u.SetToken("bar")

	r := new(library.Repo)
	r.SetID(1)
	r.SetOrg("foo")
	r.SetName("bar")
	r.SetFullName("foo/bar")

	d := new(library.Deployment)
	d.SetID(1)
	d.SetRepoID(1)
	d.SetURL("https://api.github.com/repos/foo/bar/deployments/1")
	d.SetCommit("a84d88e7554fc1fa21bcbc4efae3c782a70d2b9d")
	d.SetRef("topic-branch")
	d.SetTask("deploy")
	d.SetTarget("production")
	d.SetDescription("Deploy request from Vela")

	client, _ := NewTest(s.URL, "https://foo.bar.com")

	// run test
	err := client.CreateDeployment(context.TODO(), u, r, d)

	if resp.Code != http.StatusOK {
		t.Errorf("CreateDeployment returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("CreateDeployment returned err: %v", err)
	}
}
