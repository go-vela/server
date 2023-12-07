// SPDX-License-Identifier: Apache-2.0

package dashboard

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// swagger:operation POST /api/v1/dashboards dashboards CreateDashboard
//
// Create a dashboard in the configured backend
//
// ---
// produces:
// - application/json
// parameters:
// - in: body
//   name: body
//   description: Payload containing the dashboard to create
//   required: true
//   schema:
//     "$ref": "#/definitions/Dashboard"
// security:
//   - ApiKeyAuth: []
// responses:
//   '201':
//     description: Successfully created the dashboard
//     schema:
//       "$ref": "#/definitions/Dashboard"
//   '400':
//     description: Unable to create the dashboard
//     schema:
//       "$ref": "#/definitions/Error"
//   '403':
//     description: Unable to create the dashboard
//     schema:
//       "$ref": "#/definitions/Error"
//   '409':
//     description: Unable to create the dashboard
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to create the dashboard
//     schema:
//       "$ref": "#/definitions/Error"
//   '503':
//     description: Unable to create the dashboard
//     schema:
//       "$ref": "#/definitions/Error"

// CreateDashboard represents the API handler to
// create a dashboard in the configured backend.
func CreateDashboard(c *gin.Context) {
	// capture middleware values
	u := user.Retrieve(c)

	// capture body from API request
	input := new(library.Dashboard)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for new dashboard: %w", err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// ensure dashboard name is defined
	if input.GetName() == "" {
		util.HandleError(c, http.StatusBadRequest, fmt.Errorf("dashboard name must be set"))

		return
	}

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"user": u.GetName(),
	}).Infof("creating new dashboard %s", input.GetName())

	d := new(library.Dashboard)

	// update fields in dashboard object
	d.SetCreatedBy(u.GetName())
	d.SetName(input.GetName())
	d.SetCreatedAt(time.Now().UTC().Unix())
	d.SetUpdatedAt(time.Now().UTC().Unix())
	d.SetUpdatedBy(u.GetName())

	// validate admins to ensure they are all active users
	admins, err := validateAdminSet(c, u, input.GetAdmins())
	if err != nil {
		util.HandleError(c, http.StatusBadRequest, err)

		return
	}

	d.SetAdmins(admins)

	// validate repos to ensure they are all enabled
	err = validateRepoSet(c, input.GetRepos())
	if err != nil {
		util.HandleError(c, http.StatusBadRequest, err)

		return
	}

	d.SetRepos(input.GetRepos())

	// send API call to create the dashboard
	d, err = database.FromContext(c).CreateDashboard(c, d)
	if err != nil {
		retErr := fmt.Errorf("unable to create new dashboard %s: %w", d.GetName(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// add dashboard to claims user's dashboard set
	u.SetDashboards(append(u.GetDashboards(), d.GetID()))

	_, err = database.FromContext(c).UpdateUser(c, u)
	if err != nil {
		retErr := fmt.Errorf("unable to update user %s: %w", u.GetName(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusCreated, d)
}

func validateAdminSet(c context.Context, caller *library.User, users []string) ([]string, error) {
	// add user creating the dashboard to admin list
	admins := []string{caller.GetName()}

	// validate supplied admins are actual users
	for _, admin := range users {
		if admin == caller.GetName() {
			continue
		}

		u, err := database.FromContext(c).GetUserForName(c, admin)
		if err != nil || !u.GetActive() {
			return nil, fmt.Errorf("unable to create dashboard: %s is not an active user", admin)
		}

		admins = append(admins, admin)
	}

	return admins, nil
}

func validateRepoSet(c context.Context, repos []*library.DashboardRepo) error {
	for _, repo := range repos {
		// verify format (org/repo)
		parts := strings.Split(repo.GetName(), "/")
		if len(parts) != 2 {
			return fmt.Errorf("unable to create dashboard: %s is not a valid repo", repo.GetName())
		}

		// fetch repo from database
		dbRepo, err := database.FromContext(c).GetRepoForOrg(c, parts[0], parts[1])
		if err != nil || !dbRepo.GetActive() {
			return fmt.Errorf("unable to create dashboard: could not get repo %s: %w", repo.GetName(), err)
		}

		// override ID field if provided to match the database ID
		repo.SetID(dbRepo.GetID())
	}

	return nil
}
