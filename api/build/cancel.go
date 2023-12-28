// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/internal/token"
	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/executors"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// swagger:operation DELETE /api/v1/repos/{org}/{repo}/builds/{build}/cancel builds CancelBuild
//
// Cancel a running build
//
// ---
// produces:
// - application/json
// parameters:
// - in: path
//   name: repo
//   description: Name of the repo
//   required: true
//   type: string
// - in: path
//   name: org
//   description: Name of the org
//   required: true
//   type: string
// - in: path
//   name: build
//   description: Build number to cancel
//   required: true
//   type: integer
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully canceled the build
//     schema:
//       type: string
//   '400':
//     description: Unable to cancel build
//     schema:
//       "$ref": "#/definitions/Error"
//   '404':
//     description: Unable to cancel build
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to cancel build
//     schema:
//       "$ref": "#/definitions/Error"

// CancelBuild represents the API handler to cancel a running build.
//
//nolint:funlen // ignore statement count
func CancelBuild(c *gin.Context) {
	// capture middleware values
	b := build.Retrieve(c)
	e := executors.Retrieve(c)
	o := org.Retrieve(c)
	r := repo.Retrieve(c)
	user := user.Retrieve(c)
	ctx := c.Request.Context()

	entry := fmt.Sprintf("%s/%d", r.GetFullName(), b.GetNumber())

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"build": b.GetNumber(),
		"org":   o,
		"repo":  r.GetName(),
		"user":  user.GetName(),
	}).Infof("canceling build %s", entry)

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
			if strings.EqualFold(executor.Repo.GetFullName(), r.GetFullName()) && *executor.GetBuild().Number == b.GetNumber() {
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

	// remove build executable for clean up
	_, err = database.FromContext(c).PopBuildExecutable(ctx, b.GetID())
	if err != nil {
		retErr := fmt.Errorf("unable to pop build %s from executables table: %w", entry, err)
		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// retrieve the steps for the build from the step table
	steps := []*library.Step{}
	page := 1
	perPage := 100

	for page > 0 {
		// retrieve build steps (per page) from the database
		stepsPart, _, err := database.FromContext(c).ListStepsForBuild(b, map[string]interface{}{}, page, perPage)
		if err != nil {
			retErr := fmt.Errorf("unable to retrieve steps for build %s: %w", entry, err)
			util.HandleError(c, http.StatusNotFound, retErr)

			return
		}

		// add page of steps to list steps
		steps = append(steps, stepsPart...)

		// assume no more pages exist if under 100 results are returned
		if len(stepsPart) < 100 {
			page = 0
		} else {
			page++
		}
	}

	// iterate over each step for the build
	// setting anything running or pending to canceled
	for _, step := range steps {
		if step.GetStatus() == constants.StatusRunning || step.GetStatus() == constants.StatusPending {
			step.SetStatus(constants.StatusCanceled)

			_, err = database.FromContext(c).UpdateStep(step)
			if err != nil {
				retErr := fmt.Errorf("unable to update step %s for build %s: %w", step.GetName(), entry, err)
				util.HandleError(c, http.StatusNotFound, retErr)

				return
			}
		}
	}

	// retrieve the services for the build from the service table
	services := []*library.Service{}
	page = 1

	for page > 0 {
		// retrieve build services (per page) from the database
		servicesPart, _, err := database.FromContext(c).ListServicesForBuild(ctx, b, map[string]interface{}{}, page, perPage)
		if err != nil {
			retErr := fmt.Errorf("unable to retrieve services for build %s: %w", entry, err)
			util.HandleError(c, http.StatusNotFound, retErr)

			return
		}

		// add page of services to the list of services
		services = append(services, servicesPart...)

		// assume no more pages exist if under 100 results are returned
		if len(servicesPart) < 100 {
			page = 0
		} else {
			page++
		}
	}

	// iterate over each service for the build
	// setting anything running or pending to canceled
	for _, service := range services {
		if service.GetStatus() == constants.StatusRunning || service.GetStatus() == constants.StatusPending {
			service.SetStatus(constants.StatusCanceled)

			_, err = database.FromContext(c).UpdateService(ctx, service)
			if err != nil {
				retErr := fmt.Errorf("unable to update service %s for build %s: %w",
					service.GetName(),
					entry,
					err,
				)
				util.HandleError(c, http.StatusNotFound, retErr)

				return
			}
		}
	}

	c.JSON(http.StatusOK, b)
}
