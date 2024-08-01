// SPDX-License-Identifier: Apache-2.0

package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// swagger:operation GET /health base Health
//
// Check the Vela API health
//
// ---
// produces:
// - application/json
// parameters:
// responses:
//   '200':
//     description: Successfully 'ping'-ed Vela API
//     schema:
//       type: string

// Health represents the API handler to
// report the health status for Vela.
func Health(c *gin.Context) {
	c.JSON(http.StatusOK, "ok")
}
