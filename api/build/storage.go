package build

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/storage"
	"github.com/go-vela/server/util"
	"github.com/sirupsen/logrus"
	"net/http"
)

// swagger:operation POST /api/v1/repos/{org}/{repo}/builds/{build}/storage/upload builds UploadObject
//
// # Upload an object to a bucket
//
// ---
// produces:
// - application/json
// parameters:
//   - in: body
//     name: body
//     description: The object to be uploaded
//     required: true
//     schema:
//     type: object
//     properties:
//     bucketName:
//     type: string
//     objectName:
//     type: string
//     objectData:
//     type: string
//
// security:
//   - ApiKeyAuth: []
//
// responses:
//
//	'201':
//	  description: Successfully uploaded the object
//	'400':
//	  description: Invalid request payload
//	  schema:
//	    "$ref": "#/definitions/Error"
//	'500':
//	  description: Unexpected server error
//	  schema:
//	    "$ref": "#/definitions/Error"
//
// UploadObject represents the API handler to upload an object to a bucket.
func UploadObject(c *gin.Context) {
	l := c.MustGet("logger").(*logrus.Entry)
	ctx := c.Request.Context()

	l.Debug("platform admin: uploading object")

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
	err = storage.FromGinContext(c).Upload(ctx, input)
	if err != nil {
		retErr := fmt.Errorf("unable to upload object: %w", err)
		util.HandleError(c, http.StatusInternalServerError, retErr)
		return
	}

	c.Status(http.StatusCreated)
}
