// SPDX-License-Identifier: Apache-2.0

package github

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/go-vela/types/library"
)

func TestGithub_Changeset(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/api/v3/repos/:org/:repo/commits/:ref", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/listchanges.json")
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	want := []string{"file1.txt"}

	u := new(library.User)
	u.SetName("foo")
	u.SetToken("bar")

	r := new(library.Repo)
	r.SetOrg("repos")
	r.SetName("octocat")

	client, _ := NewTest(s.URL)

	// run test
	got, err := client.Changeset(context.TODO(), u, r, "6dcb09b5b57875f334f61aebed695e2e4193db5e")

	if resp.Code != http.StatusOK {
		t.Errorf("Changeset returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("Changeset returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Changeset is %v, want %v", got, want)
	}
}

func TestGithub_ChangesetPR(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/api/v3/repos/:org/:repo/pulls/:pull_number/files", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/listchangespr.json")
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	want := []string{"file1.txt"}

	u := new(library.User)
	u.SetName("foo")
	u.SetToken("bar")

	r := new(library.Repo)
	r.SetOrg("repos")
	r.SetName("octocat")

	client, _ := NewTest(s.URL)

	// run test
	got, err := client.ChangesetPR(context.TODO(), u, r, 1)

	if resp.Code != http.StatusOK {
		t.Errorf("ChangesetPR returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("ChangesetPR returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ChangesetPR is %v, want %v", got, want)
	}
}
