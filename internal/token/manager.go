// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package token

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Manager struct {
	// PrivateKey key used to sign tokens
	PrivateKey string

	// SignMethod method to sign tokens
	SignMethod jwt.SigningMethod

	// UserAccessTokenDuration specifies the token duration to use for users
	UserAccessTokenDuration time.Duration

	// UserRefreshTokenDuration specifies the token duration for user refresh
	UserRefreshTokenDuration time.Duration

	// BuildTokenBufferDuration specifies the additional token duration of build tokens beyond repo timeout
	BuildTokenBufferDuration time.Duration
}
