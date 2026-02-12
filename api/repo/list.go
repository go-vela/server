// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/api"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
)

// swagger:operation GET /api/v1/repos repos ListRepos
//
// Get all repositories
//
// ---
// produces:
// - application/json
// security:
//   - ApiKeyAuth: []
// parameters:
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

// ListRepos represents the API handler to get a list
// of repositories for a user.
func ListRepos(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
	u := user.Retrieve(c)
	ctx := c.Request.Context()

	l.Debugf("listing repos for user %s", u.GetName())

	pagination, err := api.ParsePagination(c)
	if err != nil {
		util.HandleError(c, http.StatusBadRequest, err)
		return
	}

	// capture the sort_by query parameter if present
	sortBy := util.QueryParameter(c, "sort_by", "name")

	// capture the query parameters if present:
	//
	// * active
	filters := map[string]interface{}{
		"active": util.QueryParameter(c, "active", "true"),
	}

	// send API call to capture the list of repos for the user
	r, err := database.FromContext(c).ListReposForUser(ctx, u, sortBy, filters, pagination.Page, pagination.PerPage)
	if err != nil {
		retErr := fmt.Errorf("unable to get repos for user %s: %w", u.GetName(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// set pagination results
	pagination.Results = len(r)
	pagination.SetHeaderLink(c)

	c.JSON(http.StatusOK, r)
}
