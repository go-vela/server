// SPDX-License-Identifier: Apache-2.0

package dashboard

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
)

// swagger:operation POST /api/v1/dashboards dashboards CreateDashboard
//
// Create a dashboard
//
// ---
// produces:
// - application/json
// parameters:
// - in: body
//   name: body
//   description: Dashboard object to create
//   required: true
//   schema:
//     "$ref": "#/definitions/Dashboard"
// security:
//   - ApiKeyAuth: []
// responses:
//   '201':
//     description: Successfully created dashboard
//     schema:
//       "$ref": "#/definitions/Dashboard"
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

// CreateDashboard represents the API handler to
// create a dashboard.
func CreateDashboard(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
	u := user.Retrieve(c)

	// capture body from API request
	input := new(types.Dashboard)

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

	l.Debugf("creating new dashboard %s", input.GetName())

	d := new(types.Dashboard)

	// update fields in dashboard object
	d.SetCreatedBy(u.GetName())
	d.SetName(input.GetName())
	d.SetCreatedAt(time.Now().UTC().Unix())
	d.SetUpdatedAt(time.Now().UTC().Unix())
	d.SetUpdatedBy(u.GetName())

	// validate admins to ensure they are all active users
	admins, err := createAdminSet(c, u, input.GetAdmins())
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

	// create dashboard in database
	d, err = database.FromContext(c).CreateDashboard(c, d)
	if err != nil {
		retErr := fmt.Errorf("unable to create new dashboard %s: %w", d.GetName(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	l.WithFields(logrus.Fields{
		"dashboard":    d.GetName(),
		"dashboard_id": d.GetID(),
	}).Info("dashboard created")

	// add dashboard to claims' user's dashboard set
	u.SetDashboards(append(u.GetDashboards(), d.GetID()))

	// update user in database
	_, err = database.FromContext(c).UpdateUser(c, u)
	if err != nil {
		retErr := fmt.Errorf("unable to update user %s: %w", u.GetName(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	l.Infof("user updated with new dashboard %s", d.GetName())

	c.JSON(http.StatusCreated, d)
}

// createAdminSet takes a slice of users, cleanses it of duplicates and throws an error
// when a user is inactive or not found in the database. It returns a sanitized slice of admins.
func createAdminSet(c context.Context, caller *types.User, users []*types.User) ([]*types.User, error) {
	// add user creating the dashboard to admin list
	admins := []*types.User{caller.Crop()}

	dupMap := make(map[string]bool)

	// validate supplied admins are actual users
	for _, u := range users {
		if u.GetName() == caller.GetName() || dupMap[u.GetName()] {
			continue
		}

		dbUser, err := database.FromContext(c).GetUserForName(c, u.GetName())
		if err != nil || !dbUser.GetActive() {
			return nil, fmt.Errorf("unable to create dashboard: %s is not an active user", u.GetName())
		}

		admins = append(admins, dbUser.Crop())

		dupMap[dbUser.GetName()] = true
	}

	return admins, nil
}

// validateRepoSet is a helper function that confirms all dashboard repos exist and are enabled
// in the database while also confirming the IDs match when saving.
func validateRepoSet(c context.Context, repos []*types.DashboardRepo) error {
	for _, repo := range repos {
		// fetch repo from database
		dbRepo, err := database.FromContext(c).GetRepoForOrg(c, repo.GetName())
		if err != nil || !dbRepo.GetActive() {
			return fmt.Errorf("unable to create dashboard: could not get repo %s: %w", repo.GetName(), err)
		}

		// override ID field if provided to match the database ID
		repo.SetID(dbRepo.GetID())
	}

	return nil
}
