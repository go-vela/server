// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/api"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/util"
)

// swagger:operation GET /api/v1/repos/{org} repos ListReposForOrg
//
// Get all repositories for an organization
//
// ---
// produces:
// - application/json
// security:
//   - ApiKeyAuth: []
// parameters:
// - in: path
//   name: org
//   description: Name of the organization
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
//         description: See https://tools.ietf.org/html/rfc5988
//         type: string
//   '400':
//     description: Invalid request payload or path
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

// ListReposForOrg represents the API handler to get a list
// of repositories for an organization.
func ListReposForOrg(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
	o := org.Retrieve(c)
	u := user.Retrieve(c)
	ctx := c.Request.Context()

	l.Debugf("listing repos for org %s", o)

	pagination, err := api.ParsePagination(c)
	if err != nil {
		util.HandleError(c, http.StatusBadRequest, err)
		return
	}

	// capture the sort_by query parameter if present
	sortBy := util.QueryParameter(c, "sort_by", "name")

	// prep filters
	filters := make(map[string]interface{})

	// capture the query parameters if present:
	//
	// * active
	active := util.QueryParameter(c, "active", "true")
	// ensure active is a boolean and add it to filters as such
	if activeBool, err := strconv.ParseBool(active); err == nil {
		filters["active"] = activeBool
	}

	// See if the user is an org admin to bypass individual permission checks
	perm, err := scm.FromContext(c).OrgAccess(ctx, u, o)
	if err != nil {
		l.Errorf("unable to get user %s access level for org %s", u.GetName(), o)
	}
	// Only show public repos to non-admins
	if perm != constants.PermissionAdmin {
		filters["visibility"] = constants.VisibilityPublic
	}

	// send API call to capture the list of repos for the org
	r, err := database.FromContext(c).ListReposForOrg(ctx, o, sortBy, filters, pagination.Page, pagination.PerPage)
	if err != nil {
		retErr := fmt.Errorf("unable to get repos for org %s: %w", o, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// set pagination results
	pagination.Results = len(r)
	// set pagination headers
	pagination.SetHeaderLink(c)

	c.JSON(http.StatusOK, r)
}
