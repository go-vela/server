// SPDX-License-Identifier: Apache-2.0

package artifact

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	artifact "github.com/go-vela/server/router/middleware/artifact"
	"github.com/go-vela/server/router/middleware/build"
)

// swagger:operation GET /api/v1/repos/{org}/{repo}/builds/{build}/artifact/{artifact} artifacts GetArtifact
//
// Get an artifact
//
// ---
// produces:
// - application/json
// parameters:
// - in: path
//   name: org
//   description: Name of the organization
//   required: true
//   type: string
// - in: path
//   name: repo
//   description: Name of the repository
//   required: true
//   type: string
// - in: path
//   name: build
//   description: Build number
//   required: true
//   type: integer
// - in: path
//   name: artifact
//   description: artifact ID
//   required: true
//   type: integer
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved the artifact
//     type: json
//     schema:
//       "$ref": "#/definitions/Artifact"
//   '400':
//     description: Invalid request payload or path
//     schema:
//       "$ref": "#/definitions/Error"
//   '401':
//     description: Unauthorized
//     schema:
//       "$ref": "#/definitions/Error"
//   '404':
//     description: Not found
//     schema:
//       "$ref": "#/definitions/Error"

// GetArtifact represents the API handler to get
// an artifact for a build.
func GetArtifact(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
	b := build.Retrieve(c)
	a := artifact.Retrieve(c)

	l.Debugf("getting artifact %d for build %d", a.GetID(), b.GetNumber())

	// return the artifact with presigned URL
	response := gin.H{
		"id":            a.GetID(),
		"build_id":      a.GetBuildID(),
		"file_name":     a.GetFileName(),
		"file_type":     a.GetFileType(),
		"file_size":     a.GetFileSize(),
		"object_path":   a.GetObjectPath(),
		"presigned_url": a.GetPresignedURL(),
		"created_at":    a.GetCreatedAt(),
	}

	c.JSON(http.StatusOK, response)
}
