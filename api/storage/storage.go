// SPDX-License-Identifier: Apache-2.0

package storage

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/api/types"
)

// swagger:operation POST /api/v1/storage/info storage Info
//
// Get storage credentials
//
// ---
// produces:
// - application/json
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved storage credentials
//     schema:
//       "$ref": "#/definitions/StorageInfo"
//   '401':
//     description: Unauthorized
//     schema:
//       "$ref": "#/definitions/Error"

// Info represents the API handler to
// retrieve storage credentials as part of worker onboarding.
func Info(c *gin.Context) {
	l := c.MustGet("logger").(*logrus.Entry)

	l.Info("requesting storage credentials with registration token")

	// extract the public key that was packed into gin context
	k := c.MustGet("access-key").(string)

	// extract the storage-address that was packed into gin context
	a := c.MustGet("storage-address").(string)

	// extract the secret key that was packed into gin context
	s := c.MustGet("secret-key").(string)

	// extract bucket name that was packed into gin context
	b := c.MustGet("storage-bucket").(string)

	wr := types.StorageInfo{
		StorageAccessKey: &k,
		StorageAddress:   &a,
		StorageSecretKey: &s,
		StorageBucket:    &b,
	}

	c.JSON(http.StatusOK, wr)
}
