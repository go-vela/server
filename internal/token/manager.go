// SPDX-License-Identifier: Apache-2.0

package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
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

	// WorkerAuthTokenDuration specifies the token duration for worker auth (check in)
	WorkerAuthTokenDuration time.Duration

	// WorkerRegisterTokenDuration specifies the token duration for worker register
	WorkerRegisterTokenDuration time.Duration
}
