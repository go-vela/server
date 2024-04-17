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

// swagger:operation PUT /api/v1/dashboards/{dashboard} dashboards UpdateDashboard
//
// Update a dashboard for the configured backend
//
// ---
// produces:
// - application/json
// parameters:
// - name: dashboard
//   in: path
//   description: ID of the dashboard
//   required: true
//   type: string
// - name: body
//   in: body
//   description: Payload containing the dashboard to update
//   required: true
//   schema:
//     $ref: '#/definitions/Dashboard'
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully updated dashboard
//     schema:
//       "$ref": "#/definitions/Dashboard"
//   '400':
//     description: Bad request when updating dashboard
//     schema:
//       "$ref": "#/definitions/Error"
//   '401':
//     description: Unauthorized to update dashboard
//     schema:
//       "$ref": "#/definitions/Error"
//   '404':
//     description: Unable to find dashboard
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Error while updating dashboard
//     schema:
//       "$ref": "#/definitions/Error"

// UpdateDashboard represents the API handler to update
// a dashboard in the configured backend.
func UpdateDashboard(c *gin.Context) {
	// capture middleware values
	d := dashboard.Retrieve(c)
	u := user.Retrieve(c)

	if !isAdmin(d, u) {
		retErr := fmt.Errorf("unable to update dashboard %s: user is not an admin", d.GetID())

		util.HandleError(c, http.StatusUnauthorized, retErr)

		return
	}

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"dashboard": d.GetID(),
	}).Infof("updating dashboard %s", d.GetID())

	// capture body from API request
	input := new(types.Dashboard)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for dashboard %s: %w", d.GetID(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	if input.GetName() != "" {
		// update name if defined
		d.SetName(input.GetName())
	}

	// validate admin set if supplied
	if len(input.GetAdmins()) > 0 {
		admins, err := createAdminSet(c, u, input.GetAdmins())
		if err != nil {
			util.HandleError(c, http.StatusBadRequest, err)

			return
		}

		d.SetAdmins(admins)
	}

	// set the updated by field using claims
	d.SetUpdatedBy(u.GetName())

	// validate repo set if supplied
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
	d, err = database.FromContext(c).UpdateDashboard(c, d)
	if err != nil {
		retErr := fmt.Errorf("unable to update dashboard %s: %w", input.GetID(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, d)
}
