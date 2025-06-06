// SPDX-License-Identifier: Apache-2.0

package perm

import (
	_context "context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/internal/token"
	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/claims"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/scm/github"
)

func TestPerm_MustPlatformAdmin(t *testing.T) {
	// setup types
	secret := "superSecret"

	tm := &token.Manager{
		PrivateKeyHMAC:           "123abc",
		UserAccessTokenDuration:  time.Minute * 5,
		UserRefreshTokenDuration: time.Minute * 30,
	}

	u := new(api.User)
	u.SetID(1)
	u.SetName("foob")
	u.SetToken("bar")
	u.SetAdmin(true)

	mto := &token.MintTokenOpts{
		User:          u,
		TokenDuration: tm.UserAccessTokenDuration,
		TokenType:     constants.UserAccessTokenType,
	}

	tok, _ := tm.MintToken(mto)

	// setup database
	db, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create test database engine: %v", err)
	}

	defer func() {
		_ = db.DeleteUser(_context.TODO(), u)
		db.Close()
	}()

	_, _ = db.CreateUser(_context.TODO(), u)

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)

	context.Request, _ = http.NewRequest(http.MethodGet, "/admin/users", nil)
	context.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tok))

	// setup github mock server
	engine.GET("/api/v3/user", func(c *gin.Context) {
		c.String(http.StatusOK, userPayload)
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup client
	client, _ := github.NewTest(s.URL)

	// setup vela mock server
	engine.Use(func(c *gin.Context) { c.Set("logger", logrus.NewEntry(logrus.StandardLogger())) })
	engine.Use(func(c *gin.Context) { c.Set("secret", secret) })
	engine.Use(func(c *gin.Context) { c.Set("token-manager", tm) })
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(func(c *gin.Context) { scm.ToContext(c, client) })
	engine.Use(claims.Establish())
	engine.Use(user.Establish())
	engine.Use(MustPlatformAdmin())
	engine.GET("/admin/users", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	s1 := httptest.NewServer(engine)
	defer s1.Close()

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("MustPlatAdmin returned %v, want %v", resp.Code, http.StatusOK)
	}
}

func TestPerm_MustPlatformAdmin_NotAdmin(t *testing.T) {
	// setup types
	secret := "superSecret"

	tm := &token.Manager{
		PrivateKeyHMAC:           "123abc",
		UserAccessTokenDuration:  time.Minute * 5,
		UserRefreshTokenDuration: time.Minute * 30,
	}

	u := new(api.User)
	u.SetID(1)
	u.SetName("foob")
	u.SetToken("bar")
	u.SetAdmin(false)

	mto := &token.MintTokenOpts{
		User:          u,
		TokenDuration: tm.UserAccessTokenDuration,
		TokenType:     constants.UserAccessTokenType,
	}

	tok, _ := tm.MintToken(mto)

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)

	// setup database
	db, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create test database engine: %v", err)
	}

	defer func() {
		_ = db.DeleteUser(_context.TODO(), u)
		db.Close()
	}()

	_, _ = db.CreateUser(_context.TODO(), u)

	context.Request, _ = http.NewRequest(http.MethodGet, "/admin/users", nil)
	context.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tok))

	// setup github mock server
	engine.GET("/api/v3/user", func(c *gin.Context) {
		c.String(http.StatusOK, userPayload)
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup client
	client, _ := github.NewTest(s.URL)

	// setup vela mock server
	engine.Use(func(c *gin.Context) { c.Set("logger", logrus.NewEntry(logrus.StandardLogger())) })
	engine.Use(func(c *gin.Context) { c.Set("secret", secret) })
	engine.Use(func(c *gin.Context) { c.Set("token-manager", tm) })
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(func(c *gin.Context) { scm.ToContext(c, client) })
	engine.Use(claims.Establish())
	engine.Use(user.Establish())
	engine.Use(MustPlatformAdmin())
	engine.GET("/admin/users", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	s1 := httptest.NewServer(engine)
	defer s1.Close()

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusUnauthorized {
		t.Errorf("MustPlatAdmin returned %v, want %v", resp.Code, http.StatusUnauthorized)
	}
}

func TestPerm_MustWorkerRegisterToken(t *testing.T) {
	// setup types
	tm := &token.Manager{
		PrivateKeyHMAC:              "123abc",
		UserAccessTokenDuration:     time.Minute * 5,
		UserRefreshTokenDuration:    time.Minute * 30,
		WorkerRegisterTokenDuration: time.Minute * 1,
		WorkerAuthTokenDuration:     time.Minute * 15,
	}

	mto := &token.MintTokenOpts{
		Hostname:      "worker",
		TokenDuration: tm.WorkerRegisterTokenDuration,
		TokenType:     constants.WorkerRegisterTokenType,
	}

	tok, _ := tm.MintToken(mto)

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)

	context.Request, _ = http.NewRequest(http.MethodGet, "/test/foo/bar", nil)
	context.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tok))

	// setup vela mock server
	engine.Use(func(c *gin.Context) { c.Set("logger", logrus.NewEntry(logrus.StandardLogger())) })
	engine.Use(func(c *gin.Context) { c.Set("token-manager", tm) })
	engine.Use(claims.Establish())
	engine.Use(user.Establish())
	engine.Use(MustWorkerRegisterToken())
	engine.GET("/test/:org/:repo", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	s1 := httptest.NewServer(engine)
	defer s1.Close()

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("MustWorkerRegisterToken returned %v, want %v", resp.Code, http.StatusOK)
	}
}

func TestPerm_MustWorkerRegisterToken_PlatAdmin(t *testing.T) {
	tm := &token.Manager{
		PrivateKeyHMAC:           "123abc",
		UserAccessTokenDuration:  time.Minute * 5,
		UserRefreshTokenDuration: time.Minute * 30,
	}

	u := new(api.User)
	u.SetID(1)
	u.SetName("vela-worker")
	u.SetToken("bar")
	u.SetAdmin(true)

	mto := &token.MintTokenOpts{
		User:          u,
		TokenDuration: tm.UserAccessTokenDuration,
		TokenType:     constants.UserAccessTokenType,
	}

	tok, _ := tm.MintToken(mto)

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)

	// setup database
	db, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create test database engine: %v", err)
	}

	defer func() {
		_ = db.DeleteUser(_context.TODO(), u)
		db.Close()
	}()

	_, _ = db.CreateUser(_context.TODO(), u)

	context.Request, _ = http.NewRequest(http.MethodGet, "/test/foo/bar", nil)
	context.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tok))

	// setup vela mock server
	engine.Use(func(c *gin.Context) { c.Set("logger", logrus.NewEntry(logrus.StandardLogger())) })
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(func(c *gin.Context) { c.Set("token-manager", tm) })
	engine.Use(claims.Establish())
	engine.Use(user.Establish())
	engine.Use(MustWorkerRegisterToken())
	engine.GET("/test/:org/:repo", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	s1 := httptest.NewServer(engine)
	defer s1.Close()

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusUnauthorized {
		t.Errorf("MustWorkerRegisterToken returned %v, want %v", resp.Code, http.StatusUnauthorized)
	}
}

func TestPerm_MustWorkerAuthToken(t *testing.T) {
	// setup types
	tm := &token.Manager{
		PrivateKeyHMAC:              "123abc",
		UserAccessTokenDuration:     time.Minute * 5,
		UserRefreshTokenDuration:    time.Minute * 30,
		WorkerRegisterTokenDuration: time.Minute * 1,
		WorkerAuthTokenDuration:     time.Minute * 15,
	}

	mto := &token.MintTokenOpts{
		Hostname:      "worker",
		TokenDuration: tm.WorkerAuthTokenDuration,
		TokenType:     constants.WorkerAuthTokenType,
	}

	tok, _ := tm.MintToken(mto)

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)

	context.Request, _ = http.NewRequest(http.MethodGet, "/test/foo/bar", nil)
	context.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tok))

	// setup vela mock server
	engine.Use(func(c *gin.Context) { c.Set("logger", logrus.NewEntry(logrus.StandardLogger())) })
	engine.Use(func(c *gin.Context) { c.Set("token-manager", tm) })
	engine.Use(claims.Establish())
	engine.Use(user.Establish())
	engine.Use(MustWorkerAuthToken())
	engine.GET("/test/:org/:repo", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	s1 := httptest.NewServer(engine)
	defer s1.Close()

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("MustWorkerAuthToken returned %v, want %v", resp.Code, http.StatusOK)
	}
}

func TestPerm_MustWorkerAuth_ServerWorkerToken(t *testing.T) {
	// setup types
	secret := "superSecret"
	tm := &token.Manager{
		PrivateKeyHMAC:              "123abc",
		UserAccessTokenDuration:     time.Minute * 5,
		UserRefreshTokenDuration:    time.Minute * 30,
		WorkerRegisterTokenDuration: time.Minute * 1,
		WorkerAuthTokenDuration:     time.Minute * 15,
	}

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)

	context.Request, _ = http.NewRequest(http.MethodGet, "/test/foo/bar", nil)
	context.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", secret))

	// setup vela mock server
	engine.Use(func(c *gin.Context) { c.Set("logger", logrus.NewEntry(logrus.StandardLogger())) })
	engine.Use(func(c *gin.Context) { c.Set("secret", secret) })
	engine.Use(func(c *gin.Context) { c.Set("token-manager", tm) })
	engine.Use(claims.Establish())
	engine.Use(user.Establish())
	engine.Use(MustWorkerAuthToken())
	engine.GET("/test/:org/:repo", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	s1 := httptest.NewServer(engine)
	defer s1.Close()

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("MustWorkerAuthToken returned %v, want %v", resp.Code, http.StatusOK)
	}
}

func TestPerm_MustBuildAccess(t *testing.T) {
	// setup types
	secret := "superSecret"

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

	b := new(api.Build)
	b.SetID(1)
	b.SetRepo(r)
	b.SetNumber(1)

	tm := &token.Manager{
		PrivateKeyHMAC:           "123abc",
		UserAccessTokenDuration:  time.Minute * 5,
		UserRefreshTokenDuration: time.Minute * 30,
	}

	mto := &token.MintTokenOpts{
		Hostname:      "worker",
		Build:         b,
		Repo:          "foo/bar",
		TokenDuration: time.Minute * 30,
		TokenType:     constants.WorkerBuildTokenType,
	}

	tok, _ := tm.MintToken(mto)

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)

	// setup database
	db, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create test database engine: %v", err)
	}

	ctx := _context.TODO()

	defer func() {
		_ = db.DeleteBuild(ctx, b)
		_ = db.DeleteRepo(ctx, r)
		db.Close()
	}()

	_, _ = db.CreateRepo(ctx, r)
	_, _ = db.CreateBuild(ctx, b)

	context.Request, _ = http.NewRequest(http.MethodGet, "/test/foo/bar/builds/1", nil)
	context.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tok))

	// setup vela mock server
	engine.Use(func(c *gin.Context) { c.Set("logger", logrus.NewEntry(logrus.StandardLogger())) })
	engine.Use(func(c *gin.Context) { c.Set("secret", secret) })
	engine.Use(func(c *gin.Context) { c.Set("token-manager", tm) })
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(claims.Establish())
	engine.Use(user.Establish())
	engine.Use(org.Establish())
	engine.Use(repo.Establish())
	engine.Use(build.Establish())
	engine.Use(MustBuildAccess())
	engine.GET("/test/:org/:repo/builds/:build", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	s1 := httptest.NewServer(engine)
	defer s1.Close()

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("MustBuildAccess returned %v, want %v", resp.Code, http.StatusOK)
	}
}

func TestPerm_MustBuildAccess_PlatAdmin(t *testing.T) {
	// setup types
	secret := "superSecret"

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

	b := new(api.Build)
	b.SetID(1)
	b.SetRepo(r)
	b.SetNumber(1)

	u := new(api.User)
	u.SetID(1)
	u.SetName("admin")
	u.SetToken("bar")
	u.SetAdmin(true)

	tm := &token.Manager{
		PrivateKeyHMAC:           "123abc",
		UserAccessTokenDuration:  time.Minute * 5,
		UserRefreshTokenDuration: time.Minute * 30,
	}

	mto := &token.MintTokenOpts{
		User:          u,
		TokenDuration: tm.UserAccessTokenDuration,
		TokenType:     constants.UserAccessTokenType,
	}

	tok, _ := tm.MintToken(mto)

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)

	ctx := _context.TODO()

	// setup database
	db, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create test database engine: %v", err)
	}

	defer func() {
		_ = db.DeleteBuild(ctx, b)
		_ = db.DeleteRepo(ctx, r)
		_ = db.DeleteUser(ctx, u)
		db.Close()
	}()

	_, _ = db.CreateRepo(ctx, r)
	_, _ = db.CreateBuild(ctx, b)
	_, _ = db.CreateUser(ctx, u)

	context.Request, _ = http.NewRequest(http.MethodGet, "/test/foo/bar/builds/1", nil)
	context.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tok))

	// setup vela mock server
	engine.Use(func(c *gin.Context) { c.Set("logger", logrus.NewEntry(logrus.StandardLogger())) })
	engine.Use(func(c *gin.Context) { c.Set("secret", secret) })
	engine.Use(func(c *gin.Context) { c.Set("token-manager", tm) })
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(claims.Establish())
	engine.Use(user.Establish())
	engine.Use(org.Establish())
	engine.Use(repo.Establish())
	engine.Use(build.Establish())
	engine.Use(MustBuildAccess())
	engine.GET("/test/:org/:repo/builds/:build", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	s1 := httptest.NewServer(engine)
	defer s1.Close()

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("MustBuildAccess returned %v, want %v", resp.Code, http.StatusOK)
	}
}

func TestPerm_MustBuildToken_WrongBuild(t *testing.T) {
	// setup types
	secret := "superSecret"

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

	b := new(api.Build)
	b.SetID(1)
	b.SetRepo(r)
	b.SetNumber(1)

	wB := new(api.Build)
	wB.SetID(2)

	tm := &token.Manager{
		PrivateKeyHMAC:           "123abc",
		UserAccessTokenDuration:  time.Minute * 5,
		UserRefreshTokenDuration: time.Minute * 30,
	}

	mto := &token.MintTokenOpts{
		Hostname:      "worker",
		Build:         wB,
		Repo:          "foo/bar",
		TokenDuration: time.Minute * 30,
		TokenType:     constants.WorkerBuildTokenType,
	}

	tok, _ := tm.MintToken(mto)

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)

	ctx := _context.TODO()

	// setup database
	db, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create test database engine: %v", err)
	}

	defer func() {
		_ = db.DeleteBuild(ctx, b)
		_ = db.DeleteRepo(ctx, r)
		db.Close()
	}()

	_, _ = db.CreateRepo(ctx, r)
	_, _ = db.CreateBuild(ctx, b)

	context.Request, _ = http.NewRequest(http.MethodGet, "/test/foo/bar/builds/1", nil)
	context.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tok))

	// setup vela mock server
	engine.Use(func(c *gin.Context) { c.Set("logger", logrus.NewEntry(logrus.StandardLogger())) })
	engine.Use(func(c *gin.Context) { c.Set("secret", secret) })
	engine.Use(func(c *gin.Context) { c.Set("token-manager", tm) })
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(claims.Establish())
	engine.Use(user.Establish())
	engine.Use(org.Establish())
	engine.Use(repo.Establish())
	engine.Use(build.Establish())
	engine.Use(MustBuildAccess())
	engine.GET("/test/:org/:repo/builds/:build", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	s1 := httptest.NewServer(engine)
	defer s1.Close()

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusUnauthorized {
		t.Errorf("MustBuildAccess returned %v, want %v", resp.Code, http.StatusOK)
	}
}

func TestPerm_MustIDRequestToken(t *testing.T) {
	// setup types
	secret := "superSecret"

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

	b := new(api.Build)
	b.SetID(1)
	b.SetRepo(r)
	b.SetNumber(1)
	b.SetStatus(constants.StatusRunning)
	b.SetCommit("456def")

	tm := &token.Manager{
		PrivateKeyHMAC:           "123abc",
		UserAccessTokenDuration:  time.Minute * 5,
		UserRefreshTokenDuration: time.Minute * 30,
	}

	mto := &token.MintTokenOpts{
		Hostname:      "foo/bar/456def",
		Build:         b,
		Repo:          r.GetFullName(),
		TokenDuration: time.Minute * 30,
		TokenType:     constants.IDRequestTokenType,
	}

	tok, err := tm.MintToken(mto)
	if err != nil {
		t.Errorf("unable to mint token: %v", err)
	}

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)

	// setup database
	db, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create test database engine: %v", err)
	}

	ctx := _context.TODO()

	defer func() {
		_ = db.DeleteBuild(ctx, b)
		_ = db.DeleteRepo(ctx, r)
		db.Close()
	}()

	_, _ = db.CreateRepo(ctx, r)
	_, _ = db.CreateBuild(ctx, b)

	context.Request, _ = http.NewRequest(http.MethodGet, "/test/foo/bar/builds/1", nil)
	context.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tok))

	// setup vela mock server
	engine.Use(func(c *gin.Context) { c.Set("logger", logrus.NewEntry(logrus.StandardLogger())) })
	engine.Use(func(c *gin.Context) { c.Set("secret", secret) })
	engine.Use(func(c *gin.Context) { c.Set("token-manager", tm) })
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(claims.Establish())
	engine.Use(user.Establish())
	engine.Use(org.Establish())
	engine.Use(repo.Establish())
	engine.Use(build.Establish())
	engine.Use(MustIDRequestToken())
	engine.GET("/test/:org/:repo/builds/:build", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	s1 := httptest.NewServer(engine)
	defer s1.Close()

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("MustIDRequestToken returned %v, want %v: %v", resp.Code, http.StatusOK, resp.Body.String())
	}
}

func TestPerm_MustIDRequestToken_NotRunning(t *testing.T) {
	// setup types
	secret := "superSecret"

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

	b := new(api.Build)
	b.SetID(1)
	b.SetRepo(r)
	b.SetNumber(1)
	b.SetStatus(constants.StatusSuccess)
	b.SetCommit("456def")

	u := new(api.User)
	u.SetID(1)
	u.SetName("admin")
	u.SetToken("bar")
	u.SetAdmin(true)

	tm := &token.Manager{
		PrivateKeyHMAC:           "123abc",
		UserAccessTokenDuration:  time.Minute * 5,
		UserRefreshTokenDuration: time.Minute * 30,
	}

	mto := &token.MintTokenOpts{
		Hostname:      "foo/bar/456def",
		Build:         b,
		Repo:          "foo/bar",
		TokenDuration: time.Minute * 30,
		TokenType:     constants.IDRequestTokenType,
	}

	tok, _ := tm.MintToken(mto)

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)

	ctx := _context.TODO()

	// setup database
	db, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create test database engine: %v", err)
	}

	defer func() {
		_ = db.DeleteBuild(ctx, b)
		_ = db.DeleteRepo(ctx, r)
		_ = db.DeleteUser(ctx, u)
		db.Close()
	}()

	_, _ = db.CreateRepo(ctx, r)
	_, _ = db.CreateBuild(ctx, b)
	_, _ = db.CreateUser(ctx, u)

	context.Request, _ = http.NewRequest(http.MethodGet, "/test/foo/bar/builds/1", nil)
	context.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tok))

	// setup vela mock server
	engine.Use(func(c *gin.Context) { c.Set("logger", logrus.NewEntry(logrus.StandardLogger())) })
	engine.Use(func(c *gin.Context) { c.Set("secret", secret) })
	engine.Use(func(c *gin.Context) { c.Set("token-manager", tm) })
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(claims.Establish())
	engine.Use(user.Establish())
	engine.Use(org.Establish())
	engine.Use(repo.Establish())
	engine.Use(build.Establish())
	engine.Use(MustIDRequestToken())
	engine.GET("/test/:org/:repo/builds/:build", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	s1 := httptest.NewServer(engine)
	defer s1.Close()

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("MustIDRequestToken returned %v, want %v", resp.Code, http.StatusOK)
	}
}

func TestPerm_MustIDRequestToken_WrongBuild(t *testing.T) {
	// setup types
	secret := "superSecret"

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

	b := new(api.Build)
	b.SetID(1)
	b.SetRepo(r)
	b.SetNumber(1)

	wB := new(api.Build)
	wB.SetID(2)

	tm := &token.Manager{
		PrivateKeyHMAC:           "123abc",
		UserAccessTokenDuration:  time.Minute * 5,
		UserRefreshTokenDuration: time.Minute * 30,
	}

	mto := &token.MintTokenOpts{
		Hostname:      "foo/bar/456def",
		Build:         wB,
		Repo:          "foo/bar",
		TokenDuration: time.Minute * 30,
		TokenType:     constants.IDRequestTokenType,
	}

	tok, _ := tm.MintToken(mto)

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)

	ctx := _context.TODO()

	// setup database
	db, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create test database engine: %v", err)
	}

	defer func() {
		_ = db.DeleteBuild(ctx, b)
		_ = db.DeleteRepo(ctx, r)
		db.Close()
	}()

	_, _ = db.CreateRepo(ctx, r)
	_, _ = db.CreateBuild(ctx, b)

	context.Request, _ = http.NewRequest(http.MethodGet, "/test/foo/bar/builds/1", nil)
	context.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tok))

	// setup vela mock server
	engine.Use(func(c *gin.Context) { c.Set("logger", logrus.NewEntry(logrus.StandardLogger())) })
	engine.Use(func(c *gin.Context) { c.Set("secret", secret) })
	engine.Use(func(c *gin.Context) { c.Set("token-manager", tm) })
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(claims.Establish())
	engine.Use(user.Establish())
	engine.Use(org.Establish())
	engine.Use(repo.Establish())
	engine.Use(build.Establish())
	engine.Use(MustIDRequestToken())
	engine.GET("/test/:org/:repo/builds/:build", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	s1 := httptest.NewServer(engine)
	defer s1.Close()

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("MustBuildAccess returned %v, want %v", resp.Code, http.StatusBadRequest)
	}
}

func TestPerm_MustSecretAdmin_BuildToken_Repo(t *testing.T) {
	// setup types
	secret := "superSecret"

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

	b := new(api.Build)
	b.SetID(1)
	b.SetRepo(r)
	b.SetNumber(1)

	tm := &token.Manager{
		PrivateKeyHMAC:           "123abc",
		UserAccessTokenDuration:  time.Minute * 5,
		UserRefreshTokenDuration: time.Minute * 30,
	}

	mto := &token.MintTokenOpts{
		Hostname:      "worker",
		Build:         b,
		Repo:          "foo/bar",
		TokenDuration: time.Minute * 30,
		TokenType:     constants.WorkerBuildTokenType,
	}

	tok, _ := tm.MintToken(mto)

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)

	ctx := _context.TODO()

	// setup database
	db, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create test database engine: %v", err)
	}

	defer func() {
		_ = db.DeleteBuild(ctx, b)
		_ = db.DeleteRepo(ctx, r)
		db.Close()
	}()

	_, _ = db.CreateRepo(ctx, r)
	_, _ = db.CreateBuild(ctx, b)

	context.Request, _ = http.NewRequest(http.MethodGet, "/test/native/repo/foo/bar/baz", nil)
	context.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tok))

	// setup vela mock server
	engine.Use(func(c *gin.Context) { c.Set("logger", logrus.NewEntry(logrus.StandardLogger())) })
	engine.Use(func(c *gin.Context) { c.Set("secret", secret) })
	engine.Use(func(c *gin.Context) { c.Set("token-manager", tm) })
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(claims.Establish())
	engine.Use(user.Establish())
	engine.Use(MustSecretAdmin())
	engine.GET("/test/:engine/:type/:org/:name/:secret", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	s1 := httptest.NewServer(engine)
	defer s1.Close()

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("MustBuildAccess returned %v, want %v", resp.Code, http.StatusOK)
	}
}

func TestPerm_MustSecretAdmin_BuildToken_Org(t *testing.T) {
	// setup types
	secret := "superSecret"

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

	b := new(api.Build)
	b.SetID(1)
	b.SetRepo(r)
	b.SetNumber(1)

	tm := &token.Manager{
		PrivateKeyHMAC:           "123abc",
		UserAccessTokenDuration:  time.Minute * 5,
		UserRefreshTokenDuration: time.Minute * 30,
	}

	mto := &token.MintTokenOpts{
		Hostname:      "worker",
		Build:         b,
		Repo:          "foo/bar",
		TokenDuration: time.Minute * 30,
		TokenType:     constants.WorkerBuildTokenType,
	}

	tok, _ := tm.MintToken(mto)

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)

	ctx := _context.TODO()

	// setup database
	db, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create test database engine: %v", err)
	}

	defer func() {
		_ = db.DeleteBuild(ctx, b)
		_ = db.DeleteRepo(ctx, r)
		db.Close()
	}()

	_, _ = db.CreateRepo(ctx, r)
	_, _ = db.CreateBuild(ctx, b)

	context.Request, _ = http.NewRequest(http.MethodGet, "/test/native/org/foo/*/baz", nil)
	context.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tok))

	// setup vela mock server
	engine.Use(func(c *gin.Context) { c.Set("logger", logrus.NewEntry(logrus.StandardLogger())) })
	engine.Use(func(c *gin.Context) { c.Set("secret", secret) })
	engine.Use(func(c *gin.Context) { c.Set("token-manager", tm) })
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(claims.Establish())
	engine.Use(user.Establish())
	engine.Use(MustSecretAdmin())
	engine.GET("/test/:engine/:type/:org/:name/:secret", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	s1 := httptest.NewServer(engine)
	defer s1.Close()

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("MustSecretAdmin returned %v, want %v", resp.Code, http.StatusOK)
	}
}

func TestPerm_MustSecretAdmin_BuildToken_Shared(t *testing.T) {
	// setup types
	secret := "superSecret"

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

	b := new(api.Build)
	b.SetID(1)
	b.SetRepo(r)
	b.SetNumber(1)

	tm := &token.Manager{
		PrivateKeyHMAC:           "123abc",
		UserAccessTokenDuration:  time.Minute * 5,
		UserRefreshTokenDuration: time.Minute * 30,
	}

	mto := &token.MintTokenOpts{
		Hostname:      "worker",
		Build:         b,
		Repo:          "foo/bar",
		TokenDuration: time.Minute * 30,
		TokenType:     constants.WorkerBuildTokenType,
	}

	tok, _ := tm.MintToken(mto)

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)

	ctx := _context.TODO()

	// setup database
	db, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create test database engine: %v", err)
	}

	defer func() {
		_ = db.DeleteBuild(ctx, b)
		_ = db.DeleteRepo(ctx, r)
		db.Close()
	}()

	_, _ = db.CreateRepo(ctx, r)
	_, _ = db.CreateBuild(ctx, b)

	context.Request, _ = http.NewRequest(http.MethodGet, "/test/native/shared/foo/*/*", nil)
	context.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tok))

	// setup vela mock server
	engine.Use(func(c *gin.Context) { c.Set("logger", logrus.NewEntry(logrus.StandardLogger())) })
	engine.Use(func(c *gin.Context) { c.Set("secret", secret) })
	engine.Use(func(c *gin.Context) { c.Set("token-manager", tm) })
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(claims.Establish())
	engine.Use(user.Establish())
	engine.Use(MustSecretAdmin())
	engine.GET("/test/:engine/:type/:org/:name/:secret", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	s1 := httptest.NewServer(engine)
	defer s1.Close()

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("MustSecretAdmin returned %v, want %v", resp.Code, http.StatusOK)
	}
}

func TestPerm_MustAdmin(t *testing.T) {
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
	u.SetName("foob")
	u.SetToken("bar")
	u.SetAdmin(false)

	mto := &token.MintTokenOpts{
		User:          u,
		TokenDuration: tm.UserAccessTokenDuration,
		TokenType:     constants.UserAccessTokenType,
	}

	tok, _ := tm.MintToken(mto)

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)

	// setup database
	db, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create test database engine: %v", err)
	}

	defer func() {
		_ = db.DeleteRepo(_context.TODO(), r)
		_ = db.DeleteUser(_context.TODO(), u)
		db.Close()
	}()

	_, _ = db.CreateRepo(_context.TODO(), r)
	_, _ = db.CreateUser(_context.TODO(), u)

	context.Request, _ = http.NewRequest(http.MethodGet, "/test/foo/bar", nil)
	context.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tok))

	// setup github mock server
	engine.GET("/api/v3/repos/:org/:repo/collaborators/:username/permission", func(c *gin.Context) {
		c.String(http.StatusOK, permAdminPayload)
	})
	engine.GET("/api/v3/user", func(c *gin.Context) {
		c.String(http.StatusOK, userPayload)
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup client
	client, _ := github.NewTest(s.URL)

	// setup vela mock server
	engine.Use(func(c *gin.Context) { c.Set("logger", logrus.NewEntry(logrus.StandardLogger())) })
	engine.Use(func(c *gin.Context) { c.Set("secret", secret) })
	engine.Use(func(c *gin.Context) { c.Set("token-manager", tm) })
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(func(c *gin.Context) { scm.ToContext(c, client) })
	engine.Use(claims.Establish())
	engine.Use(user.Establish())
	engine.Use(org.Establish())
	engine.Use(repo.Establish())
	engine.Use(MustAdmin())
	engine.GET("/test/:org/:repo", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	s1 := httptest.NewServer(engine)
	defer s1.Close()

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("MustAdmin returned %v, want %v", resp.Code, http.StatusOK)
	}
}

func TestPerm_MustAdmin_PlatAdmin(t *testing.T) {
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
	u.SetName("foob")
	u.SetToken("bar")
	u.SetAdmin(true)

	mto := &token.MintTokenOpts{
		User:          u,
		TokenDuration: tm.UserAccessTokenDuration,
		TokenType:     constants.UserAccessTokenType,
	}

	tok, _ := tm.MintToken(mto)

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)

	// setup database
	db, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create test database engine: %v", err)
	}

	defer func() {
		_ = db.DeleteRepo(_context.TODO(), r)
		_ = db.DeleteUser(_context.TODO(), u)
		db.Close()
	}()

	_, _ = db.CreateRepo(_context.TODO(), r)
	_, _ = db.CreateUser(_context.TODO(), u)

	context.Request, _ = http.NewRequest(http.MethodGet, "/test/foo/bar", nil)
	context.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tok))

	// setup github mock server
	engine.GET("/api/v3/repos/:org/:repo/collaborators/:username/permission", func(c *gin.Context) {
		c.String(http.StatusOK, permWritePayload)
	})
	engine.GET("/api/v3/user", func(c *gin.Context) {
		c.String(http.StatusOK, userPayload)
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup client
	client, _ := github.NewTest(s.URL)

	// setup vela mock server
	engine.Use(func(c *gin.Context) { c.Set("logger", logrus.NewEntry(logrus.StandardLogger())) })
	engine.Use(func(c *gin.Context) { c.Set("secret", secret) })
	engine.Use(func(c *gin.Context) { c.Set("token-manager", tm) })
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(func(c *gin.Context) { scm.ToContext(c, client) })
	engine.Use(claims.Establish())
	engine.Use(user.Establish())
	engine.Use(org.Establish())
	engine.Use(repo.Establish())
	engine.Use(MustAdmin())
	engine.GET("/test/:org/:repo", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	s1 := httptest.NewServer(engine)
	defer s1.Close()

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("MustAdmin returned %v, want %v", resp.Code, http.StatusOK)
	}
}

func TestPerm_MustAdmin_NotAdmin(t *testing.T) {
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
	u.SetName("foob")
	u.SetToken("bar")
	u.SetAdmin(false)

	mto := &token.MintTokenOpts{
		User:          u,
		TokenDuration: tm.UserAccessTokenDuration,
		TokenType:     constants.UserAccessTokenType,
	}

	tok, _ := tm.MintToken(mto)

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)

	// setup database
	db, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create test database engine: %v", err)
	}

	defer func() {
		_ = db.DeleteRepo(_context.TODO(), r)
		_ = db.DeleteUser(_context.TODO(), u)
		db.Close()
	}()

	_, _ = db.CreateRepo(_context.TODO(), r)
	_, _ = db.CreateUser(_context.TODO(), u)

	context.Request, _ = http.NewRequest(http.MethodGet, "/test/foo/bar", nil)
	context.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tok))

	// setup github mock server
	engine.GET("/api/v3/repos/:org/:repo/collaborators/:username/permission", func(c *gin.Context) {
		c.String(http.StatusOK, permWritePayload)
	})
	engine.GET("/api/v3/user", func(c *gin.Context) {
		c.String(http.StatusOK, userPayload)
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup client
	client, _ := github.NewTest(s.URL)

	// setup vela mock server
	engine.Use(func(c *gin.Context) { c.Set("logger", logrus.NewEntry(logrus.StandardLogger())) })
	engine.Use(func(c *gin.Context) { c.Set("secret", secret) })
	engine.Use(func(c *gin.Context) { c.Set("token-manager", tm) })
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(func(c *gin.Context) { scm.ToContext(c, client) })
	engine.Use(claims.Establish())
	engine.Use(user.Establish())
	engine.Use(org.Establish())
	engine.Use(repo.Establish())
	engine.Use(MustAdmin())
	engine.GET("/test/:org/:repo", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	s1 := httptest.NewServer(engine)
	defer s1.Close()

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusUnauthorized {
		t.Errorf("MustAdmin returned %v, want %v", resp.Code, http.StatusUnauthorized)
	}
}

func TestPerm_MustWrite(t *testing.T) {
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
	u.SetName("foob")
	u.SetToken("bar")
	u.SetAdmin(false)

	mto := &token.MintTokenOpts{
		User:          u,
		TokenDuration: tm.UserAccessTokenDuration,
		TokenType:     constants.UserAccessTokenType,
	}

	tok, _ := tm.MintToken(mto)

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)

	// setup database
	db, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create test database engine: %v", err)
	}

	defer func() {
		_ = db.DeleteRepo(_context.TODO(), r)
		_ = db.DeleteUser(_context.TODO(), u)
		db.Close()
	}()

	_, _ = db.CreateRepo(_context.TODO(), r)
	_, _ = db.CreateUser(_context.TODO(), u)

	context.Request, _ = http.NewRequest(http.MethodGet, "/test/foo/bar", nil)
	context.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tok))

	// setup github mock server
	engine.GET("/api/v3/repos/:org/:repo/collaborators/:username/permission", func(c *gin.Context) {
		c.String(http.StatusOK, permWritePayload)
	})
	engine.GET("/api/v3/user", func(c *gin.Context) {
		c.String(http.StatusOK, userPayload)
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup client
	client, _ := github.NewTest(s.URL)

	// setup vela mock server
	engine.Use(func(c *gin.Context) { c.Set("logger", logrus.NewEntry(logrus.StandardLogger())) })
	engine.Use(func(c *gin.Context) { c.Set("secret", secret) })
	engine.Use(func(c *gin.Context) { c.Set("token-manager", tm) })
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(func(c *gin.Context) { scm.ToContext(c, client) })
	engine.Use(claims.Establish())
	engine.Use(user.Establish())
	engine.Use(org.Establish())
	engine.Use(repo.Establish())
	engine.Use(MustWrite())
	engine.GET("/test/:org/:repo", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	s1 := httptest.NewServer(engine)
	defer s1.Close()

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("MustWrite returned %v, want %v", resp.Code, http.StatusOK)
	}
}

func TestPerm_MustWrite_PlatAdmin(t *testing.T) {
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
	u.SetName("foob")
	u.SetToken("bar")
	u.SetAdmin(true)

	mto := &token.MintTokenOpts{
		User:          u,
		TokenDuration: tm.UserAccessTokenDuration,
		TokenType:     constants.UserAccessTokenType,
	}

	tok, _ := tm.MintToken(mto)

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)

	// setup database
	db, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create test database engine: %v", err)
	}

	defer func() {
		_ = db.DeleteRepo(_context.TODO(), r)
		_ = db.DeleteUser(_context.TODO(), u)
		db.Close()
	}()

	_, _ = db.CreateRepo(_context.TODO(), r)
	_, _ = db.CreateUser(_context.TODO(), u)

	context.Request, _ = http.NewRequest(http.MethodGet, "/test/foo/bar", nil)
	context.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tok))

	// setup github mock server
	engine.GET("/api/v3/repos/:org/:repo/collaborators/:username/permission", func(c *gin.Context) {
		c.String(http.StatusOK, permWritePayload)
	})
	engine.GET("/api/v3/user", func(c *gin.Context) {
		c.String(http.StatusOK, userPayload)
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup client
	client, _ := github.NewTest(s.URL)

	// setup vela mock server
	engine.Use(func(c *gin.Context) { c.Set("logger", logrus.NewEntry(logrus.StandardLogger())) })
	engine.Use(func(c *gin.Context) { c.Set("secret", secret) })
	engine.Use(func(c *gin.Context) { c.Set("token-manager", tm) })
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(func(c *gin.Context) { scm.ToContext(c, client) })
	engine.Use(claims.Establish())
	engine.Use(user.Establish())
	engine.Use(org.Establish())
	engine.Use(repo.Establish())
	engine.Use(MustWrite())
	engine.GET("/test/:org/:repo", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	s1 := httptest.NewServer(engine)
	defer s1.Close()

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("MustWrite returned %v, want %v", resp.Code, http.StatusOK)
	}
}

func TestPerm_MustWrite_RepoAdmin(t *testing.T) {
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
	u.SetName("foob")
	u.SetToken("bar")
	u.SetAdmin(false)

	mto := &token.MintTokenOpts{
		User:          u,
		TokenDuration: tm.UserAccessTokenDuration,
		TokenType:     constants.UserAccessTokenType,
	}

	tok, _ := tm.MintToken(mto)

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)

	// setup database
	db, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create test database engine: %v", err)
	}

	defer func() {
		_ = db.DeleteRepo(_context.TODO(), r)
		_ = db.DeleteUser(_context.TODO(), u)
		db.Close()
	}()

	_, _ = db.CreateRepo(_context.TODO(), r)
	_, _ = db.CreateUser(_context.TODO(), u)

	context.Request, _ = http.NewRequest(http.MethodGet, "/test/foo/bar", nil)
	context.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tok))

	// setup github mock server
	engine.GET("/api/v3/repos/:org/:repo/collaborators/:username/permission", func(c *gin.Context) {
		c.String(http.StatusOK, permAdminPayload)
	})
	engine.GET("/api/v3/user", func(c *gin.Context) {
		c.String(http.StatusOK, userPayload)
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup client
	client, _ := github.NewTest(s.URL)

	// setup vela mock server
	engine.Use(func(c *gin.Context) { c.Set("logger", logrus.NewEntry(logrus.StandardLogger())) })
	engine.Use(func(c *gin.Context) { c.Set("secret", secret) })
	engine.Use(func(c *gin.Context) { c.Set("token-manager", tm) })
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(func(c *gin.Context) { scm.ToContext(c, client) })
	engine.Use(claims.Establish())
	engine.Use(user.Establish())
	engine.Use(org.Establish())
	engine.Use(repo.Establish())
	engine.Use(MustWrite())
	engine.GET("/test/:org/:repo", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	s1 := httptest.NewServer(engine)
	defer s1.Close()

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("MustWrite returned %v, want %v", resp.Code, http.StatusOK)
	}
}

func TestPerm_MustWrite_NotWrite(t *testing.T) {
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
	u.SetName("foob")
	u.SetToken("bar")
	u.SetAdmin(false)

	mto := &token.MintTokenOpts{
		User:          u,
		TokenDuration: tm.UserAccessTokenDuration,
		TokenType:     constants.UserAccessTokenType,
	}

	tok, _ := tm.MintToken(mto)

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)

	// setup database
	db, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create test database engine: %v", err)
	}

	defer func() {
		_ = db.DeleteRepo(_context.TODO(), r)
		_ = db.DeleteUser(_context.TODO(), u)
		db.Close()
	}()

	_, _ = db.CreateRepo(_context.TODO(), r)
	_, _ = db.CreateUser(_context.TODO(), u)

	context.Request, _ = http.NewRequest(http.MethodGet, "/test/foo/bar", nil)
	context.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tok))

	// setup github mock server
	engine.GET("/api/v3/repos/:org/:repo/collaborators/:username/permission", func(c *gin.Context) {
		c.String(http.StatusOK, permReadPayload)
	})
	engine.GET("/api/v3/user", func(c *gin.Context) {
		c.String(http.StatusOK, userPayload)
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup client
	client, _ := github.NewTest(s.URL)

	// setup vela mock server
	engine.Use(func(c *gin.Context) { c.Set("logger", logrus.NewEntry(logrus.StandardLogger())) })
	engine.Use(func(c *gin.Context) { c.Set("secret", secret) })
	engine.Use(func(c *gin.Context) { c.Set("token-manager", tm) })
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(func(c *gin.Context) { scm.ToContext(c, client) })
	engine.Use(claims.Establish())
	engine.Use(user.Establish())
	engine.Use(org.Establish())
	engine.Use(repo.Establish())
	engine.Use(MustWrite())
	engine.GET("/test/:org/:repo", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	s1 := httptest.NewServer(engine)
	defer s1.Close()

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusUnauthorized {
		t.Errorf("MustWrite returned %v, want %v", resp.Code, http.StatusUnauthorized)
	}
}

func TestPerm_MustRead(t *testing.T) {
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
	r.SetVisibility("private")

	u := new(api.User)
	u.SetID(1)
	u.SetName("foob")
	u.SetToken("bar")
	u.SetAdmin(false)

	mto := &token.MintTokenOpts{
		User:          u,
		TokenDuration: tm.UserAccessTokenDuration,
		TokenType:     constants.UserAccessTokenType,
	}

	tok, _ := tm.MintToken(mto)

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)

	// setup database
	db, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create test database engine: %v", err)
	}

	defer func() {
		_ = db.DeleteRepo(_context.TODO(), r)
		_ = db.DeleteUser(_context.TODO(), u)
		db.Close()
	}()

	_, _ = db.CreateRepo(_context.TODO(), r)
	_, _ = db.CreateUser(_context.TODO(), u)

	context.Request, _ = http.NewRequest(http.MethodGet, "/test/foo/bar", nil)
	context.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tok))

	// setup github mock server
	engine.GET("/api/v3/repos/:org/:repo/collaborators/:username/permission", func(c *gin.Context) {
		c.String(http.StatusOK, permReadPayload)
	})
	engine.GET("/api/v3/user", func(c *gin.Context) {
		c.String(http.StatusOK, userPayload)
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup client
	client, _ := github.NewTest(s.URL)

	// setup vela mock server
	engine.Use(func(c *gin.Context) { c.Set("logger", logrus.NewEntry(logrus.StandardLogger())) })
	engine.Use(func(c *gin.Context) { c.Set("secret", secret) })
	engine.Use(func(c *gin.Context) { c.Set("token-manager", tm) })
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(func(c *gin.Context) { scm.ToContext(c, client) })
	engine.Use(claims.Establish())
	engine.Use(user.Establish())
	engine.Use(org.Establish())
	engine.Use(repo.Establish())
	engine.Use(MustRead())
	engine.GET("/test/:org/:repo", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	s1 := httptest.NewServer(engine)
	defer s1.Close()

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("MustRead returned %v, want %v", resp.Code, http.StatusOK)
	}
}

func TestPerm_MustRead_PlatAdmin(t *testing.T) {
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
	r.SetVisibility("private")

	u := new(api.User)
	u.SetID(1)
	u.SetName("foob")
	u.SetToken("bar")
	u.SetAdmin(true)

	mto := &token.MintTokenOpts{
		User:          u,
		TokenDuration: tm.UserAccessTokenDuration,
		TokenType:     constants.UserAccessTokenType,
	}

	tok, _ := tm.MintToken(mto)

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)

	// setup database
	db, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create test database engine: %v", err)
	}

	defer func() {
		_ = db.DeleteRepo(_context.TODO(), r)
		_ = db.DeleteUser(_context.TODO(), u)
		db.Close()
	}()

	_, _ = db.CreateRepo(_context.TODO(), r)
	_, _ = db.CreateUser(_context.TODO(), u)

	context.Request, _ = http.NewRequest(http.MethodGet, "/test/foo/bar", nil)
	context.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tok))

	// setup github mock server
	engine.GET("/api/v3/repos/:org/:repo/collaborators/:username/permission", func(c *gin.Context) {
		c.String(http.StatusOK, permReadPayload)
	})
	engine.GET("/api/v3/user", func(c *gin.Context) {
		c.String(http.StatusOK, userPayload)
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup client
	client, _ := github.NewTest(s.URL)

	// setup vela mock server
	engine.Use(func(c *gin.Context) { c.Set("logger", logrus.NewEntry(logrus.StandardLogger())) })
	engine.Use(func(c *gin.Context) { c.Set("secret", secret) })
	engine.Use(func(c *gin.Context) { c.Set("token-manager", tm) })
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(func(c *gin.Context) { scm.ToContext(c, client) })
	engine.Use(claims.Establish())
	engine.Use(user.Establish())
	engine.Use(org.Establish())
	engine.Use(repo.Establish())
	engine.Use(MustRead())
	engine.GET("/test/:org/:repo", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	s1 := httptest.NewServer(engine)
	defer s1.Close()

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("MustRead returned %v, want %v", resp.Code, http.StatusOK)
	}
}

func TestPerm_MustRead_WorkerBuildToken(t *testing.T) {
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
	r.SetVisibility("private")

	b := new(api.Build)
	b.SetID(1)
	b.SetRepo(r)
	b.SetNumber(1)

	mto := &token.MintTokenOpts{
		Hostname:      "worker",
		TokenDuration: time.Minute * 35,
		TokenType:     constants.WorkerBuildTokenType,
		Build:         b,
		Repo:          "foo/bar",
	}

	tok, _ := tm.MintToken(mto)

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)

	ctx := _context.TODO()

	// setup database
	db, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create test database engine: %v", err)
	}

	defer func() {
		_ = db.DeleteBuild(ctx, b)
		_ = db.DeleteRepo(ctx, r)
		db.Close()
	}()

	_, _ = db.CreateBuild(ctx, b)
	_, _ = db.CreateRepo(ctx, r)

	context.Request, _ = http.NewRequest(http.MethodGet, "/test/foo/bar/builds/1", nil)
	context.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tok))

	// setup vela mock server
	engine.Use(func(c *gin.Context) { c.Set("logger", logrus.NewEntry(logrus.StandardLogger())) })
	engine.Use(func(c *gin.Context) { c.Set("secret", secret) })
	engine.Use(func(c *gin.Context) { c.Set("token-manager", tm) })
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(claims.Establish())
	engine.Use(user.Establish())
	engine.Use(org.Establish())
	engine.Use(repo.Establish())
	engine.Use(build.Establish())
	engine.Use(MustRead())
	engine.GET("/test/:org/:repo/builds/:build", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	s1 := httptest.NewServer(engine)
	defer s1.Close()

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("MustRead returned %v, want %v", resp.Code, http.StatusOK)
	}
}

func TestPerm_MustRead_RepoAdmin(t *testing.T) {
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
	r.SetVisibility("private")

	u := new(api.User)
	u.SetID(1)
	u.SetName("foob")
	u.SetToken("bar")
	u.SetAdmin(false)

	mto := &token.MintTokenOpts{
		User:          u,
		TokenDuration: tm.UserAccessTokenDuration,
		TokenType:     constants.UserAccessTokenType,
	}

	tok, _ := tm.MintToken(mto)

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)

	// setup database
	db, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create test database engine: %v", err)
	}

	defer func() {
		_ = db.DeleteRepo(_context.TODO(), r)
		_ = db.DeleteUser(_context.TODO(), u)
		db.Close()
	}()

	_, _ = db.CreateRepo(_context.TODO(), r)
	_, _ = db.CreateUser(_context.TODO(), u)

	context.Request, _ = http.NewRequest(http.MethodGet, "/test/foo/bar", nil)
	context.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tok))

	// setup github mock server
	engine.GET("/api/v3/repos/:org/:repo/collaborators/:username/permission", func(c *gin.Context) {
		c.String(http.StatusOK, permAdminPayload)
	})
	engine.GET("/api/v3/user", func(c *gin.Context) {
		c.String(http.StatusOK, userPayload)
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup client
	client, _ := github.NewTest(s.URL)

	// setup vela mock server
	engine.Use(func(c *gin.Context) { c.Set("logger", logrus.NewEntry(logrus.StandardLogger())) })
	engine.Use(func(c *gin.Context) { c.Set("secret", secret) })
	engine.Use(func(c *gin.Context) { c.Set("token-manager", tm) })
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(func(c *gin.Context) { scm.ToContext(c, client) })
	engine.Use(claims.Establish())
	engine.Use(user.Establish())
	engine.Use(org.Establish())
	engine.Use(repo.Establish())
	engine.Use(MustRead())
	engine.GET("/test/:org/:repo", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	s1 := httptest.NewServer(engine)
	defer s1.Close()

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("MustRead returned %v, want %v", resp.Code, http.StatusOK)
	}
}

func TestPerm_MustRead_RepoWrite(t *testing.T) {
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
	r.SetVisibility("private")

	u := new(api.User)
	u.SetID(1)
	u.SetName("foob")
	u.SetToken("bar")
	u.SetAdmin(false)

	mto := &token.MintTokenOpts{
		User:          u,
		TokenDuration: tm.UserAccessTokenDuration,
		TokenType:     constants.UserAccessTokenType,
	}

	tok, _ := tm.MintToken(mto)

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)

	// setup database
	db, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create test database engine: %v", err)
	}

	defer func() {
		_ = db.DeleteRepo(_context.TODO(), r)
		_ = db.DeleteUser(_context.TODO(), u)
		db.Close()
	}()

	_, _ = db.CreateRepo(_context.TODO(), r)
	_, _ = db.CreateUser(_context.TODO(), u)

	context.Request, _ = http.NewRequest(http.MethodGet, "/test/foo/bar", nil)
	context.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tok))

	// setup github mock server
	engine.GET("/api/v3/repos/:org/:repo/collaborators/:username/permission", func(c *gin.Context) {
		c.String(http.StatusOK, permWritePayload)
	})
	engine.GET("/api/v3/user", func(c *gin.Context) {
		c.String(http.StatusOK, userPayload)
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup client
	client, _ := github.NewTest(s.URL)

	// setup vela mock server
	engine.Use(func(c *gin.Context) { c.Set("logger", logrus.NewEntry(logrus.StandardLogger())) })
	engine.Use(func(c *gin.Context) { c.Set("secret", secret) })
	engine.Use(func(c *gin.Context) { c.Set("token-manager", tm) })
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(func(c *gin.Context) { scm.ToContext(c, client) })
	engine.Use(claims.Establish())
	engine.Use(user.Establish())
	engine.Use(org.Establish())
	engine.Use(repo.Establish())
	engine.Use(MustRead())
	engine.GET("/test/:org/:repo", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	s1 := httptest.NewServer(engine)
	defer s1.Close()

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("MustRead returned %v, want %v", resp.Code, http.StatusOK)
	}
}

func TestPerm_MustRead_RepoPublic(t *testing.T) {
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
	u.SetName("foob")
	u.SetToken("bar")
	u.SetAdmin(false)

	mto := &token.MintTokenOpts{
		User:          u,
		TokenDuration: tm.UserAccessTokenDuration,
		TokenType:     constants.UserAccessTokenType,
	}

	tok, _ := tm.MintToken(mto)

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)

	// setup database
	db, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create test database engine: %v", err)
	}

	defer func() {
		_ = db.DeleteRepo(_context.TODO(), r)
		_ = db.DeleteUser(_context.TODO(), u)
		db.Close()
	}()

	_, _ = db.CreateRepo(_context.TODO(), r)
	_, _ = db.CreateUser(_context.TODO(), u)

	context.Request, _ = http.NewRequest(http.MethodGet, "/test/foo/bar", nil)
	context.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tok))

	// setup github mock server
	engine.GET("/api/v3/repos/:org/:repo/collaborators/:username/permission", func(c *gin.Context) {
		c.String(http.StatusOK, permNonePayload)
	})
	engine.GET("/api/v3/user", func(c *gin.Context) {
		c.String(http.StatusOK, userPayload)
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup client
	client, _ := github.NewTest(s.URL)

	// setup vela mock server
	engine.Use(func(c *gin.Context) { c.Set("logger", logrus.NewEntry(logrus.StandardLogger())) })
	engine.Use(func(c *gin.Context) { c.Set("secret", secret) })
	engine.Use(func(c *gin.Context) { c.Set("token-manager", tm) })
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(func(c *gin.Context) { scm.ToContext(c, client) })
	engine.Use(claims.Establish())
	engine.Use(user.Establish())
	engine.Use(org.Establish())
	engine.Use(repo.Establish())
	engine.Use(MustRead())
	engine.GET("/test/:org/:repo", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	s1 := httptest.NewServer(engine)
	defer s1.Close()

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("MustRead returned %v, want %v", resp.Code, http.StatusOK)
	}
}

func TestPerm_MustRead_NotRead(t *testing.T) {
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
	r.SetVisibility("private")

	u := new(api.User)
	u.SetID(1)
	u.SetName("foob")
	u.SetToken("bar")
	u.SetAdmin(false)

	mto := &token.MintTokenOpts{
		User:          u,
		TokenDuration: tm.UserAccessTokenDuration,
		TokenType:     constants.UserAccessTokenType,
	}

	tok, _ := tm.MintToken(mto)

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)

	// setup database
	db, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create test database engine: %v", err)
	}

	defer func() {
		_ = db.DeleteRepo(_context.TODO(), r)
		_ = db.DeleteUser(_context.TODO(), u)
		db.Close()
	}()

	_, _ = db.CreateRepo(_context.TODO(), r)
	_, _ = db.CreateUser(_context.TODO(), u)

	context.Request, _ = http.NewRequest(http.MethodGet, "/test/foo/bar", nil)
	context.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tok))

	// setup github mock server
	engine.GET("/api/v3/repos/:org/:repo/collaborators/:username/permission", func(c *gin.Context) {
		c.String(http.StatusOK, permNonePayload)
	})
	engine.GET("/api/v3/user", func(c *gin.Context) {
		c.String(http.StatusOK, userPayload)
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup client
	client, _ := github.NewTest(s.URL)

	// setup vela mock server
	engine.Use(func(c *gin.Context) { c.Set("logger", logrus.NewEntry(logrus.StandardLogger())) })
	engine.Use(func(c *gin.Context) { c.Set("secret", secret) })
	engine.Use(func(c *gin.Context) { c.Set("token-manager", tm) })
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(func(c *gin.Context) { scm.ToContext(c, client) })
	engine.Use(claims.Establish())
	engine.Use(user.Establish())
	engine.Use(org.Establish())
	engine.Use(repo.Establish())
	engine.Use(MustRead())
	engine.GET("/test/:org/:repo", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	s1 := httptest.NewServer(engine)
	defer s1.Close()

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusUnauthorized {
		t.Errorf("MustRead returned %v, want %v", resp.Code, http.StatusUnauthorized)
	}
}

const permAdminPayload = `
{
  "permission": "admin",
  "role_name": "admin",
  "user": {
    "login": "foo",
    "id": 1,
    "node_id": "MDQ6VXNlcjE=",
    "avatar_url": "https://github.com/images/error/octocat_happy.gif",
    "gravatar_id": "",
    "url": "https://api.github.com/users/foo",
    "html_url": "https://github.com/octocat",
    "followers_url": "https://api.github.com/users/foo/followers",
    "following_url": "https://api.github.com/users/foo/following{/other_user}",
    "gists_url": "https://api.github.com/users/foo/gists{/gist_id}",
    "starred_url": "https://api.github.com/users/foo/starred{/org}{/repo}",
    "subscriptions_url": "https://api.github.com/users/foo/subscriptions",
    "orgs_url": "https://api.github.com/users/foo/orgs",
    "repos_url": "https://api.github.com/users/foo/repos",
    "events_url": "https://api.github.com/users/foo/events{/privacy}",
    "received_events_url": "https://api.github.com/users/foo/received_events",
    "type": "User",
    "site_admin": false
  }
}
`

const permWritePayload = `
{
  "permission": "write",
  "role_name": "maintain",
  "user": {
    "login": "foo",
    "id": 1,
    "node_id": "MDQ6VXNlcjE=",
    "avatar_url": "https://github.com/images/error/octocat_happy.gif",
    "gravatar_id": "",
    "url": "https://api.github.com/users/foo",
    "html_url": "https://github.com/octocat",
    "followers_url": "https://api.github.com/users/foo/followers",
    "following_url": "https://api.github.com/users/foo/following{/other_user}",
    "gists_url": "https://api.github.com/users/foo/gists{/gist_id}",
    "starred_url": "https://api.github.com/users/foo/starred{/org}{/repo}",
    "subscriptions_url": "https://api.github.com/users/foo/subscriptions",
    "orgs_url": "https://api.github.com/users/foo/orgs",
    "repos_url": "https://api.github.com/users/foo/repos",
    "events_url": "https://api.github.com/users/foo/events{/privacy}",
    "received_events_url": "https://api.github.com/users/foo/received_events",
    "type": "User",
    "site_admin": false
  }
}
`

const permReadPayload = `
{
  "permission": "read",
  "role_name": "triage",
  "user": {
    "login": "foo",
    "id": 1,
    "node_id": "MDQ6VXNlcjE=",
    "avatar_url": "https://github.com/images/error/octocat_happy.gif",
    "gravatar_id": "",
    "url": "https://api.github.com/users/foo",
    "html_url": "https://github.com/octocat",
    "followers_url": "https://api.github.com/users/foo/followers",
    "following_url": "https://api.github.com/users/foo/following{/other_user}",
    "gists_url": "https://api.github.com/users/foo/gists{/gist_id}",
    "starred_url": "https://api.github.com/users/foo/starred{/org}{/repo}",
    "subscriptions_url": "https://api.github.com/users/foo/subscriptions",
    "orgs_url": "https://api.github.com/users/foo/orgs",
    "repos_url": "https://api.github.com/users/foo/repos",
    "events_url": "https://api.github.com/users/foo/events{/privacy}",
    "received_events_url": "https://api.github.com/users/foo/received_events",
    "type": "User",
    "site_admin": false
  }
}
`

const permNonePayload = `
{
  "permission": "none",
  "role_name": "none",
  "user": {
    "login": "foo",
    "id": 1,
    "node_id": "MDQ6VXNlcjE=",
    "avatar_url": "https://github.com/images/error/octocat_happy.gif",
    "gravatar_id": "",
    "url": "https://api.github.com/users/foo",
    "html_url": "https://github.com/octocat",
    "followers_url": "https://api.github.com/users/foo/followers",
    "following_url": "https://api.github.com/users/foo/following{/other_user}",
    "gists_url": "https://api.github.com/users/foo/gists{/gist_id}",
    "starred_url": "https://api.github.com/users/foo/starred{/org}{/repo}",
    "subscriptions_url": "https://api.github.com/users/foo/subscriptions",
    "orgs_url": "https://api.github.com/users/foo/orgs",
    "repos_url": "https://api.github.com/users/foo/repos",
    "events_url": "https://api.github.com/users/foo/events{/privacy}",
    "received_events_url": "https://api.github.com/users/foo/received_events",
    "type": "User",
    "site_admin": false
  }
}
`

const userPayload = `
{
  "login": "foob",
  "id": 1,
  "node_id": "MDQ6VXNlcjE=",
  "avatar_url": "https://github.com/images/error/octocat_happy.gif",
  "gravatar_id": "",
  "url": "https://api.github.com/users/foo",
  "html_url": "https://github.com/octocat",
  "followers_url": "https://api.github.com/users/foo/followers",
  "following_url": "https://api.github.com/users/foo/following{/other_user}",
  "gists_url": "https://api.github.com/users/foo/gists{/gist_id}",
  "starred_url": "https://api.github.com/users/foo/starred{/org}{/repo}",
  "subscriptions_url": "https://api.github.com/users/foo/subscriptions",
  "orgs_url": "https://api.github.com/users/foo/orgs",
  "repos_url": "https://api.github.com/users/foo/repos",
  "events_url": "https://api.github.com/users/foo/events{/privacy}",
  "received_events_url": "https://api.github.com/users/foo/received_events",
  "type": "User",
  "site_admin": false,
  "name": "monalisa foo",
  "company": "GitHub",
  "blog": "https://github.com/blog",
  "location": "San Francisco",
  "email": "foo@github.com",
  "hireable": false,
  "bio": "There once was...",
  "public_repos": 2,
  "public_gists": 1,
  "followers": 20,
  "following": 0,
  "created_at": "2008-01-14T04:33:35Z",
  "updated_at": "2008-01-14T04:33:35Z",
  "private_gists": 81,
  "total_private_repos": 100,
  "owned_private_repos": 100,
  "disk_usage": 10000,
  "collaborators": 8,
  "two_factor_authentication": true,
  "plan": {
    "name": "Medium",
    "space": 400,
    "private_repos": 20,
    "collaborators": 0
  }
}
`
