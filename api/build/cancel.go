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

	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/cache"
	"github.com/go-vela/server/cache/models"
	"github.com/go-vela/server/compiler/types/yaml"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/internal/token"
	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/util"
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
func CancelBuild(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
	b := build.Retrieve(c)
	r := repo.Retrieve(c)
	user := user.Retrieve(c)
	ctx := c.Request.Context()

	entry := fmt.Sprintf("%s/%d", r.GetFullName(), b.GetNumber())

	l.Debugf("canceling build %s", entry)

	switch b.GetStatus() {
	case constants.StatusRunning:
		build, err := CancelRunning(c, b)
		if err != nil {
			retErr := fmt.Errorf("unable to cancel running build %s: %w", entry, err)
			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}

		if build == nil {
			c.JSON(http.StatusOK, "no running build found to cancel")

			return
		}

		build.SetError(fmt.Sprintf("build was canceled by %s", user.GetName()))

		build, err = database.FromContext(c).UpdateBuild(ctx, build)
		if err != nil {
			retErr := fmt.Errorf("unable to update status for build %s: %w", entry, err)
			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}

		l.WithFields(logrus.Fields{
			"build":    build.GetNumber(),
			"build_id": build.GetID(),
		}).Info("build updated - build canceled")

		c.JSON(http.StatusOK, build)

		return
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

	var (
		checks   []models.CheckRun
		scmToken string
	)

	if b.GetRepo().GetInstallID() != 0 {
		checks, err = cache.FromContext(c).GetCheckRuns(ctx, b)
		if err != nil {
			l.Errorf("unable to retrieve check runs for build %s: %v", entry, err)
		}

		scmToken, _, err = scm.FromContext(c).GetNetrcPassword(ctx, database.FromContext(c), cache.FromContext(c), b, yaml.Git{})
		if err != nil {
			l.Errorf("unable to generate new installation token for build %s: %v", entry, err)

			return
		}
	} else {
		scmToken = b.GetRepo().GetOwner().GetToken()
	}

	// send API call to set the status on the commit
	_, err = scm.FromContext(c).Status(ctx, b, scmToken, checks)
	if err != nil {
		l.Errorf("unable to set commit status for build %s: %v", entry, err)
	}

	// update component statuses to canceled
	err = UpdateComponentStatuses(c, b, constants.StatusCanceled, scmToken)
	if err != nil {
		retErr := fmt.Errorf("unable to update component statuses for build %s: %w", entry, err)
		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, b)
}

// CancelRunning is a helper function that determines the executor currently running a build and sends an API call
// to that executor's worker to cancel the build.
func CancelRunning(c *gin.Context, b *types.Build) (*types.Build, error) {
	l := c.MustGet("logger").(*logrus.Entry)

	// retrieve the worker
	w, err := database.FromContext(c).GetWorkerForHostname(c, b.GetHost())
	if err != nil {
		return nil, err
	}

	// retrieve the executors from the worker
	e, err := getWorkerExecutors(c, w)
	if err != nil {
		return nil, err
	}

	for _, executor := range *e {
		// check each executor on the worker running the build to see if it's running the build we want to cancel
		if executor.Build.GetID() == b.GetID() {
			build, err := cancelBuildForExecutor(c, l, w, &executor)
			if err != nil {
				return nil, err
			}

			return build, nil
		}
	}

	return nil, nil
}

// getWorkerExecutors is a helper function that retrieves the list of executors from a worker.
func getWorkerExecutors(c *gin.Context, w *types.Worker) (*[]types.Executor, error) {
	e := new([]types.Executor)

	// prepare the request to the worker to retrieve executors
	client := http.DefaultClient
	client.Timeout = 30 * time.Second
	endpoint := fmt.Sprintf("%s/api/v1/executors", w.GetAddress())

	req, err := createWorkerRequest(c, "GET", endpoint)
	if err != nil {
		return nil, err
	}

	// make the request to the worker and check the response
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// Read Response Body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// parse response and validate at least one item was returned
	err = json.Unmarshal(respBody, e)
	if err != nil {
		return nil, err
	}

	return e, nil
}

// cancelBuildForExecutor is a helper function that sends an API call to a specific executor to cancel the build.
func cancelBuildForExecutor(c *gin.Context, l *logrus.Entry, w *types.Worker, executor *types.Executor) (*types.Build, error) {
	b := new(types.Build)

	client := http.DefaultClient
	client.Timeout = 30 * time.Second

	// set the API endpoint path we send the request to
	u := fmt.Sprintf("%s/api/v1/executors/%d/build/cancel", w.GetAddress(), executor.GetID())

	req, err := createWorkerRequest(c, "DELETE", u)
	if err != nil {
		return nil, err
	}

	// perform the request to the worker
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	l.Debugf("sent cancel request to worker %s (executor %d) for build %d", w.GetHostname(), executor.GetID(), b.GetID())

	// Read Response Body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(respBody, b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// createWorkerRequest is a helper function that creates an authenticated HTTP request to a worker from the server.
func createWorkerRequest(c *gin.Context, method, endpoint string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(context.Background(), method, endpoint, nil)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	// add the token to authenticate to the worker
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tkn))

	return req, nil
}
