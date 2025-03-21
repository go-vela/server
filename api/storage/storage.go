// SPDX-License-Identifier: Apache-2.0

package storage

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/storage"
	"github.com/go-vela/server/util"
)

// swagger:operation POST /api/v1/storage/info storage StorageInfo
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

// swagger:operation GET /api/v1/storage/{bucket}/objects storage ListObjects
//
// List objects in a bucket
//
// ---
// produces:
// - application/json
// parameters:
// - in: path
//   name: bucket
//   description: Name of the bucket
//   required: true
//   type: string
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully listed objects in the bucket
//     schema:
//       type: array
//       items:
//         type: string
//   '500':
//     description: Unexpected server error
//     schema:
//       "$ref": "#/definitions/Error"

// ListObjects represents the API handler to list objects in a bucket.
func ListObjects(c *gin.Context) {
	l := c.MustGet("logger").(*logrus.Entry)

	l.Debug("listing objects in bucket")

	// extract the bucket name from the request
	bucketName := util.PathParameter(c, "bucket")

	// create a new bucket object
	b := &types.Bucket{
		BucketName: bucketName,
	}

	// list objects in the bucket
	objects, err := storage.FromGinContext(c).ListObjects(c.Request.Context(), b)
	if err != nil {
		l.Errorf("unable to list objects in bucket %s: %v", bucketName, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, gin.H{"objects": objects})
}
