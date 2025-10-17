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
	enable := c.MustGet("storage-enable").(bool)
	if !enable {
		l := c.MustGet("logger").(*logrus.Entry)
		l.Info("storage is not enabled, sending storage disabled response")
		e := c.MustGet("storage-enable").(bool)
		wr := types.StorageInfo{
			StorageEnabled: &e,
		}

		c.JSON(http.StatusOK, wr)

		return
	}

	l := c.MustGet("logger").(*logrus.Entry)

	l.Info("requesting storage credentials with registration token")

	// extract the storage-enable that was packed into gin context
	e := c.MustGet("storage-enable").(bool)

	// extract the public key that was packed into gin context
	k := c.MustGet("storage-access-key").(string)

	// extract the storage-address that was packed into gin context
	a := c.MustGet("storage-address").(string)

	// extract the secret key that was packed into gin context
	s := c.MustGet("storage-secret-key").(string)

	// extract bucket name that was packed into gin context
	b := c.MustGet("storage-bucket").(string)

	wr := types.StorageInfo{
		StorageEnabled:   &e,
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

// ListObjectNames represents the API handler to list only the names of objects in a bucket.
func ListObjectNames(c *gin.Context) {
	enable := c.MustGet("storage-enable").(bool)
	if !enable {
		l := c.MustGet("logger").(*logrus.Entry)
		l.Info("storage is not enabled, skipping credentials request")
		c.JSON(http.StatusForbidden, gin.H{"error": "storage is not enabled"})

		return
	}

	l := c.MustGet("logger").(*logrus.Entry)
	l.Debug("listing object names in bucket")

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

	// Extract just the names from the objects
	names := make([]string, 0, len(objects))
	for _, obj := range objects {
		names = append(names, obj.Key)
	}

	c.JSON(http.StatusOK, gin.H{"names": names})
}

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
	org := util.PathParameter(c, "org")
	repo := util.PathParameter(c, "repo")
	buildNum := util.PathParameter(c, "build")

	// Validate parameters
	if org == "" || repo == "" || buildNum == "" {
		l.Error("missing required parameters (org, repo, or build)")
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing required parameters"})

		return
	}

	l.Debugf("listing object names in bucket %s for %s/%s build #%s", bucketName, org, repo, buildNum)

	// Create a new bucket object
	b := &types.Bucket{
		BucketName: bucketName,
		Recursive:  true,
	}

	// Call the ListBuildObjectNames method that handles prefix filtering
	objectNames, err := storage.FromGinContext(c).ListBuildObjectNames(
		c.Request.Context(),
		b,
		org,
		repo,
		buildNum,
	)
	if err != nil {
		l.Errorf("unable to list objects for %s/%s build #%s: %v", org, repo, buildNum, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, gin.H{"names": objectNames})
}
