// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package repo

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
	"github.com/sirupsen/logrus"
)

// swagger:operation GET /api/v1/repos/{org} repos ListReposForOrg
//
// Get all repos for the provided org in the configured backend
//
// ---
// produces:
// - application/json
// security:
//   - ApiKeyAuth: []
// parameters:
// - in: path
//   name: org
//   description: Name of the org
//   required: true
//   type: string
// - in: query
//   name: active
//   description: Filter active repos
//   type: boolean
//   default: true
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
//   name: sort_by
//   description: How to sort the results
//   type: string
//   enum:
//   - name
//   - latest
//   default: name
// responses:
//   '200':
//     description: Successfully retrieved the repo
//     schema:
//       type: array
//       items:
//         "$ref": "#/definitions/Repo"
//     headers:
//       X-Total-Count:
//         description: Total number of results
//         type: integer
//       Link:
//         description: see https://tools.ietf.org/html/rfc5988
//         type: string
//   '400':
//     description: Unable to retrieve the org
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to retrieve the org
//     schema:
//       "$ref": "#/definitions/Error"

// ListReposForOrg represents the API handler to capture a list
// of repos for an org from the configured backend.
func ListReposForOrg(c *gin.Context) {
	// capture middleware values
	o := org.Retrieve(c)
	u := user.Retrieve(c)
	ctx := c.Request.Context()

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"org":  o,
		"user": u.GetName(),
	}).Infof("listing repos for org %s", o)

	// capture page query parameter if present
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		retErr := fmt.Errorf("unable to convert page query parameter for user %s: %w", u.GetName(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// capture per_page query parameter if present
	perPage, err := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	if err != nil {
		retErr := fmt.Errorf("unable to convert per_page query parameter for user %s: %w", u.GetName(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// ensure per_page isn't above or below allowed values
	perPage = util.MaxInt(1, util.MinInt(100, perPage))

	// capture the sort_by query parameter if present
	sortBy := util.QueryParameter(c, "sort_by", "name")

	// capture the query parameters if present:
	//
	// * active
	filters := map[string]interface{}{
		"active": util.QueryParameter(c, "active", "true"),
	}

	// See if the user is an org admin to bypass individual permission checks
	perm, err := scm.FromContext(c).OrgAccess(ctx, u, o)
	if err != nil {
		logrus.Errorf("unable to get user %s access level for org %s", u.GetName(), o)
	}
	// Only show public repos to non-admins
	if perm != "admin" {
		filters["visibility"] = constants.VisibilityPublic
	}

	// send API call to capture the list of repos for the org
	r, t, err := database.FromContext(c).ListReposForOrg(ctx, o, sortBy, filters, page, perPage)
	if err != nil {
		retErr := fmt.Errorf("unable to get repos for org %s: %w", o, err)

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

	c.JSON(http.StatusOK, r)
}
