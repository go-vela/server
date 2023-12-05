// SPDX-License-Identifier: Apache-2.0

package dashboard

import (
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

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"dashboard": d.GetID(),
		"user":      u.GetName(),
	}).Infof("reading dashboard %d", d.GetID())

	dashboard := new(DashCard)
	dashboard.Dashboard = d

	var repos []RepoPartial

	for _, r := range d.Repos {
		repo := RepoPartial{}

		dbRepo, err := database.FromContext(c).GetRepo(c, r.GetID())
		if err != nil {
			retErr := fmt.Errorf("unable to get repo for dashboard %d: %w", d.GetID(), err)

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}

		repo.Org = dbRepo.GetOrg()
		repo.Name = dbRepo.GetName()
		repo.Counter = dbRepo.GetCounter()

		builds, err := database.FromContext(c).ListBuildsForDashboardRepo(c, dbRepo, r.GetBranches(), r.GetEvents())
		if err != nil {
			retErr := fmt.Errorf("unable to list builds for repo %s in dashboard %d: %w", dbRepo.GetFullName(), d.GetID(), err)

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}

		bPartials := []BuildPartial{}

		for _, build := range builds {
			bPartial := BuildPartial{}
			bPartial.Number = build.GetNumber()
			bPartial.Status = build.GetStatus()
			bPartial.Started = build.GetStarted()
			bPartial.Finished = build.GetFinished()
			bPartial.Sender = build.GetSender()

			bPartials = append(bPartials, bPartial)
		}

		repo.Builds = bPartials

		repos = append(repos, repo)
	}

	dashboard.Repos = repos

	c.JSON(http.StatusOK, dashboard)
}
