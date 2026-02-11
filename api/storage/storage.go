// SPDX-License-Identifier: Apache-2.0

package storage

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/storage"
	"github.com/go-vela/server/util"
)

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
	enable := c.MustGet("storage-enable").(bool)
	if !enable {
		l := c.MustGet("logger").(*logrus.Entry)
		l.Info("storage is not enabled, skipping credentials request")
		c.JSON(http.StatusForbidden, gin.H{"error": "storage is not enabled"})

		return
	}

	l := c.MustGet("logger").(*logrus.Entry)

	l.Debug("listing objects in bucket")

	// extract the bucket name from the request
	bucketName := util.PathParameter(c, "bucket")

	// create a new bucket object
	b := &types.Bucket{
		BucketName: bucketName,
		Recursive:  true,
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

// swagger:operation GET /api/v1/repos/{org}/{repo}/builds/{build}/storage/{bucket}/names storage ListBuildObjectNames
//
// List object names for a specific build in a bucket.
//
// ---
// produces:
// - application/json
// parameters:
//   - name: org
//     in: path
//     description: Organization name
//     required: true
//     type: string
//   - name: repo
//     in: path
//     description: Repository name
//     required: true
//     type: string
//   - name: build
//     in: path
//     description: Build number
//     required: true
//     type: integer
//     format: int64
//   - name: bucket
//     in: path
//     description: Name of the bucket
//     required: true
//     type: string
// security:
//   - ApiKeyAuth: []
// responses:
//   200:
//     description: Successfully listed object names for the build
//   500:
//     description: Unexpected server error
//     schema:
//       $ref: '#/definitions/Error'

// ListBuildObjectNames represents the API handler to list object names for a specific build.
func ListBuildObjectNames(c *gin.Context) {
	enable := c.MustGet("storage-enable").(bool)
	if !enable {
		l := c.MustGet("logger").(*logrus.Entry)
		l.Info("storage is not enabled, skipping credentials request")
		c.JSON(http.StatusForbidden, gin.H{"error": "storage is not enabled"})

		return
	}

	l := c.MustGet("logger").(*logrus.Entry)

	// Extract path parameters
	bucketName := util.PathParameter(c, "bucket")
	r := repo.Retrieve(c)
	b := build.Retrieve(c)
	org := r.GetOrg()
	buildNum := b.GetNumber()

	l.Debugf("listing object names in bucket %s for %s/%s build #%d", bucketName, org, r.GetName(), buildNum)

	// Create a new bucket object
	bObject := &types.Bucket{
		BucketName: bucketName,
		Recursive:  true,
	}

	// Call the ListBuildObjectNames method that handles prefix filtering
	objectNames, err := storage.FromGinContext(c).ListBuildObjectNames(
		c.Request.Context(),
		bObject,
		org,
		r.GetName(),
		strconv.FormatInt(buildNum, 10),
	)
	if err != nil {
		l.Errorf("unable to list objects for %s/%s build #%d: %v", org, r.GetName(), buildNum, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, gin.H{"names": objectNames})
}
