// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/internal/token"
)

// helper function to setup the tokenmanager from the CLI arguments.
func setupTokenManager(ctx context.Context, c *cli.Command, db database.Interface) (*token.Manager, error) {
	logrus.Debug("creating token manager from CLI configuration")

	tm := &token.Manager{
		PrivateKeyHMAC:              c.String("vela-server-private-key"),
		UserAccessTokenDuration:     c.Duration("user-access-token-duration"),
		UserRefreshTokenDuration:    c.Duration("user-refresh-token-duration"),
		BuildTokenBufferDuration:    c.Duration("build-token-buffer-duration"),
		WorkerAuthTokenDuration:     c.Duration("worker-auth-token-duration"),
		WorkerRegisterTokenDuration: c.Duration("worker-register-token-duration"),
		IDTokenDuration:             c.Duration("id-token-duration"),
	}

	// generate a new RSA key pair
	err := tm.GenerateRSA(ctx, db)
	if err != nil {
		return nil, err
	}

	return tm, nil
}
