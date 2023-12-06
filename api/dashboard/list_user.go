// SPDX-License-Identifier: Apache-2.0

package dashboard

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
	"github.com/sirupsen/logrus"
)

// swagger:operation GET /api/v1/user/dashboards dashboards ListUserDashboards
//
// Get all dashboards for the claims users in the configured backend
//
// ---
// produces:
// - application/json
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved the dashboards
//     schema:
//       type: array
//       items:
//         "$ref": "#/definitions/Dashboard"
//   '400':
//     description: Unable to retrieve the org
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to retrieve the org
//     schema:
//       "$ref": "#/definitions/Error"

// ListUserDashboards represents the API handler to capture a list
// of dashboards for a user from the configured backend.
func ListUserDashboards(c *gin.Context) {
	// capture middleware values
	u := user.Retrieve(c)

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"user": u.GetName(),
	}).Infof("listing dashboards for user %s", u.GetName())

	var dashCards []DashCard

	for _, dashboard := range u.GetDashboards() {
		dashCard := DashCard{}

		d, err := database.FromContext(c).GetDashboard(c, dashboard)
		if err != nil {
			retErr := fmt.Errorf("unable to get dashboard %s: %w", dashboard, err)

			util.HandleError(c, http.StatusBadRequest, retErr)

			return
		}

		dashCard.Dashboard = d

		dashCard.Repos, err = buildRepoPartials(c, d.Repos)
		if err != nil {
			util.HandleError(c, http.StatusInternalServerError, err)

			return
		}

		dashCards = append(dashCards, dashCard)
	}

	c.JSON(http.StatusOK, dashCards)
}
