// SPDX-License-Identifier: Apache-2.0

package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// swagger:operation PUT /api/v1/admin/deployment admin AdminUpdateDeployment
//
// Get All (Not Implemented)
//
// ---
// produces:
// - application/json
// parameters:
// responses:
//   '401':
//     description: Unauthorized to update deployment
//     schema:
//       "$ref": "#/definitions/Error
//   '501':
//     description: This endpoint is not implemented
//     schema:
//       type: string

// UpdateDeployment represents the API handler to
// update any deployment stored in the database.
func UpdateDeployment(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "The server does not support the functionality required to fulfill the request.")
}
