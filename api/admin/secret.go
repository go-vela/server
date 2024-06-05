// SPDX-License-Identifier: Apache-2.0

package admin

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/library"
)

// swagger:operation PUT /api/v1/admin/secret admin AdminUpdateSecret
//
// Update a secret in the database
//
// ---
// produces:
// - application/json
// parameters:
// - in: body
//   name: body
//   description: Payload containing secret to update
//   required: true
//   schema:
//     "$ref": "#/definitions/Secret"
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully updated the secret in the database
//     schema:
//       "$ref": "#/definitions/Secret"
//   '401':
//     description: Unauthorized to update the secret in the database
//     schema:
//       "$ref": "#/definitions/Error
//   '400':
//     description: Unable to update the secret in the database - bad request
//     schema:
//       "$ref": "#/definitions/Error"
//   '501':
//     description: Unable to update the secret in the database
//     schema:
//       "$ref": "#/definitions/Error"

// UpdateSecret represents the API handler to
// update any secret stored in the database.
func UpdateSecret(c *gin.Context) {
	// capture middleware values
	ctx := c.Request.Context()
	u := user.Retrieve(c)

	logger := logrus.WithFields(logrus.Fields{
		"ip":      util.EscapeValue(c.ClientIP()),
		"path":    util.EscapeValue(c.Request.URL.Path),
		"user":    u.GetName(),
		"user_id": u.GetID(),
	})

	logrus.Debug("platform admin: updating secret")

	// capture body from API request
	input := new(library.Secret)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for secret %d: %w", input.GetID(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	logger.WithFields(logrus.Fields{
		"secret_id": input.GetID(),
		"org":       util.EscapeValue(input.GetOrg()),
		"repo":      util.EscapeValue(input.GetRepo()),
		"type":      util.EscapeValue(input.GetType()),
		"name":      util.EscapeValue(input.GetName()),
		"team":      util.EscapeValue(input.GetTeam()),
	}).Debug("attempting to update secret")

	// send API call to update the secret
	s, err := database.FromContext(c).UpdateSecret(ctx, input)
	if err != nil {
		retErr := fmt.Errorf("unable to update secret %d: %w", input.GetID(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	logger.WithFields(logrus.Fields{
		"secret_id": s.GetID(),
		"org":       s.GetOrg(),
		"repo":      s.GetRepo(),
		"type":      s.GetType(),
		"name":      s.GetName(),
		"team":      s.GetTeam(),
	}).Info("secret updated")

	c.JSON(http.StatusOK, s)
}
