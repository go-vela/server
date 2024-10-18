// SPDX-License-Identifier: Apache-2.0

package pipeline

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/compiler"
	"github.com/go-vela/server/compiler/native"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/database/testutils"
	"github.com/go-vela/server/internal"
	"github.com/go-vela/server/internal/token"
	"github.com/go-vela/server/router/middleware/claims"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/scm/github"
)

func TestPipeline_Retrieve(t *testing.T) {
	// setup types
	_pipeline := new(api.Pipeline)

	gin.SetMode(gin.TestMode)
	_context, _ := gin.CreateTestContext(nil)

	// setup tests
	tests := []struct {
		name    string
		context *gin.Context
		want    *api.Pipeline
	}{
		{
			name:    "context",
			context: _context,
			want:    _pipeline,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ToContext(test.context, test.want)

			got := Retrieve(test.context)

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("Retrieve for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}

func TestPipeline_Establish(t *testing.T) {
	// setup types
	owner := testutils.APIUser().Crop()
	owner.SetID(1)
	owner.SetName("octocat")
	owner.SetToken("foo")

	r := testutils.APIRepo()
	r.SetID(1)
	r.SetOwner(owner)
	r.SetHash("baz")
	r.SetOrg("foo")
	r.SetName("bar")
	r.SetFullName("foo/bar")
	r.SetVisibility("public")

	want := new(api.Pipeline)
	want.SetID(1)
	want.SetRepo(r)
	want.SetCommit("48afb5bdc41ad69bf22588491333f7cf71135163")
	want.SetFlavor("")
	want.SetPlatform("")
	want.SetRef("refs/heads/main")
	want.SetType("yaml")
	want.SetVersion("1")
	want.SetExternalSecrets(false)
	want.SetInternalSecrets(false)
	want.SetServices(false)
	want.SetStages(false)
	want.SetSteps(false)
	want.SetTemplates(false)
	want.SetData([]byte{})

	got := new(api.Pipeline)

	// setup database
	db, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create test database engine: %v", err)
	}

	defer func() {
		_ = db.DeleteUser(context.TODO(), owner)
		_ = db.DeletePipeline(context.TODO(), want)
		_ = db.DeleteRepo(context.TODO(), r)
		db.Close()
	}()

	_, _ = db.CreateUser(context.TODO(), owner)
	_, _ = db.CreateRepo(context.TODO(), r)
	_, _ = db.CreatePipeline(context.TODO(), want)

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/pipelines/foo/bar/48afb5bdc41ad69bf22588491333f7cf71135163", nil)

	// setup mock server
	engine.Use(func(c *gin.Context) { c.Set("logger", logrus.NewEntry(logrus.StandardLogger())) })
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(org.Establish())
	engine.Use(repo.Establish())
	engine.Use(Establish())
	engine.GET("/pipelines/:org/:repo/:pipeline", func(c *gin.Context) {
		got = Retrieve(c)

		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("Establish returned %v, want %v", resp.Code, http.StatusOK)
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Establish mismatch (-got +want):\n%v", diff)
	}
}

func TestPipeline_Establish_NoRepo(t *testing.T) {
	// setup database
	db, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create test database engine: %v", err)
	}
	defer db.Close()

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/pipelines/foo/bar/48afb5bdc41ad69bf22588491333f7cf71135163", nil)

	// setup mock server
	engine.Use(func(c *gin.Context) { c.Set("logger", logrus.NewEntry(logrus.StandardLogger())) })
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(Establish())

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusNotFound {
		t.Errorf("Establish returned %v, want %v", resp.Code, http.StatusNotFound)
	}
}

func TestPipeline_Establish_NoPipelineParameter(t *testing.T) {
	// setup types
	owner := new(api.User)
	owner.SetID(1)

	r := new(api.Repo)
	r.SetID(1)
	r.SetOwner(owner)
	r.SetHash("baz")
	r.SetOrg("foo")
	r.SetName("bar")
	r.SetFullName("foo/bar")
	r.SetVisibility("public")

	// setup database
	db, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create test database engine: %v", err)
	}

	defer func() {
		_ = db.DeleteRepo(context.TODO(), r)
		db.Close()
	}()

	_, _ = db.CreateRepo(context.TODO(), r)

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/pipelines/foo/bar", nil)

	// setup mock server
	engine.Use(func(c *gin.Context) { c.Set("logger", logrus.NewEntry(logrus.StandardLogger())) })
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(org.Establish())
	engine.Use(repo.Establish())
	engine.Use(Establish())
	engine.GET("/pipelines/:org/:repo", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("Establish returned %v, want %v", resp.Code, http.StatusBadRequest)
	}
}

func TestPipeline_Establish_NoPipeline(t *testing.T) {
	// setup types
	secret := "superSecret"

	tm := &token.Manager{
		PrivateKeyHMAC:           "123abc",
		UserAccessTokenDuration:  time.Minute * 5,
		UserRefreshTokenDuration: time.Minute * 30,
	}

	owner := new(api.User)
	owner.SetID(1)

	r := new(api.Repo)
	r.SetID(1)
	r.SetOwner(owner)
	r.SetHash("baz")
	r.SetOrg("foo")
	r.SetName("bar")
	r.SetFullName("foo/bar")
	r.SetVisibility("public")

	u := new(api.User)
	u.SetID(1)
	u.SetName("foo")
	u.SetToken("bar")
	u.SetAdmin(true)

	m := &internal.Metadata{
		Database: &internal.Database{
			Driver: "foo",
			Host:   "foo",
		},
		Queue: &internal.Queue{
			Channel: "foo",
			Driver:  "foo",
			Host:    "foo",
		},
		Source: &internal.Source{
			Driver: "foo",
			Host:   "foo",
		},
		Vela: &internal.Vela{
			Address:    "foo",
			WebAddress: "foo",
		},
	}

	mto := &token.MintTokenOpts{
		User:          u,
		TokenDuration: tm.UserAccessTokenDuration,
		TokenType:     constants.UserAccessTokenType,
	}

	at, err := tm.MintToken(mto)
	if err != nil {
		t.Errorf("unable to mint user access token: %v", err)
	}

	set := flag.NewFlagSet("test", 0)
	set.String("clone-image", "target/vela-git-slim:latest", "doc")

	comp, err := native.FromCLIContext(cli.NewContext(nil, set, nil))
	if err != nil {
		t.Errorf("unable to create compiler: %v", err)
	}

	// setup database
	db, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create test database engine: %v", err)
	}

	defer func() {
		_ = db.DeleteRepo(context.TODO(), r)
		_ = db.DeleteUser(context.TODO(), u)
		db.Close()
	}()

	_, _ = db.CreateRepo(context.TODO(), r)
	_, _ = db.CreateUser(context.TODO(), u)

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/pipelines/foo/bar/148afb5bdc41ad69bf22588491333f7cf71135163", nil)
	context.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", at))

	// setup github mock server
	engine.GET("/api/v3/repos/:org/:repo/contents/:path", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/yml.json")
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup client
	client, _ := github.NewTest(s.URL)

	// setup vela mock server
	engine.Use(func(c *gin.Context) { c.Set("logger", logrus.NewEntry(logrus.StandardLogger())) })
	engine.Use(func(c *gin.Context) { c.Set("metadata", m) })
	engine.Use(func(c *gin.Context) { c.Set("token-manager", tm) })
	engine.Use(func(c *gin.Context) { c.Set("secret", secret) })
	engine.Use(func(c *gin.Context) { compiler.WithGinContext(c, comp) })
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(func(c *gin.Context) { scm.ToContext(c, client) })
	engine.Use(claims.Establish())
	engine.Use(org.Establish())
	engine.Use(repo.Establish())
	engine.Use(user.Establish())
	engine.Use(Establish())
	engine.GET("/pipelines/:org/:repo/:pipeline", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("Establish returned %v, want %v", resp.Code, http.StatusOK)
	}
}
