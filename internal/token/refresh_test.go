// SPDX-License-Identifier: Apache-2.0

package token

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/database"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/golang-jwt/jwt/v5"
)

func TestTokenManager_Refresh(t *testing.T) {
	// setup types
	u := new(library.User)
	u.SetID(1)
	u.SetName("foo")
	u.SetToken("bar")
	u.SetHash("baz")

	tm := &Manager{
		PrivateKey:               "123abc",
		SignMethod:               jwt.SigningMethodHS256,
		UserAccessTokenDuration:  time.Minute * 5,
		UserRefreshTokenDuration: time.Minute * 30,
	}

	mto := &MintTokenOpts{
		User:          u,
		TokenType:     constants.UserRefreshTokenType,
		TokenDuration: tm.UserRefreshTokenDuration,
	}

	rt, err := tm.MintToken(mto)
	if err != nil {
		t.Errorf("unable to create refresh token")
	}

	u.SetRefreshToken(rt)

	// setup database
	db, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create test database engine: %v", err)
	}

	defer func() {
		db.DeleteUser(context.TODO(), u)
		db.Close()
	}()

	_, _ = db.CreateUser(context.TODO(), u)

	// set up context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(resp)
	context.Set("database", db)
	req, _ := http.NewRequestWithContext(context, "", "", nil)
	context.Request = req

	// run tests
	got, err := tm.Refresh(context, rt)
	if err != nil {
		t.Error("Refresh should not error")
	}

	if len(got) == 0 {
		t.Errorf("Refresh should have returned an access token")
	}
}

func TestTokenManager_Refresh_Expired(t *testing.T) {
	// setup types
	u := new(library.User)
	u.SetID(1)
	u.SetName("foo")
	u.SetToken("bar")
	u.SetHash("baz")

	tm := &Manager{
		PrivateKey:               "123abc",
		SignMethod:               jwt.SigningMethodHS256,
		UserAccessTokenDuration:  time.Minute * 5,
		UserRefreshTokenDuration: time.Minute * 30,
	}

	mto := &MintTokenOpts{
		User:          u,
		TokenType:     constants.UserRefreshTokenType,
		TokenDuration: time.Minute * -1,
	}

	rt, err := tm.MintToken(mto)
	if err != nil {
		t.Errorf("unable to create refresh token")
	}

	u.SetRefreshToken(rt)

	// setup database
	db, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create test database engine: %v", err)
	}

	defer func() {
		db.DeleteUser(context.TODO(), u)
		db.Close()
	}()

	_, _ = db.CreateUser(context.TODO(), u)

	// set up context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(resp)
	context.Set("database", db)
	req, _ := http.NewRequestWithContext(context, "", "", nil)
	context.Request = req

	// run tests
	_, err = tm.Refresh(context, rt)
	if err == nil {
		t.Error("Refresh with expired token should error")
	}
}
