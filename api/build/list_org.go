// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package build

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/api"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// swagger:operation GET /api/v1/repos/{org} builds ListBuildsForOrg
//
// Get a list of builds by org in the configured backend
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
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved build list
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

// ListBuildsForOrg represents the API handler to capture a
// list of builds associated with an org from the configured backend.
func ListBuildsForOrg(c *gin.Context) {
	// variables that will hold the build list, build list filters and total count
	var (
		filters = map[string]interface{}{}
		b       []*library.Build
		t       int64
	)

	// capture middleware values
	o := org.Retrieve(c)
	u := user.Retrieve(c)

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"org":  o,
		"user": u.GetName(),
	}).Infof("listing builds for org %s", o)

	// capture the branch name parameter
	branch := c.Query("branch")
	// capture the event type parameter
	event := c.Query("event")
	// capture the status type parameter
	status := c.Query("status")

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
			event != constants.EventTag {
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

	// capture page query parameter if present
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		retErr := fmt.Errorf("unable to convert page query parameter for org %s: %w", o, err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// capture per_page query parameter if present
	perPage, err := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	if err != nil {
		retErr := fmt.Errorf("unable to convert per_page query parameter for Org %s: %w", o, err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// ensure per_page isn't above or below allowed values
	perPage = util.MaxInt(1, util.MinInt(100, perPage))

	// See if the user is an org admin to bypass individual permission checks
	perm, err := scm.FromContext(c).OrgAccess(u, o)
	if err != nil {
		logrus.Errorf("unable to get user %s access level for org %s", u.GetName(), o)
	}
	// Only show public repos to non-admins
	//nolint:goconst // ignore need for constant
	if perm != "admin" {
		filters["visibility"] = constants.VisibilityPublic
	}

	// send API call to capture the list of builds for the org (and event type if passed in)
	b, t, err = database.FromContext(c).ListBuildsForOrg(o, filters, page, perPage)

	if err != nil {
		retErr := fmt.Errorf("unable to list builds for org %s: %w", o, err)

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
