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

// swagger:operation PUT /api/v1/dashboards/{dashboard} dashboards UpdateDashboard
//
// Update a dashboard for the configured backend
//
// ---
// produces:
// - application/json
// parameters:
// - in: path
//   name: dashboard
//   description: ID of the dashboard
//   required: true
//   type: int
// - in: body
//   name: body
//   description: Payload containing the dashboard to update
//   required: true
//   schema:
//     "$ref": "#/definitions/Dashboard"
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully updated the dashboard
//     schema:
//       "$ref": "#/definitions/Dashboard"
//   '400':
//     description: Unable to update the dashboard
//     schema:
//       "$ref": "#/definitions/Error"
//   '404':
//     description: Unable to update the dashboard
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to update the dashboard
//     schema:
//       "$ref": "#/definitions/Error"

// UpdateDashboard represents the API handler to update
// a dashboard in the configured backend.
func UpdateDashboard(c *gin.Context) {
	// capture middleware values
	d := dashboard.Retrieve(c)
	ctx := c.Request.Context()
	u := user.Retrieve(c)

	admin := false

	for _, a := range d.GetAdmins() {
		if u.GetName() == a {
			admin = true
			break
		}
	}

	if !admin {
		retErr := fmt.Errorf("unable to update dashboard %d: user is not an admin", d.GetID())

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"dashboard": d.GetID(),
	}).Infof("updating dashboard %d", d.GetID())

	// capture body from API request
	input := new(library.Dashboard)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for dashboard %d: %w", d.GetID(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	if input.GetName() != "" {
		// update name if defined
		d.SetName(input.GetName())
	}

	if len(input.GetAdmins()) > 0 {
		admins, err := validateAdminSet(c, u, input.GetAdmins())
		if err != nil {
			util.HandleError(c, http.StatusBadRequest, err)

			return
		}

		d.SetAdmins(admins)
	}

	// set the updated by field using claims
	d.SetUpdatedBy(u.GetName())

	if len(input.GetRepos()) > 0 {
		// validate supplied repo list
		err = validateRepoSet(c, input.GetRepos())
		if err != nil {
			util.HandleError(c, http.StatusBadRequest, err)

			return
		}

		d.SetRepos(input.GetRepos())
	}

	// update the dashboard within the database
	d, err = database.FromContext(c).UpdateDashboard(ctx, d)
	if err != nil {
		retErr := fmt.Errorf("unable to update dashboard %d: %w", input.GetID(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, d)
}
