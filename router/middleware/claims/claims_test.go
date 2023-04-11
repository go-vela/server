// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package claims

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/database/sqlite"
	"github.com/go-vela/server/internal/token"
	"github.com/golang-jwt/jwt/v4"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"

	"github.com/gin-gonic/gin"
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
		PrivateKey:               "123abc",
		SignMethod:               jwt.SigningMethodHS256,
		UserAccessTokenDuration:  time.Minute * 5,
		UserRefreshTokenDuration: time.Minute * 30,
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
					IssuedAt:  nil,
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
					IssuedAt:  nil,
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

	// setup database
	db, _ := sqlite.NewTest()

	defer func() {
		db.Sqlite.Exec("delete from users;")
		_sql, _ := db.Sqlite.DB()
		_sql.Close()
	}()

	_ = db.CreateUser(user)

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
			} else {
				tkn, _ = tm.MintToken(tt.Mto)
			}

			context.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tkn))

			// setup context
			gin.SetMode(gin.TestMode)

			// setup vela mock server
			engine.Use(func(c *gin.Context) { c.Set("token-manager", tm) })
			engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
			engine.Use(func(c *gin.Context) { c.Set("secret", "very-secret") })
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
	db, _ := sqlite.NewTest()

	defer func() {
		db.Sqlite.Exec("delete from users;")
		_sql, _ := db.Sqlite.DB()
		_sql.Close()
	}()

	_ = db.CreateUser(u)

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
