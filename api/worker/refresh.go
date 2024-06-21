// SPDX-License-Identifier: Apache-2.0

package worker

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/internal/token"
	"github.com/go-vela/server/router/middleware/claims"
	"github.com/go-vela/server/router/middleware/worker"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
)

// swagger:operation POST /api/v1/workers/{worker}/refresh workers RefreshWorkerAuth
//
// Refresh authorization token for a worker
//
// ---
// produces:
// - application/json
// parameters:
// - in: path
//   name: worker
//   description: Name of the worker
//   required: true
//   type: string
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully refreshed auth
//     schema:
//       "$ref": "#/definitions/Token"
//   '400':
//     description: Invalid request payload or path
//     schema:
//       "$ref": "#/definitions/Error"
//   '401':
//     description: Unauthorized
//     schema:
//       "$ref": "#/definitions/Error"
//   '404':
//     description: Not found
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unexpected server error
//     schema:
//       "$ref": "#/definitions/Error"

// Refresh represents the API handler to
// refresh the auth token for a worker.
func Refresh(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
	w := worker.Retrieve(c)
	cl := claims.Retrieve(c)
	ctx := c.Request.Context()

	// if we are not using a symmetric token, and the subject does not match the input, request should be denied
	if !strings.EqualFold(cl.TokenType, constants.ServerWorkerTokenType) && !strings.EqualFold(cl.Subject, w.GetHostname()) {
		retErr := fmt.Errorf("unable to refresh worker auth: claims subject %s does not match worker hostname %s", cl.Subject, w.GetHostname())

		l.Warnf("attempted refresh of worker %s using token from worker %s", w.GetHostname(), cl.Subject)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// set last checked in time
	w.SetLastCheckedIn(time.Now().Unix())

	// send API call to update the worker
	_, err := database.FromContext(c).UpdateWorker(ctx, w)
	if err != nil {
		retErr := fmt.Errorf("unable to update worker %s: %w", w.GetHostname(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	l.Info("worker updated - check-in time updated")

	l.Debugf("refreshing worker %s authentication", w.GetHostname())

	switch cl.TokenType {
	// if symmetric token configured, send back symmetric token
	case constants.ServerWorkerTokenType:
		if secret, ok := c.Value("secret").(string); ok {
			tkn := new(library.Token)
			tkn.SetToken(secret)

			c.JSON(http.StatusOK, tkn)

			return
		}

		retErr := fmt.Errorf("symmetric token provided but not configured in server")
		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	// if worker auth / register token, send back auth token
	case constants.WorkerAuthTokenType, constants.WorkerRegisterTokenType:
		tm := c.MustGet("token-manager").(*token.Manager)

		wmto := &token.MintTokenOpts{
			TokenType:     constants.WorkerAuthTokenType,
			TokenDuration: tm.WorkerAuthTokenDuration,
			Hostname:      cl.Subject,
		}

		tkn := new(library.Token)

		wt, err := tm.MintToken(wmto)
		if err != nil {
			retErr := fmt.Errorf("unable to generate auth token for worker %s: %w", w.GetHostname(), err)

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}

		tkn.SetToken(wt)

		c.JSON(http.StatusOK, tkn)
	}
}
