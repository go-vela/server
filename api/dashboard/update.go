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
// Update a dashboard
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
//   description: The dashboard object with the fields to be updated
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

// UpdateDashboard represents the API handler to update a dashboard.
func UpdateDashboard(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
	d := dashboard.Retrieve(c)
	u := user.Retrieve(c)

	l.Debugf("updating dashboard %s", d.GetID())

	if !isAdmin(d, u) {
		retErr := fmt.Errorf("unable to update dashboard %s: user is not an admin", d.GetID())

		util.HandleError(c, http.StatusUnauthorized, retErr)

		return
	}

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
