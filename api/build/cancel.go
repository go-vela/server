// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/internal/token"
	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/executors"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/constants"
)

// swagger:operation DELETE /api/v1/repos/{org}/{repo}/builds/{build}/cancel builds CancelBuild
//
// Cancel a build
//
// ---
// produces:
// - application/json
// parameters:
// - in: path
//   name: org
//   description: Name of the organization
//   required: true
//   type: string
// - in: path
//   name: repo
//   description: Name of the repository
//   required: true
//   type: string
// - in: path
//   name: build
//   description: Build number
//   required: true
//   type: integer
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully canceled the build
//     schema:
//       "$ref": "#/definitions/Build"
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

// CancelBuild represents the API handler to cancel a build.
//
//nolint:funlen // ignore statement count
func CancelBuild(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
	b := build.Retrieve(c)
	e := executors.Retrieve(c)
	r := repo.Retrieve(c)
	user := user.Retrieve(c)
	ctx := c.Request.Context()

	entry := fmt.Sprintf("%s/%d", r.GetFullName(), b.GetNumber())

	l.Debugf("canceling build %s", entry)

	switch b.GetStatus() {
	case constants.StatusRunning:
		// retrieve the worker info
		w, err := database.FromContext(c).GetWorkerForHostname(ctx, b.GetHost())
		if err != nil {
			retErr := fmt.Errorf("unable to get worker for build %s: %w", entry, err)
			util.HandleError(c, http.StatusNotFound, retErr)

			return
		}

		for _, executor := range e {
			// check each executor on the worker running the build to see if it's running the build we want to cancel
			if executor.Build.GetID() == b.GetID() {
				// prepare the request to the worker
				client := http.DefaultClient
				client.Timeout = 30 * time.Second

				// set the API endpoint path we send the request to
				u := fmt.Sprintf("%s/api/v1/executors/%d/build/cancel", w.GetAddress(), executor.GetID())

				req, err := http.NewRequestWithContext(context.Background(), "DELETE", u, nil)
				if err != nil {
					retErr := fmt.Errorf("unable to form a request to %s: %w", u, err)
					util.HandleError(c, http.StatusBadRequest, retErr)

					return
				}

				tm := c.MustGet("token-manager").(*token.Manager)

				// set mint token options
				mto := &token.MintTokenOpts{
					Hostname:      "vela-server",
					TokenType:     constants.WorkerAuthTokenType,
					TokenDuration: time.Minute * 1,
				}

				// mint token
				tkn, err := tm.MintToken(mto)
				if err != nil {
					retErr := fmt.Errorf("unable to generate auth token: %w", err)
					util.HandleError(c, http.StatusInternalServerError, retErr)

					return
				}

				// add the token to authenticate to the worker
				req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tkn))

				// perform the request to the worker
				resp, err := client.Do(req)
				if err != nil {
					retErr := fmt.Errorf("unable to connect to %s: %w", u, err)
					util.HandleError(c, http.StatusBadRequest, retErr)

					return
				}
				defer resp.Body.Close()

				// Read Response Body
				respBody, err := io.ReadAll(resp.Body)
				if err != nil {
					retErr := fmt.Errorf("unable to read response from %s: %w", u, err)
					util.HandleError(c, http.StatusBadRequest, retErr)

					return
				}

				err = json.Unmarshal(respBody, b)
				if err != nil {
					retErr := fmt.Errorf("unable to parse response from %s: %w", u, err)
					util.HandleError(c, http.StatusBadRequest, retErr)

					return
				}

				b.SetError(fmt.Sprintf("build was canceled by %s", user.GetName()))

				b, err = database.FromContext(c).UpdateBuild(ctx, b)
				if err != nil {
					retErr := fmt.Errorf("unable to update status for build %s: %w", entry, err)
					util.HandleError(c, http.StatusInternalServerError, retErr)

					return
				}

				l.WithFields(logrus.Fields{
					"build":    b.GetNumber(),
					"build_id": b.GetID(),
				}).Info("build updated - build canceled")

				c.JSON(resp.StatusCode, b)

				return
			}
		}
	case constants.StatusPending, constants.StatusPendingApproval:
		break

	default:
		retErr := fmt.Errorf("found build %s but its status was %s", entry, b.GetStatus())

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// build has been abandoned
	// update the status in the build table
	b.SetStatus(constants.StatusCanceled)
	b.SetError(fmt.Sprintf("build was canceled by %s", user.GetName()))

	b, err := database.FromContext(c).UpdateBuild(ctx, b)
	if err != nil {
		retErr := fmt.Errorf("unable to update status for build %s: %w", entry, err)
		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	l.WithFields(logrus.Fields{
		"build":    b.GetNumber(),
		"build_id": b.GetID(),
	}).Info("build updated - build canceled")

	// remove build executable for clean up
	_, err = database.FromContext(c).PopBuildExecutable(ctx, b.GetID())
	if err != nil {
		retErr := fmt.Errorf("unable to pop build %s from executables table: %w", entry, err)
		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// update component statuses to canceled
	err = UpdateComponentStatuses(c, b, constants.StatusCanceled)
	if err != nil {
		retErr := fmt.Errorf("unable to update component statuses for build %s: %w", entry, err)
		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, b)
}
