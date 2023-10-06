// SPDX-License-Identifier: Apache-2.0

package build

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/api"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// swagger:operation GET /api/v1/repos/{org}/{repo}/builds builds ListBuildsForRepo
//
// Get builds from the configured backend
//
// ---
// produces:
// - application/json
// parameters:
// - in: path
//   name: org
//   description: Name of the org
//   required: true
//   type: string
// - in: path
//   name: repo
//   description: Name of the repo
//   required: true
//   type: string
// - in: query
//   name: event
//   description: Filter by build event
//   type: string
//   enum:
//   - comment
//   - deployment
//   - pull_request
//   - push
//   - schedule
//   - tag
// - in: query
//   name: commit
//   description: Filter builds based on the commit hash
//   type: string
// - in: query
//   name: branch
//   description: Filter builds by branch
//   type: string
// - in: query
//   name: status
//   description: Filter by build status
//   type: string
//   enum:
//   - canceled
//   - error
//   - failure
//   - killed
//   - pending
//   - running
//   - success
// - in: query
//   name: page
//   description: The page of results to retrieve
//   type: integer
//   default: 1
// - in: query
//   name: per_page
//   description: How many results per page to return
//   type: integer
//   maximum: 100
//   default: 10
// - in: query
//   name: before
//   description: filter builds created before a certain time
//   type: integer
//   default: 1
// - in: query
//   name: after
//   description: filter builds created after a certain time
//   type: integer
//   default: 0
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved the builds
//     schema:
//       type: array
//       items:
//         "$ref": "#/definitions/Build"
//     headers:
//       X-Total-Count:
//         description: Total number of results
//         type: integer
//       Link:
//         description: see https://tools.ietf.org/html/rfc5988
//         type: string
//   '400':
//     description: Unable to retrieve the list of builds
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to retrieve the list of builds
//     schema:
//       "$ref": "#/definitions/Error"

// ListBuildsForRepo represents the API handler to capture a
// list of builds for a repo from the configured backend.
func ListBuildsForRepo(c *gin.Context) {
	// variables that will hold the build list, build list filters and total count
	var (
		filters = map[string]interface{}{}
		b       []*library.Build
		t       int64
	)

	// capture middleware values
	o := org.Retrieve(c)
	r := repo.Retrieve(c)
	u := user.Retrieve(c)
	ctx := c.Request.Context()

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"org":  o,
		"repo": r.GetName(),
		"user": u.GetName(),
	}).Infof("listing builds for repo %s", r.GetFullName())

	// capture the branch name parameter
	branch := c.Query("branch")
	// capture the event type parameter
	event := c.Query("event")
	// capture the status type parameter
	status := c.Query("status")
	// capture the commit hash parameter
	commit := c.Query("commit")

	// check if branch filter was provided
	if len(branch) > 0 {
		// add branch to filters map
		filters["branch"] = branch
	}
	// check if event filter was provided
	if len(event) > 0 {
		// verify the event provided is a valid event type
		if event != constants.EventComment && event != constants.EventDeploy &&
			event != constants.EventPush && event != constants.EventPull &&
			event != constants.EventTag && event != constants.EventSchedule {
			retErr := fmt.Errorf("unable to process event %s: invalid event type provided", event)

			util.HandleError(c, http.StatusBadRequest, retErr)

			return
		}

		// add event to filters map
		filters["event"] = event
	}
	// check if status filter was provided
	if len(status) > 0 {
		// verify the status provided is a valid status type
		if status != constants.StatusCanceled && status != constants.StatusError &&
			status != constants.StatusFailure && status != constants.StatusKilled &&
			status != constants.StatusPending && status != constants.StatusRunning &&
			status != constants.StatusSuccess {
			retErr := fmt.Errorf("unable to process status %s: invalid status type provided", status)

			util.HandleError(c, http.StatusBadRequest, retErr)

			return
		}

		// add status to filters map
		filters["status"] = status
	}

	// check if commit hash filter was provided
	if len(commit) > 0 {
		// add commit to filters map
		filters["commit"] = commit
	}

	// capture page query parameter if present
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		retErr := fmt.Errorf("unable to convert page query parameter for repo %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// capture per_page query parameter if present
	perPage, err := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	if err != nil {
		retErr := fmt.Errorf("unable to convert per_page query parameter for repo %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// ensure per_page isn't above or below allowed values
	perPage = util.MaxInt(1, util.MinInt(100, perPage))

	// capture before query parameter if present, default to now
	before, err := strconv.ParseInt(c.DefaultQuery("before", strconv.FormatInt(time.Now().UTC().Unix(), 10)), 10, 64)
	if err != nil {
		retErr := fmt.Errorf("unable to convert before query parameter for repo %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// capture after query parameter if present, default to 0
	after, err := strconv.ParseInt(c.DefaultQuery("after", "0"), 10, 64)
	if err != nil {
		retErr := fmt.Errorf("unable to convert after query parameter for repo %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	b, t, err = database.FromContext(c).ListBuildsForRepo(ctx, r, filters, before, after, page, perPage)
	if err != nil {
		retErr := fmt.Errorf("unable to list builds for repo %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// create pagination object
	pagination := api.Pagination{
		Page:    page,
		PerPage: perPage,
		Total:   t,
	}
	// set pagination headers
	pagination.SetHeaderLink(c)

	c.JSON(http.StatusOK, b)
}
