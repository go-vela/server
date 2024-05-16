// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/internal/token"
)

// helper function to setup the tokenmanager from the CLI arguments.
func setupTokenManager(c *cli.Context, db database.Interface) (*token.Manager, error) {
	logrus.Debug("Creating token manager from CLI configuration")

	tm := &token.Manager{
		PrivateKeyHMAC:              c.String("vela-server-private-key"),
		UserAccessTokenDuration:     c.Duration("user-access-token-duration"),
		UserRefreshTokenDuration:    c.Duration("user-refresh-token-duration"),
		BuildTokenBufferDuration:    c.Duration("build-token-buffer-duration"),
		WorkerAuthTokenDuration:     c.Duration("worker-auth-token-duration"),
		WorkerRegisterTokenDuration: c.Duration("worker-register-token-duration"),
		IDTokenDuration:             c.Duration("id-token-duration"),
		Issuer:                      c.String("server-addr"),
	}

	// generate a new RSA key pair
	err := tm.GenerateRSA(db)
	if err != nil {
		return nil, err
	}

	return tm, nil
}
