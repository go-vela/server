// SPDX-License-Identifier: Apache-2.0

package claims

import (
	_context "context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/internal/token"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/golang-jwt/jwt/v5"
)

func TestClaims_Retrieve(t *testing.T) {
	// setup types
	now := time.Now()
	want := &token.Claims{
		TokenType: constants.UserAccessTokenType,
		IsAdmin:   false,
		IsActive:  true,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "octocat",
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute * 1)),
		},
	}

	// setup context
	gin.SetMode(gin.TestMode)

	context, _ := gin.CreateTestContext(nil)
	ToContext(context, want)

	// run test
	got := Retrieve(context)

	if got != want {
		t.Errorf("Retrieve is %v, want %v", got, want)
	}
}

func TestClaims_Establish(t *testing.T) {
	// setup types
	user := new(library.User)
	user.SetID(1)
	user.SetName("foo")
	user.SetRefreshToken("fresh")
	user.SetToken("bar")
	user.SetHash("baz")
	user.SetActive(true)
	user.SetAdmin(false)
	user.SetFavorites([]string{})

	tm := &token.Manager{
		PrivateKey:                  "123abc",
		SignMethod:                  jwt.SigningMethodHS256,
		UserAccessTokenDuration:     time.Minute * 5,
		UserRefreshTokenDuration:    time.Minute * 30,
		WorkerAuthTokenDuration:     time.Minute * 20,
		WorkerRegisterTokenDuration: time.Minute * 1,
	}

	now := time.Now()

	tests := []struct {
		TokenType  string
		WantClaims *token.Claims
		Mto        *token.MintTokenOpts
		CtxRequest string
		Endpoint   string
	}{
		{
			TokenType: constants.UserAccessTokenType,
			WantClaims: &token.Claims{
				TokenType: constants.UserAccessTokenType,
				IsAdmin:   false,
				IsActive:  true,
				RegisteredClaims: jwt.RegisteredClaims{
					Subject:   "foo",
					IssuedAt:  jwt.NewNumericDate(now),
					ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute * 5)),
				},
			},
			Mto: &token.MintTokenOpts{
				User:          user,
				TokenDuration: tm.UserAccessTokenDuration,
				TokenType:     constants.UserAccessTokenType,
			},
			CtxRequest: "/repos/foo/bar/builds/1",
			Endpoint:   "repos/:org/:repo/builds/:build",
		},
		{
			TokenType: constants.WorkerBuildTokenType,
			WantClaims: &token.Claims{
				TokenType: constants.WorkerBuildTokenType,
				BuildID:   1,
				Repo:      "foo/bar",
				RegisteredClaims: jwt.RegisteredClaims{
					Subject:   "host",
					IssuedAt:  jwt.NewNumericDate(now),
					ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute * 35)),
				},
			},
			Mto: &token.MintTokenOpts{
				Hostname:      "host",
				BuildID:       1,
				Repo:          "foo/bar",
				TokenDuration: time.Minute * 35,
				TokenType:     constants.WorkerBuildTokenType,
			},
			CtxRequest: "/repos/foo/bar/builds/1",
			Endpoint:   "repos/:org/:repo/builds/:build",
		},
		{
			TokenType: constants.WorkerAuthTokenType,
			WantClaims: &token.Claims{
				TokenType: constants.WorkerAuthTokenType,
				RegisteredClaims: jwt.RegisteredClaims{
					Subject:   "host",
					IssuedAt:  jwt.NewNumericDate(now),
					ExpiresAt: jwt.NewNumericDate(now.Add(tm.WorkerAuthTokenDuration)),
				},
			},
			Mto: &token.MintTokenOpts{
				Hostname:      "host",
				TokenDuration: tm.WorkerAuthTokenDuration,
				TokenType:     constants.WorkerAuthTokenType,
			},
			CtxRequest: "/workers/host",
			Endpoint:   "/workers/:hostname",
		},
		{
			TokenType: constants.WorkerRegisterTokenType,
			WantClaims: &token.Claims{
				TokenType: constants.WorkerRegisterTokenType,
				RegisteredClaims: jwt.RegisteredClaims{
					Subject:   "host",
					IssuedAt:  jwt.NewNumericDate(now),
					ExpiresAt: jwt.NewNumericDate(now.Add(tm.WorkerRegisterTokenDuration)),
				},
			},
			Mto: &token.MintTokenOpts{
				Hostname:      "host",
				TokenDuration: tm.WorkerRegisterTokenDuration,
				TokenType:     constants.WorkerRegisterTokenType,
			},
			CtxRequest: "/workers/host/register",
			Endpoint:   "workers/:hostname/register",
		},
		{
			TokenType: constants.ServerWorkerTokenType,
			WantClaims: &token.Claims{
				TokenType: constants.ServerWorkerTokenType,
				RegisteredClaims: jwt.RegisteredClaims{
					Subject: "vela-worker",
				},
			},
			CtxRequest: "/repos/foo/bar/builds/1",
			Endpoint:   "repos/:org/:repo/builds/:build",
		},
	}

	got := new(token.Claims)

	gin.SetMode(gin.TestMode)

	for _, tt := range tests {
		t.Run(tt.TokenType, func(t *testing.T) {
			resp := httptest.NewRecorder()
			context, engine := gin.CreateTestContext(resp)
			context.Request, _ = http.NewRequest(http.MethodPut, tt.CtxRequest, nil)

			var tkn string

			if strings.EqualFold(tt.TokenType, constants.ServerWorkerTokenType) {
				tkn = "very-secret"
				engine.Use(func(c *gin.Context) { c.Set("secret", "very-secret") })
			} else {
				tkn, _ = tm.MintToken(tt.Mto)
			}

			context.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tkn))

			// setup context
			gin.SetMode(gin.TestMode)

			// setup vela mock server
			engine.Use(func(c *gin.Context) { c.Set("token-manager", tm) })
			engine.Use(Establish())
			engine.PUT(tt.Endpoint, func(c *gin.Context) {
				got = Retrieve(c)

				c.Status(http.StatusOK)
			})

			s1 := httptest.NewServer(engine)

			// run test
			engine.ServeHTTP(context.Writer, context.Request)

			if resp.Code != http.StatusOK {
				t.Errorf("Establish returned %v, want %v", resp.Code, http.StatusOK)
			}

			if !reflect.DeepEqual(got, tt.WantClaims) {
				t.Errorf("Establish is %v, want %v", got, tt.WantClaims)
			}

			s1.Close()
		})
	}
}

func TestClaims_Establish_NoToken(t *testing.T) {
	// setup types
	tm := &token.Manager{
		PrivateKey:               "123abc",
		SignMethod:               jwt.SigningMethodHS256,
		UserAccessTokenDuration:  time.Minute * 5,
		UserRefreshTokenDuration: time.Minute * 30,
	}
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/workers/host", nil)

	engine.Use(func(c *gin.Context) { c.Set("token-manager", tm) })
	engine.Use(Establish())

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusUnauthorized {
		t.Errorf("Establish returned %v, want %v", resp.Code, http.StatusUnauthorized)
	}
}

func TestClaims_Establish_BadToken(t *testing.T) {
	// setup types
	tm := &token.Manager{
		PrivateKey:               "123abc",
		SignMethod:               jwt.SigningMethodHS256,
		UserAccessTokenDuration:  time.Minute * 5,
		UserRefreshTokenDuration: time.Minute * 30,
	}
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/workers/host", nil)

	u := new(library.User)
	u.SetID(1)
	u.SetName("octocat")
	u.SetHash("abc")

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

	mto := &token.MintTokenOpts{
		User:          u,
		TokenDuration: time.Minute * -1,
		TokenType:     constants.UserRefreshTokenType,
	}

	tkn, _ := tm.MintToken(mto)

	context.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tkn))

	engine.Use(func(c *gin.Context) { c.Set("token-manager", tm) })
	engine.Use(func(c *gin.Context) { c.Set("secret", "very-secret") })
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(Establish())

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusUnauthorized {
		t.Errorf("Establish returned %v, want %v", resp.Code, http.StatusUnauthorized)
	}
}
