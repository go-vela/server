// SPDX-License-Identifier: Apache-2.0

package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/go-vela/server/version"
)

// swagger:operation GET /version base Version
//
// Get the version of the Vela API
//
// ---
// produces:
// - application/json
// parameters:
// responses:
//   '200':
//     description: Successfully retrieved the Vela API version
//     schema:
//       "$ref": "#/definitions/Version"

// Version represents the API handler to
// report the version number for Vela.
func Version(c *gin.Context) {
	c.JSON(http.StatusOK, version.New())
}
