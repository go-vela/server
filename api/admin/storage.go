// SPDX-License-Identifier: Apache-2.0

package admin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/storage"
	"github.com/go-vela/server/util"
	"github.com/sirupsen/logrus"
	"net/http"
)

// swagger:operation POST /api/v1/admin/storage/bucket admin CreateBucket

//
// Create a new bucket
//
// ---
// produces:
// - application/json
// parameters:
// - in: body
//   name: body
//   description: The bucket name to be created
//   required: true
//   schema:
//     type: object
//     properties:
//       bucketName:
//         type: string
// security:
//   - ApiKeyAuth: []
// responses:
//   '201':
//     description: Successfully created the bucket
//   '400':
//     description: Invalid request payload
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unexpected server error
//     schema:
//       "$ref": "#/definitions/Error"

// CreateBucket represents the API handler to create a new bucket.
func CreateBucket(c *gin.Context) {
	l := c.MustGet("logger").(*logrus.Entry)
	ctx := c.Request.Context()

	l.Debug("platform admin: creating bucket")

	// capture body from API request
	input := new(types.Bucket)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for bucket %s: %w", input.BucketName, err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	err = storage.FromGinContext(c).CreateBucket(ctx, input)
	if err != nil {
		retErr := fmt.Errorf("unable to create bucket: %w", err)
		util.HandleError(c, http.StatusInternalServerError, retErr)
		return
	}

	c.Status(http.StatusCreated)
}

// swagger:operation DELETE /api/v1/admin/storage/bucket admin DeleteBucket
//
// Delete a bucket
//
// ---
// produces:
// - application/json
// parameters:
// - in: body
//   name: body
//   description: The bucket name to be deleted
//   required: true
//   schema:
//     type: object
//     properties:
//       bucketName:
//         type: string
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully deleted the bucket
//   '400':
//     description: Invalid request payload
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unexpected server error
//     schema:
//       "$ref": "#/definitions/Error"

// DeleteBucket represents the API handler to delete a bucket.
func DeleteBucket(c *gin.Context) {
	l := c.MustGet("logger").(*logrus.Entry)
	ctx := c.Request.Context()

	l.Debug("platform admin: deleting bucket")

	// capture body from API request
	input := new(types.Bucket)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for bucket %s: %w", input.BucketName, err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	err = storage.FromGinContext(c).DeleteBucket(ctx, input)
	if err != nil {
		retErr := fmt.Errorf("unable to delete bucket: %w", err)
		util.HandleError(c, http.StatusInternalServerError, retErr)
		return
	}

	c.Status(http.StatusOK)
}

// swagger:operation PUT /api/v1/admin/storage/bucket/lifecycle admin AdminSetBucketLifecycle
//
// Set bucket lifecycle configuration
//
// ---
// produces:
// - application/json
// parameters:
// - in: body
//   name: body
//   description: The bucket lifecycle configuration
//   required: true
//   schema:
//     type: object
//     properties:
//       bucketName:
//         type: string
//       lifecycleConfig:
//         type: string
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully set the bucket lifecycle configuration
//   '400':
//     description: Invalid request payload
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unexpected server error
//     schema:
//       "$ref": "#/definitions/Error"

// SetBucketLifecycle represents the API handler to set bucket lifecycle configuration.
func SetBucketLifecycle(c *gin.Context) {
	l := c.MustGet("logger").(*logrus.Entry)
	ctx := c.Request.Context()

	l.Debug("platform admin: setting bucket lifecycle configuration")

	// capture body from API request
	input := new(types.Bucket)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for bucket %s: %w", input.BucketName, err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	err = storage.FromGinContext(c).SetBucketLifecycle(ctx, input)
	if err != nil {
		retErr := fmt.Errorf("unable to set bucket lifecycle configuration: %w", err)
		util.HandleError(c, http.StatusInternalServerError, retErr)
		return
	}

	c.Status(http.StatusOK)
}
