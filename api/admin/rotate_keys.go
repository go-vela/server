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

// swagger:operation POST /api/v1/admin/rotate_oidc_keys admin AdminRotateOIDCKeys
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
//     description: Successfully rotated OIDC provider keys
//     schema:
//       type: string
//   '401':
//     description: Unauthorized
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unexpected server error
//     schema:
//       "$ref": "#/definitions/Error"

// RotateOIDCKeys represents the API handler to
// rotate RSA keys in the OIDC provider service.
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
