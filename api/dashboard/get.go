// SPDX-License-Identifier: Apache-2.0

package dashboard

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/dashboard"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

type RepoPartial struct {
	Org     string         `json:"org,omitempty"`
	Name    string         `json:"name,omitempty"`
	Counter int            `json:"counter,omitempty"`
	Builds  []BuildPartial `json:"builds,omitempty"`
}

type BuildPartial struct {
	Number   int    `json:"number,omitempty"`
	Started  int64  `json:"started,omitempty"`
	Finished int64  `json:"finished,omitempty"`
	Sender   string `json:"sender,omitempty"`
	Status   string `json:"status,omitempty"`
	Event    string `json:"event,omitempty"`
	Branch   string `json:"branch,omitempty"`
	Link     string `json:"link,omitempty"`
}

type DashCard struct {
	Dashboard *library.Dashboard `json:"dashboard,omitempty"`
	Repos     []RepoPartial      `json:"repos,omitempty"`
}

// swagger:operation GET /api/v1/dashboards/{dashboard} dashboards GetDashboard
//
// Get a dashboard in the configured backend
//
// ---
// produces:
// - application/json
// parameters:
// - in: path
//   name: dashboard
//   description: Dashboard id to retrieve
//   required: true
//   type: integer
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved the dashboard
//     type: json
//     schema:
//       "$ref": "#/definitions/Dashboard"

// GetDashboard represents the API handler to capture
// a dashboard for a repo from the configured backend.
func GetDashboard(c *gin.Context) {
	// capture middleware values
	d := dashboard.Retrieve(c)
	u := user.Retrieve(c)

	var err error

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"dashboard": d.GetID(),
		"user":      u.GetName(),
	}).Infof("reading dashboard %s", d.GetID())

	dashboard := new(DashCard)
	dashboard.Dashboard = d

	dashboard.Repos, err = buildRepoPartials(c, d.Repos)
	if err != nil {
		util.HandleError(c, http.StatusInternalServerError, err)

		return
	}

	c.JSON(http.StatusOK, dashboard)
}

func buildRepoPartials(c context.Context, repos []*library.DashboardRepo) ([]RepoPartial, error) {
	var result []RepoPartial

	for _, r := range repos {
		repo := RepoPartial{}

		dbRepo, err := database.FromContext(c).GetRepo(c, r.GetID())
		if err != nil {
			return nil, fmt.Errorf("unable to get repo %s for dashboard: %w", r.GetName(), err)
		}

		repo.Org = dbRepo.GetOrg()
		repo.Name = dbRepo.GetName()
		repo.Counter = dbRepo.GetCounter()

		builds, err := database.FromContext(c).ListBuildsForDashboardRepo(c, dbRepo, r.GetBranches(), r.GetEvents())
		if err != nil {
			return nil, fmt.Errorf("unable to list builds for repo %s in dashboard: %w", dbRepo.GetFullName(), err)
		}

		bPartials := []BuildPartial{}

		for _, build := range builds {
			bPartial := BuildPartial{
				Number:   build.GetNumber(),
				Status:   build.GetStatus(),
				Started:  build.GetStarted(),
				Finished: build.GetFinished(),
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
