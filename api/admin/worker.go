// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package admin

import (
	"fmt"
	"net/http"

	"github.com/go-vela/server/internal/token"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// swagger:operation POST /api/v1/admin/workers/{worker}/register-token admin RegisterToken
//
// Get a worker registration token
//
// ---
// produces:
// - application/json
// parameters:
// - in: path
//   name: worker
//   description: Hostname of the worker
//   required: true
//   type: string
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully generated registration token
//     schema:
//       "$ref": "#/definitions/Token"
//   '401':
//     description: Unauthorized
//     schema:
//       "$ref": "#/definitions/Error"

// RegisterToken represents the API handler to
// generate a registration token for onboarding a worker.
func RegisterToken(c *gin.Context) {
	logrus.Info("Admin: generating registration token")

	host := util.PathParameter(c, "worker")

	tm := c.MustGet("token-manager").(*token.Manager)
	rmto := &token.MintTokenOpts{
		Hostname:      host,
		TokenType:     constants.WorkerRegisterTokenType,
		TokenDuration: tm.WorkerRegisterTokenDuration,
	}

	rt, err := tm.MintToken(rmto)
	if err != nil {
		retErr := fmt.Errorf("unable to generate registration token: %w", err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, library.Token{Token: &rt})
}
