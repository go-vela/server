// SPDX-License-Identifier: Apache-2.0

package token

import (
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/internal"
	"github.com/go-vela/types/constants"
)

func TestToken_Compose(t *testing.T) {
	// setup types
	u := new(api.User)
	u.SetID(1)
	u.SetName("foo")
	u.SetToken("bar")

	tm := &Manager{
		PrivateKey:               "123abc",
		SignMethod:               jwt.SigningMethodHS256,
		UserAccessTokenDuration:  time.Minute * 5,
		UserRefreshTokenDuration: time.Minute * 30,
	}

	d := time.Minute * 5
	now := time.Now()
	exp := now.Add(d)

	claims := &Claims{
		IsActive:  u.GetActive(),
		IsAdmin:   u.GetAdmin(),
		TokenType: constants.UserAccessTokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   u.GetName(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}

	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	want, err := tkn.SignedString([]byte(tm.PrivateKey))
	if err != nil {
		t.Errorf("Unable to create test token: %v", err)
	}

	m := &internal.Metadata{
		Vela: &internal.Vela{
			AccessTokenDuration: d,
		},
	}

	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(resp)
	context.Set("metadata", m)
	context.Set("securecookie", false)

	// run test
	_, got, err := tm.Compose(context, u)
	if err != nil {
		t.Errorf("Compose returned err: %v", err)
	}

	if !strings.EqualFold(got, want) {
		t.Errorf("Compose is %v, want %v", got, want)
	}
}
