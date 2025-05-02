// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/api"
	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/team"
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
// - in: path
//   name: team
//   description: Slug of team
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

// ListReposForTeam represents the API handler to get a list
// of repositories for a team and whether they are connected to Vela.
func ListReposForTeam(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
	o := org.Retrieve(c)
	t := team.Retrieve(c)
	u := user.Retrieve(c)
	ctx := c.Request.Context()

	l.Debugf("listing repos for team %s", o)

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
	perPage = max(1, min(100, perPage))

	repoList, code, err := scm.FromContext(c).ListTeamRepositories(ctx, u.GetToken(), o, t, page, perPage)
	if err != nil {
		retErr := fmt.Errorf("unable to get SCM repos for team %s: %w", t, err)
		util.HandleError(c, code, retErr)

		return
	}

	dbRepoList, err := database.FromContext(c).ReposInList(ctx, repoList)
	if err != nil {
		retErr := fmt.Errorf("unable to get database repos for team %s: %w", t, err)
		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	result := []*types.Repo{}

	for _, r := range repoList {
		var vela bool

		for _, dbRepo := range dbRepoList {
			if dbRepo.GetFullName() == r {
				result = append(result, dbRepo)

				vela = true

				break
			}
		}

		if !vela {
			result = append(result, &types.Repo{
				FullName: &r,
				Active:   &vela,
			})
		}
	}

	// create pagination object
	pagination := api.Pagination{
		Page:    page,
		PerPage: perPage,
		Results: len(result),
	}
	// set pagination headers
	pagination.SetHeaderLink(c)

	c.JSON(http.StatusOK, result)
}
