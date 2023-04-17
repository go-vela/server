// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/golang-jwt/jwt/v5"
)

// Claims struct is an extension of the JWT standard claims. It
// includes information about the user.
type Claims struct {
	BuildID   int64  `json:"build_id"`
	IsActive  bool   `json:"is_active"`
	IsAdmin   bool   `json:"is_admin"`
	Repo      string `json:"repo"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

// MintTokenOpts is a type to inform the token minter how to construct
// the token.
type MintTokenOpts struct {
	BuildID       int64
	Hostname      string
	Repo          string
	TokenDuration time.Duration
	TokenType     string
	User          *library.User
}

// MintToken mints a Vela JWT Token given a set of options.
func (tm *Manager) MintToken(mto *MintTokenOpts) (string, error) {
	// initialize claims struct
	var claims = new(Claims)

	// apply claims based on token type
	switch mto.TokenType {
	case constants.UserAccessTokenType, constants.UserRefreshTokenType:
		if mto.User == nil {
			return "", fmt.Errorf("no user provided for user access token")
		}

		claims.IsActive = mto.User.GetActive()
		claims.IsAdmin = mto.User.GetAdmin()
		claims.Subject = mto.User.GetName()

	case constants.WorkerBuildTokenType:
		if mto.BuildID == 0 {
			return "", errors.New("missing build id for build token")
		}

		if len(mto.Repo) == 0 {
			return "", errors.New("missing repo for build token")
		}

		if len(mto.Hostname) == 0 {
			return "", errors.New("missing host name for build token")
		}

		claims.BuildID = mto.BuildID
		claims.Repo = mto.Repo
		claims.Subject = mto.Hostname

	case constants.WorkerAuthTokenType, constants.WorkerRegisterTokenType:
		if len(mto.Hostname) == 0 {
			return "", fmt.Errorf("missing host name for %s token", mto.TokenType)
		}

		claims.Subject = mto.Hostname

	default:
		return "", errors.New("invalid token type")
	}

	claims.IssuedAt = jwt.NewNumericDate(time.Now())
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(mto.TokenDuration))
	claims.TokenType = mto.TokenType

	tk := jwt.NewWithClaims(tm.SignMethod, claims)

	//sign token with configured private signing key
	token, err := tk.SignedString([]byte(tm.PrivateKey))
	if err != nil {
		return "", fmt.Errorf("unable to sign token: %w", err)
	}

	return token, nil
}
