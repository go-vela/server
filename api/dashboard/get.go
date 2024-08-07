// SPDX-License-Identifier: Apache-2.0

package dashboard

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/dashboard"
	"github.com/go-vela/server/util"
)

// swagger:operation GET /api/v1/dashboards/{dashboard} dashboards GetDashboard
//
// Get a dashboard
//
// ---
// produces:
// - application/json
// parameters:
// - in: path
//   name: dashboard
//   description: Dashboard id to retrieve
//   required: true
//   type: string
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved dashboard
//     type: json
//     schema:
//       "$ref": "#/definitions/Dashboard"
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

// GetDashboard represents the API handler to get
// a dashboard for a repository.
func GetDashboard(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
	d := dashboard.Retrieve(c)

	var err error

	l.Debugf("reading dashboard %s", d.GetID())

	// initialize DashCard and set dashboard to the dashboard info pulled from database
	dashboard := new(types.DashCard)
	dashboard.Dashboard = d

	// build RepoPartials referenced in the dashboard
	dashboard.Repos, err = buildRepoPartials(c, d.GetRepos())
	if err != nil {
		util.HandleError(c, http.StatusInternalServerError, err)

		return
	}

	c.JSON(http.StatusOK, dashboard)
}

// buildRepoPartials is a helper function which takes the dashboard repo list
// and builds a list of RepoPartials with information about the associated
// repository and its latest five builds.
func buildRepoPartials(c context.Context, repos []*types.DashboardRepo) ([]types.RepoPartial, error) {
	var result []types.RepoPartial

	for _, r := range repos {
		repo := types.RepoPartial{}

		// fetch repo from database
		dbRepo, err := database.FromContext(c).GetRepo(c, r.GetID())
		if err != nil {
			return nil, fmt.Errorf("unable to get repo %s for dashboard: %w", r.GetName(), err)
		}

		// set values for RepoPartial
		repo.Org = dbRepo.GetOrg()
		repo.Name = dbRepo.GetName()
		repo.Counter = dbRepo.GetCounter()
		repo.Active = dbRepo.GetActive()

		// list last 5 builds for repo given the branch and event filters
		builds, err := database.FromContext(c).ListBuildsForDashboardRepo(c, dbRepo, r.GetBranches(), r.GetEvents())
		if err != nil {
			return nil, fmt.Errorf("unable to list builds for repo %s in dashboard: %w", dbRepo.GetFullName(), err)
		}

		bPartials := []types.BuildPartial{}

		// populate BuildPartials with info from builds list
		for _, build := range builds {
			bPartial := types.BuildPartial{
				Number:   build.GetNumber(),
				Status:   build.GetStatus(),
				Started:  build.GetStarted(),
				Finished: build.GetFinished(),
				Ref:      build.GetRef(),
				Sender:   build.GetSender(),
				Branch:   build.GetBranch(),
				Event:    build.GetEvent(),
				Link:     build.GetLink(),
			}

			bPartials = append(bPartials, bPartial)
		}

		repo.Builds = bPartials

		result = append(result, repo)
	}

	return result, nil
}
