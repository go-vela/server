// SPDX-License-Identifier: Apache-2.0

package admin

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/library"
)

// swagger:operation PUT /api/v1/admin/secret admin AdminUpdateSecret
//
// Update a secret
//
// ---
// produces:
// - application/json
// parameters:
// - in: body
//   name: body
//   description: The secret object with the fields to be updated
//   required: true
//   schema:
//     "$ref": "#/definitions/Secret"
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully updated the secret
//     schema:
//       "$ref": "#/definitions/Secret"
//   '401':
//     description: Unauthorized
//     schema:
//       "$ref": "#/definitions/Error"
//   '400':
//     description: Invalid request payload
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unexpected server error
//     schema:
//       "$ref": "#/definitions/Error"

// UpdateSecret represents the API handler to update a secret.
func UpdateSecret(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
	ctx := c.Request.Context()

	l.Debug("platform admin: updating secret")

	// capture body from API request
	input := new(library.Secret)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for secret %d: %w", input.GetID(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	l.WithFields(logrus.Fields{
		"secret_id":   input.GetID(),
		"secret_org":  util.EscapeValue(input.GetOrg()),
		"secret_repo": util.EscapeValue(input.GetRepo()),
		"secret_type": util.EscapeValue(input.GetType()),
		"secret_name": util.EscapeValue(input.GetName()),
		"secret_team": util.EscapeValue(input.GetTeam()),
	}).Debug("platform admin: attempting to update secret")

	// send API call to update the secret
	s, err := database.FromContext(c).UpdateSecret(ctx, input)
	if err != nil {
		retErr := fmt.Errorf("unable to update secret %d: %w", input.GetID(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	l.WithFields(logrus.Fields{
		"secret_id":   s.GetID(),
		"secret_org":  s.GetOrg(),
		"secret_repo": s.GetRepo(),
		"secret_type": s.GetType(),
		"secret_name": s.GetName(),
		"secret_team": s.GetTeam(),
	}).Info("platform admin: secret updated")

	c.JSON(http.StatusOK, s)
}
