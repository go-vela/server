// SPDX-License-Identifier: Apache-2.0

package token

import (
	"crypto/rsa"
	"time"
)

type RSAKeySet struct {
	PrivateKey *rsa.PrivateKey
	KID        string
}

type Manager struct {
	// PrivateKeyHMAC is the private key used to sign and validate closed-system tokens
	PrivateKeyHMAC string

	// RSAKeySet is the private key used to sign and validate open-system tokens (OIDC)
	RSAKeySet RSAKeySet

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

	// IDTokenDuration specifies the token duration for ID tokens
	IDTokenDuration time.Duration

	// Issuer specifies the issuer of the token
	Issuer string
}
