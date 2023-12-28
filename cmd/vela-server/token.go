// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/golang-jwt/jwt/v5"

	"github.com/sirupsen/logrus"

	"github.com/urfave/cli/v2"

	"github.com/go-vela/server/internal/token"
)

// helper function to setup the tokenmanager from the CLI arguments.
func setupTokenManager(c *cli.Context) *token.Manager {
	logrus.Debug("Creating token manager from CLI configuration")

	tm := &token.Manager{
		PrivateKey:                  c.String("vela-server-private-key"),
		SignMethod:                  jwt.SigningMethodHS256,
		UserAccessTokenDuration:     c.Duration("user-access-token-duration"),
		UserRefreshTokenDuration:    c.Duration("user-refresh-token-duration"),
		BuildTokenBufferDuration:    c.Duration("build-token-buffer-duration"),
		WorkerAuthTokenDuration:     c.Duration("worker-auth-token-duration"),
		WorkerRegisterTokenDuration: c.Duration("worker-register-token-duration"),
	}

	return tm
}
