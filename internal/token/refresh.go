// SPDX-License-Identifier: Apache-2.0

package token

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database"
)

// Refresh returns a new access token, if the provided refreshToken is valid.
func (tm *Manager) Refresh(c *gin.Context, refreshToken string) (string, error) {
	// capture middleware values
	ctx := c.Request.Context()

	// retrieve claims from token
	claims, err := tm.ParseToken(refreshToken)
	if err != nil {
		return "", err
	}

	// look up user in database given claims subject
	u, err := database.FromContext(c).GetUserForName(ctx, claims.Subject)
	if err != nil {
		return "", fmt.Errorf("unable to retrieve user %s from database from claims subject: %w", claims.Subject, err)
	}

	// options for user access token minting
	amto := &MintTokenOpts{
		User:          u,
		TokenType:     constants.UserAccessTokenType,
		TokenDuration: tm.UserAccessTokenDuration,
	}

	// create a new access token
	at, err := tm.MintToken(amto)
	if err != nil {
		return "", err
	}

	return at, nil
}
