// SPDX-License-Identifier: Apache-2.0

package worker

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/internal/token"
	"github.com/go-vela/server/router/middleware/claims"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
)

// swagger:operation POST /api/v1/workers workers CreateWorker
//
// Create a worker
//
// ---
// produces:
// - application/json
// parameters:
// - in: body
//   name: body
//   description: Worker object to create
//   required: true
//   schema:
//     "$ref": "#/definitions/Worker"
// security:
//   - ApiKeyAuth: []
// responses:
//   '201':
//     description: Successfully created the worker and retrieved auth token
//     schema:
//       "$ref": "#/definitions/Token"
//   '400':
//     description: Invalid request payload
//     schema:
//       "$ref": "#/definitions/Error"
//   '401':
//     description: Unauthorized
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unexpected server error
//     schema:
//       "$ref": "#/definitions/Error"

// CreateWorker represents the API handler to
// create a worker.
func CreateWorker(c *gin.Context) {
	// capture middleware values
	u := user.Retrieve(c)
	cl := claims.Retrieve(c)
	ctx := c.Request.Context()

	// capture body from API request
	input := new(types.Worker)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for new worker: %w", err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// verify input host name matches worker hostname
	if !strings.EqualFold(cl.TokenType, constants.ServerWorkerTokenType) && !strings.EqualFold(cl.Subject, input.GetHostname()) {
		retErr := fmt.Errorf("unable to add worker; claims subject %s does not match worker hostname %s", cl.Subject, input.GetHostname())

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	input.SetLastCheckedIn(time.Now().Unix())

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"user":   u.GetName(),
		"worker": input.GetHostname(),
	}).Debugf("creating new worker %s", input.GetHostname())

	_, err = database.FromContext(c).CreateWorker(ctx, input)
	if err != nil {
		retErr := fmt.Errorf("unable to create worker: %w", err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	switch cl.TokenType {
	// if symmetric token configured, send back symmetric token
	case constants.ServerWorkerTokenType:
		if secret, ok := c.Value("secret").(string); ok {
			tkn := new(library.Token)
			tkn.SetToken(secret)
			c.JSON(http.StatusCreated, tkn)

			return
		}

		retErr := fmt.Errorf("symmetric token provided but not configured in server")
		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	// if worker register token, send back auth token
	default:
		tm := c.MustGet("token-manager").(*token.Manager)

		wmto := &token.MintTokenOpts{
			TokenType:     constants.WorkerAuthTokenType,
			TokenDuration: tm.WorkerAuthTokenDuration,
			Hostname:      cl.Subject,
		}

		tkn := new(library.Token)

		wt, err := tm.MintToken(wmto)
		if err != nil {
			retErr := fmt.Errorf("unable to generate auth token for worker %s: %w", input.GetHostname(), err)

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}

		tkn.SetToken(wt)

		c.JSON(http.StatusCreated, tkn)
	}
}
