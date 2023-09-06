// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

//nolint:dupl // ignore similar code with delete token
package user

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/internal/token"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// swagger:operation POST /api/v1/user/token users CreateToken
//
// Create a token for the current authenticated user
//
// ---
// produces:
// - application/json
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully created a token for the current user
//     schema:
//       "$ref": "#/definitions/Token"
//   '503':
//     description: Unable to create a token for the current user
//     schema:
//       "$ref": "#/definitions/Error"

// CreateToken represents the API handler to create
// a user token in the configured backend.
func CreateToken(c *gin.Context) {
	// capture middleware values
	u := user.Retrieve(c)
	ctx := c.Request.Context()

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"user": u.GetName(),
	}).Infof("composing token for user %s", u.GetName())

	tm := c.MustGet("token-manager").(*token.Manager)

	// compose JWT token for user
	rt, at, err := tm.Compose(c, u)
	if err != nil {
		retErr := fmt.Errorf("unable to compose token for user %s: %w", u.GetName(), err)

		util.HandleError(c, http.StatusServiceUnavailable, retErr)

		return
	}

	u.SetRefreshToken(rt)

	// send API call to update the user
	_, err = database.FromContext(c).UpdateUser(ctx, u)
	if err != nil {
		retErr := fmt.Errorf("unable to update user %s: %w", u.GetName(), err)

		util.HandleError(c, http.StatusServiceUnavailable, retErr)

		return
	}

	c.JSON(http.StatusOK, library.Token{Token: &at})
}
