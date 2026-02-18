// SPDX-License-Identifier: Apache-2.0

package storage

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/storage"
)

// swagger:operation GET /api/v1/repos/{org}/{repo}/builds/{build}/storage/sts storage GetSTSCreds
//
// Get temporary STS credentials for build storage uploads.
//
// Generates temporary AWS STS credentials scoped to allow PUT operations
// into the configured storage bucket under the build-specific prefix.
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
// security:
//   - ApiKeyAuth: []
// responses:
//   200:
//     description: Successfully generated temporary STS credentials
//     schema:
//       $ref: '#/definitions/STSCreds'
//   400:
//     description: Bad request due to invalid parameters
//     schema:
//       $ref: '#/definitions/Error'
//   403:
//     description: Storage is not enabled or invalid token
//     schema:
//       $ref: '#/definitions/Error'
//   404:
//     description: Repo not found
//     schema:
//       $ref: '#/definitions/Error'
//   500:
//     description: Unable to assume role or generate credentials
//     schema:
//       $ref: '#/definitions/Error'

// GetSTSCreds represents the API handler to generate temporary STS credentials for build storage uploads.
func GetSTSCreds(c *gin.Context) {
	l := c.MustGet("logger").(*logrus.Entry)

	enabled := c.MustGet("storage-enable").(bool)
	if !enabled {
		l.Info("storage is not enabled, skipping credentials request")
		c.JSON(http.StatusForbidden, gin.H{"error": "storage is not enabled"})

		return
	}

	r := repo.Retrieve(c)
	org := r.GetOrg()
	b := build.Retrieve(c)
	repoName := r.GetName()
	buildNum := b.GetNumber()
	ctx := c.Request.Context()

	prefix := fmt.Sprintf("%s/%s/%d/", org, repoName, buildNum)

	sessionName := fmt.Sprintf("vela-%s-%s-%d", org, repoName, buildNum)

	creds, err := storage.FromGinContext(c).AssumeRole(ctx, int(r.GetTimeout())*60, prefix, sessionName)
	if creds == nil {
		l.Errorf("unable to assume role and generate temporary credentials without error %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to assume role and generate temporary credentials"})

		return
	}

	if err != nil {
		l.Errorf("unable to assume role: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, creds)
}
