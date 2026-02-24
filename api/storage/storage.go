// SPDX-License-Identifier: Apache-2.0

package storage

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/storage"
)

// swagger:operation GET /api/v1/repos/{org}/{repo}/builds/{build}/storage storage ListBuildObjectNames
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
// security:
//   - ApiKeyAuth: []
// responses:
//   200:
//     description: Successfully listed object names for the build
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

	r := repo.Retrieve(c)
	b := build.Retrieve(c)
	org := r.GetOrg()
	buildNum := b.GetNumber()

	l.Debugf("listing object names in bucket for %s/%s build #%d", org, r.GetName(), buildNum)

	// Call the ListBuildObjectNames method that handles prefix filtering
	objectNames, err := storage.FromGinContext(c).ListBuildObjectNames(
		c.Request.Context(),
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
