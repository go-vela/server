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
	"strings"
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
	l.Debugf("bucket name: %s", input.BucketName)
	err = storage.FromGinContext(c).CreateBucket(ctx, input)
	if err != nil {
		retErr := fmt.Errorf("unable to create bucket: %w", err)
		util.HandleError(c, http.StatusInternalServerError, retErr)
		return
	}

	c.Status(http.StatusCreated)
}

// swagger:operation GET /api/v1/admin/storage/bucket/download admin DownloadObject
//
// # Download an object from a bucket
//
// ---
// produces:
// - application/json
// parameters:
//   - in: query
//     name: bucketName
//     description: The name of the bucket
//     required: true
//     type: string
//   - in: query
//     name: objectName
//     description: The name of the object
//     required: true
//     type: string
//
// security:
//   - ApiKeyAuth: []
//
// responses:
//
//	'200':
//	  description: Successfully downloaded the object
//	'400':
//	  description: Invalid request payload
//	  schema:
//	    "$ref": "#/definitions/Error"
//	'500':
//	  description: Unexpected server error
//	  schema:
//	    "$ref": "#/definitions/Error"
//
// DownloadObject represents the API handler to download an object from a bucket.
func DownloadObject(c *gin.Context) {
	l := c.MustGet("logger").(*logrus.Entry)
	ctx := c.Request.Context()

	l.Debug("platform admin: downloading object")

	// capture body from API request
	input := new(types.Object)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for object %s: %w", input.ObjectName, err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}
	if input.Bucket.BucketName == "" || input.ObjectName == "" {
		retErr := fmt.Errorf("bucketName and objectName are required")
		util.HandleError(c, http.StatusBadRequest, retErr)
		return
	}
	if input.FilePath == "" {
		retErr := fmt.Errorf("file path is required")
		util.HandleError(c, http.StatusBadRequest, retErr)
		return
	}
	if strings.ContainsAny(input.FilePath, "/\\") || strings.Contains(input.FilePath, "..") || strings.TrimSpace(input.FilePath) == "" {
		retErr := fmt.Errorf("invalid file path")
		util.HandleError(c, http.StatusBadRequest, retErr)
		return
	}
	err = storage.FromGinContext(c).Download(ctx, input)
	if err != nil {
		retErr := fmt.Errorf("unable to download object: %w", err)
		util.HandleError(c, http.StatusInternalServerError, retErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("File has been downloaded to %s", input.FilePath)})
}

// swagger:operation GET /api/v1/admin/storage/presign admin GetPresignedURL
//
// # Generate a presigned URL for an object
//
// ---
// produces:
// - application/json
// parameters:
//   - in: query
//     name: bucketName
//     description: The name of the bucket
//     required: true
//     type: string
//   - in: query
//     name: objectName
//     description: The name of the object
//     required: true
//     type: string
//
// security:
//   - ApiKeyAuth: []
//
// responses:
//
//	'200':
//	  description: Successfully generated the presigned URL
//	'400':
//	  description: Invalid request payload
//	  schema:
//	    "$ref": "#/definitions/Error"
//	'500':
//	  description: Unexpected server error
//	  schema:
//	    "$ref": "#/definitions/Error"
func GetPresignedURL(c *gin.Context) {
	l := c.MustGet("logger").(*logrus.Entry)
	ctx := c.Request.Context()

	l.Debug("platform admin: generating presigned URL")

	// capture query parameters from API request
	bucketName := c.Query("bucketName")
	objectName := c.Query("objectName")

	if bucketName == "" || objectName == "" {
		retErr := fmt.Errorf("bucketName and objectName are required")
		util.HandleError(c, http.StatusBadRequest, retErr)
		return
	}

	input := &types.Object{
		Bucket:     types.Bucket{BucketName: bucketName},
		ObjectName: objectName,
	}

	url, err := storage.FromGinContext(c).PresignedGetObject(ctx, input)
	if err != nil || url == "" {
		retErr := fmt.Errorf("unable to generate presigned URL: %w", err)
		util.HandleError(c, http.StatusBadRequest, retErr)
		return
	}

	c.JSON(http.StatusOK, url)
}
