// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package worker

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/internal/token"
	"github.com/go-vela/server/router/middleware/claims"
	"github.com/go-vela/server/router/middleware/worker"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// swagger:operation POST /api/v1/workers/{worker}/refresh workers RefreshWorkerAuth
//
// Refresh authorization token for worker
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
//     description: Unable to refresh worker auth
//     schema:
//       "$ref": "#/definitions/Error"
//   '404':
//     description: Unable to refresh worker auth
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to refresh worker auth
//     schema:
//       "$ref": "#/definitions/Error"

// Refresh represents the API handler to
// refresh the auth token for a worker.
func Refresh(c *gin.Context) {
	// capture middleware values
	w := worker.Retrieve(c)
	cl := claims.Retrieve(c)

	// if we are not using a symmetric token, and the subject does not match the input, request should be denied
	if !strings.EqualFold(cl.TokenType, constants.ServerWorkerTokenType) && !strings.EqualFold(cl.Subject, w.GetHostname()) {
		retErr := fmt.Errorf("unable to refresh worker auth: claims subject %s does not match worker hostname %s", cl.Subject, w.GetHostname())

		logrus.WithFields(logrus.Fields{
			"subject": cl.Subject,
			"worker":  w.GetHostname(),
		}).Warnf("attempted refresh of worker %s using token from worker %s", w.GetHostname(), cl.Subject)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// set last checked in time
	w.SetLastCheckedIn(time.Now().Unix())

	// send API call to update the worker
	err := database.FromContext(c).UpdateWorker(w)
	if err != nil {
		retErr := fmt.Errorf("unable to update worker %s: %w", w.GetHostname(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"worker": w.GetHostname(),
	}).Infof("refreshing worker %s authentication", w.GetHostname())

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
