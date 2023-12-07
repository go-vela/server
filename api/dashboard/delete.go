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
	"github.com/sirupsen/logrus"
)

// swagger:operation DELETE /api/v1/dashboards/{dashboard} dashboards DeleteDashboard
//
// Delete a dashboard in the configured backend
//
// ---
// produces:
// - application/json
// parameters:
// - in: path
//   name: dashboard
//   description: id of the dashboard
//   required: true
//   type: string
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully deleted the dashboard
//     schema:
//       type: string
//   '500':
//     description: Unable to  deleted the dashboard
//     schema:
//       "$ref": "#/definitions/Error"
//   '510':
//     description: Unable to  deleted the dashboard
//     schema:
//       "$ref": "#/definitions/Error"

// DeleteDashboard represents the API handler to remove
// a dashboard from the configured backend.
func DeleteDashboard(c *gin.Context) {
	// capture middleware values
	d := dashboard.Retrieve(c)
	u := user.Retrieve(c)

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"dashboard": d.GetID(),
		"user":      u.GetName(),
	}).Infof("deleting dashboard %s", d.GetID())

	if !isAdmin(u.GetName(), d.GetAdmins()) {
		retErr := fmt.Errorf("unable to delete dashboard %s: user is not an admin", d.GetName())

		util.HandleError(c, http.StatusUnauthorized, retErr)

		return
	}

	// Comment out actual delete until delete mechanism is fleshed out
	err := database.FromContext(c).DeleteDashboard(c, d)
	if err != nil {
		retErr := fmt.Errorf("error while deleting dashboard %s: %w", d.GetID(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("dashboard %s deleted", d.GetName()))
}
