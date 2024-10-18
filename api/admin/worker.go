// SPDX-License-Identifier: Apache-2.0

package admin

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/internal/token"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/library"
)

// swagger:operation POST /api/v1/admin/workers/{worker}/register admin RegisterToken
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
//   '500':
//     description: Unexpected server error
//     schema:
//       "$ref": "#/definitions/Error"

// RegisterToken represents the API handler to
// generate a registration token for onboarding a worker.
func RegisterToken(c *gin.Context) {
	l := c.MustGet("logger").(*logrus.Entry)

	// capture middleware values
	host := util.PathParameter(c, "worker")

	tm := c.MustGet("token-manager").(*token.Manager)
	rmto := &token.MintTokenOpts{
		Hostname:      host,
		TokenType:     constants.WorkerRegisterTokenType,
		TokenDuration: tm.WorkerRegisterTokenDuration,
	}

	l.Debug("platform admin: generating worker registration token")

	rt, err := tm.MintToken(rmto)
	if err != nil {
		retErr := fmt.Errorf("unable to generate registration token: %w", err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	l.Infof("platform admin: generated worker registration token for %s", host)

	c.JSON(http.StatusOK, library.Token{Token: &rt})
}
