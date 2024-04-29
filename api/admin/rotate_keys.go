// SPDX-License-Identifier: Apache-2.0

package admin

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/util"
)

// swagger:operation POST /api/v1/admin/rotate_oidc admin AdminRotateOIDCKeys
//
// Rotate RSA Keys
//
// ---
// produces:
// - application/json
// parameters:
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully updated the repo in the database
//     schema:
//       "$ref": "#/definitions/Repo"
//   '501':
//     description: Unable to update the repo in the database
//     schema:
//       "$ref": "#/definitions/Error"

// RotateOIDCKeys represents the API handler to
// rotate RSA keys in OIDC provider service.
func RotateOIDCKeys(c *gin.Context) {
	logrus.Info("Admin: rotating keys for OIDC provider")

	// capture middleware values
	ctx := c.Request.Context()

	err := database.FromContext(c).RotateKeys(ctx)
	if err != nil {
		retErr := fmt.Errorf("unable to rotate keys: %w", err)
		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, "keys rotated successfully")
}
