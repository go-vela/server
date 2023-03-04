// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package claims

import (
	"testing"
	"time"

	"github.com/go-vela/server/internal/token"
	"github.com/go-vela/types/constants"
	"github.com/golang-jwt/jwt/v4"

	"github.com/gin-gonic/gin"
)

func TestClaims_FromContext(t *testing.T) {
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
	context.Set(key, want)

	// run test
	got := FromContext(context)

	if got != want {
		t.Errorf("FromContext is %v, want %v", got, want)
	}
}

func TestClaims_FromContext_Bad(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)
	context, _ := gin.CreateTestContext(nil)
	context.Set(key, nil)

	// run test
	got := FromContext(context)

	if got != nil {
		t.Errorf("FromContext is %v, want nil", got)
	}
}

func TestClaims_FromContext_WrongType(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)
	context, _ := gin.CreateTestContext(nil)
	context.Set(key, 1)

	// run test
	got := FromContext(context)

	if got != nil {
		t.Errorf("FromContext is %v, want nil", got)
	}
}

func TestClaims_FromContext_Empty(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)
	context, _ := gin.CreateTestContext(nil)

	// run test
	got := FromContext(context)

	if got != nil {
		t.Errorf("FromContext is %v, want nil", got)
	}
}

func TestClaims_ToContext(t *testing.T) {
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
	got := context.Value(key)

	if got != want {
		t.Errorf("ToContext is %v, want %v", got, want)
	}
}
