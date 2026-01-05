// SPDX-License-Identifier: Apache-2.0

package testattachment

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/testattachment"
)

// swagger:operation GET /api/v1/repos/{org}/{repo}/builds/{build}/reports/testattachment/{attachment} testreports GetTestAttachment
//
// Get a test attachment
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
//   name: attachment
//   description: Test attachment ID
//   required: true
//   type: integer
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved the test attachment
//     type: json
//     schema:
//       "$ref": "#/definitions/TestAttachment"
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

// GetTestAttachment represents the API handler to get
// a test attachment for a build.
func GetTestAttachment(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
	b := build.Retrieve(c)
	ta := testattachment.Retrieve(c)

	l.Debugf("getting test attachment %d for build %d", ta.GetID(), b.GetNumber())

	// return the test attachment with presigned URL
	response := gin.H{
		"id":             ta.GetID(),
		"test_report_id": ta.GetTestReportID(),
		"file_name":      ta.GetFileName(),
		"file_type":      ta.GetFileType(),
		"file_size":      ta.GetFileSize(),
		"object_path":    ta.GetObjectPath(),
		"presigned_url":  ta.GetPresignedURL(),
		"created_at":     ta.GetCreatedAt(),
	}

	c.JSON(http.StatusOK, response)
}
