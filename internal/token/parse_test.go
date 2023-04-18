// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package token

import (
	"reflect"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"

	jwt "github.com/golang-jwt/jwt/v5"
)

func TestTokenManager_ParseToken(t *testing.T) {
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

	now := time.Now()

	tests := []struct {
		TokenType string
		Mto       *MintTokenOpts
		Want      *Claims
	}{
		{
			TokenType: constants.UserAccessTokenType,
			Mto: &MintTokenOpts{
				User:          u,
				TokenType:     constants.UserAccessTokenType,
				TokenDuration: tm.UserAccessTokenDuration,
			},
			Want: &Claims{
				IsActive:  u.GetActive(),
				IsAdmin:   u.GetAdmin(),
				TokenType: constants.UserAccessTokenType,
				RegisteredClaims: jwt.RegisteredClaims{
					Subject:   u.GetName(),
					IssuedAt:  jwt.NewNumericDate(now),
					ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute * 5)),
				},
			},
		},
		{
			TokenType: constants.UserRefreshTokenType,
			Mto: &MintTokenOpts{
				User:          u,
				TokenType:     constants.UserRefreshTokenType,
				TokenDuration: tm.UserRefreshTokenDuration,
			},
			Want: &Claims{
				IsActive:  u.GetActive(),
				IsAdmin:   u.GetAdmin(),
				TokenType: constants.UserRefreshTokenType,
				RegisteredClaims: jwt.RegisteredClaims{
					Subject:   u.GetName(),
					IssuedAt:  jwt.NewNumericDate(now),
					ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute * 30)),
				},
			},
		},
		{
			TokenType: constants.WorkerBuildTokenType,
			Mto: &MintTokenOpts{
				BuildID:       1,
				Repo:          "foo/bar",
				Hostname:      "worker",
				TokenType:     constants.WorkerBuildTokenType,
				TokenDuration: time.Minute * 90,
			},
			Want: &Claims{
				BuildID:   1,
				Repo:      "foo/bar",
				TokenType: constants.WorkerBuildTokenType,
				RegisteredClaims: jwt.RegisteredClaims{
					Subject:   "worker",
					IssuedAt:  jwt.NewNumericDate(now),
					ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute * 90)),
				},
			},
		},
	}

	gin.SetMode(gin.TestMode)

	for _, tt := range tests {
		t.Run(tt.TokenType, func(t *testing.T) {
			tkn, err := tm.MintToken(tt.Mto)
			if err != nil {
				t.Errorf("Unable to create token: %v", err)
			}
			// run test
			got, err := tm.ParseToken(tkn)
			if err != nil {
				t.Errorf("Parse returned err: %v", err)
			}

			if !reflect.DeepEqual(got, tt.Want) {
				t.Errorf("Parse is %v, want %v", got, tt.Want)
			}
		})
	}
}

func TestTokenManager_ParseToken_Error_NoParse(t *testing.T) {
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

	// run test
	got, err := tm.ParseToken("!@#$%^&*()")
	if err == nil {
		t.Errorf("Parse should have returned err")
	}

	if got != nil {
		t.Errorf("Parse is %v, want nil", got)
	}
}

func TestTokenManager_ParseToken_Expired(t *testing.T) {
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
		TokenType:     constants.UserAccessTokenType,
		TokenDuration: time.Minute * -1,
	}

	tkn, err := tm.MintToken(mto)
	if err != nil {
		t.Errorf("Unable to create token: %v", err)
	}

	// run test
	_, err = tm.ParseToken(tkn)
	if err == nil {
		t.Errorf("Parse should return error due to expiration")
	}
}

func TestTokenManager_ParseToken_NoSubject(t *testing.T) {
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

	claims := &Claims{
		IsActive:  u.GetActive(),
		IsAdmin:   u.GetAdmin(),
		TokenType: constants.UserRefreshTokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now()),
		},
	}
	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := tkn.SignedString([]byte(tm.PrivateKey))
	if err != nil {
		t.Errorf("Unable to create test token: %v", err)
	}

	// run test
	got, err := tm.ParseToken(token)
	if err == nil {
		t.Errorf("Parse should have returned err")
	}

	if got != nil {
		t.Errorf("Parse is %v, want nil", got)
	}
}

func TestTokenManager_ParseToken_Error_InvalidSignature(t *testing.T) {
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

	claims := &Claims{
		IsActive:  u.GetActive(),
		IsAdmin:   u.GetAdmin(),
		TokenType: constants.UserAccessTokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   u.GetName(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 1)),
		},
	}
	tkn := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	token, err := tkn.SignedString([]byte(tm.PrivateKey))
	if err != nil {
		t.Errorf("Unable to create test token: %v", err)
	}

	// run test
	got, err := tm.ParseToken(token)
	if err == nil {
		t.Errorf("Parse should have returned err")
	}

	if got != nil {
		t.Errorf("Parse is %v, want nil", got)
	}
}

func TestToken_Parse_AccessToken_NoExpiration(t *testing.T) {
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

	claims := &Claims{
		TokenType: constants.UserAccessTokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: "user",
		},
	}
	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := tkn.SignedString([]byte(u.GetHash()))
	if err != nil {
		t.Errorf("Unable to create test token: %v", err)
	}

	// run test
	got, err := tm.ParseToken(token)
	if err == nil {
		t.Errorf("Parse should have returned err")
	}

	if got != nil {
		t.Errorf("Parse is %v, want nil", got)
	}
}
