// SPDX-License-Identifier: Apache-2.0

package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/util"
)

// swagger:operation GET /badge/{org}/{repo}/status.svg base GetBadge
//
// Get a build status badge for a repository
//
// ---
// produces:
// - image/svg+xml
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
// - in: query
//   name: branch
//   description: Name of the branch
//   type: string
// responses:
//   '200':
//     description: Successfully retrieved the build status badge
//     schema:
//       type: string
//   '400':
//     description: Invalid request payload or path
//     schema:
//       "$ref": "#/definitions/Error"
//   '404':
//     description: Not found
//     schema:
//       "$ref": "#/definitions/Error"

// GetBadge represents the API handler to
// return a build status badge.
func GetBadge(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
	r := repo.Retrieve(c)
	ctx := c.Request.Context()

	branch := util.QueryParameter(c, "branch", r.GetBranch())

	l.Debugf("creating latest build badge for repo %s on branch %s", r.GetFullName(), branch)

	// send API call to capture the last build for the repo and branch
	b, err := database.FromContext(c).LastBuildForRepo(ctx, r, branch)
	if err != nil {
		c.String(http.StatusOK, constants.BadgeUnknown)
		return
	}

	badge := badgeForStatus(b.GetStatus())

	// set headers to prevent caching
	c.Header("Content-Type", "image/svg+xml")
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Expires", "0") // passing invalid date sets resource as expired

	c.String(http.StatusOK, badge)
}

// badgeForStatus is a helper to match the build status with a badge.
func badgeForStatus(s string) string {
	switch s {
	case constants.StatusRunning, constants.StatusPending:
		return constants.BadgeRunning
	case constants.StatusFailure, constants.StatusKilled:
		return constants.BadgeFailed
	case constants.StatusSuccess:
		return constants.BadgeSuccess
	case constants.StatusError:
		return constants.BadgeError
	default:
		return constants.BadgeUnknown
	}
}
