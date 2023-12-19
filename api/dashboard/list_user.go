// SPDX-License-Identifier: Apache-2.0

package dashboard

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
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

	// iterate through user dashboards and build a list of DashCards
	for _, dashboard := range u.GetDashboards() {
		dashCard := DashCard{}

		d, err := database.FromContext(c).GetDashboard(c, dashboard)
		if err != nil {
			// check if the query returned a record not found error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				d = new(library.Dashboard)
				d.SetID(dashboard)

				dashCard.Dashboard = d
				// if user dashboard has been deleted, append empty dashboard
				// to set and continue
				dashCards = append(dashCards, dashCard)

				continue
			}

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
