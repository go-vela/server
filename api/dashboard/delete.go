// SPDX-License-Identifier: Apache-2.0

package dashboard

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/dashboard"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
)

// swagger:operation DELETE /api/v1/dashboards/{dashboard} dashboards DeleteDashboard
//
// Delete a dashboard
//
// ---
// produces:
// - application/json
// parameters:
// - in: path
//   name: dashboard
//   description: Dashboard ID
//   required: true
//   type: string
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully deleted dashboard
//     schema:
//       type: string
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

// DeleteDashboard represents the API handler to remove a dashboard.
func DeleteDashboard(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
	d := dashboard.Retrieve(c)
	u := user.Retrieve(c)
	ctx := c.Request.Context()

	l.Debugf("deleting dashboard %s", d.GetID())

	if !isAdmin(d, u) {
		retErr := fmt.Errorf("unable to delete dashboard %s: user is not an admin", d.GetID())

		util.HandleError(c, http.StatusUnauthorized, retErr)

		return
	}

	// Remove the dashboard ID from the user's dashboards
	dashboards := u.GetDashboards()

	updatedDashboards := []string{}
	for _, id := range dashboards {
		if id != d.GetID() {
			updatedDashboards = append(updatedDashboards, id)
		}
	}

	u.SetDashboards(updatedDashboards)

	// send API call to update the user
	u, err := database.FromContext(c).UpdateUser(ctx, u)
	if err != nil {
		retErr := fmt.Errorf("unable to update user %s: %w", u.GetName(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	err = database.FromContext(c).DeleteDashboard(c, d)
	if err != nil {
		retErr := fmt.Errorf("error while deleting dashboard %s: %w", d.GetID(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("dashboard %s deleted", d.GetName()))
}

// isAdmin is a helper function that iterates through the dashboard admins
// and confirms if the user is in the slice.
func isAdmin(d *types.Dashboard, u *types.User) bool {
	for _, admin := range d.GetAdmins() {
		if admin.GetID() == u.GetID() {
			return true
		}
	}

	return false
}
