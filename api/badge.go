// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"net/http"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/types/constants"
	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

// swagger:operation GET /badge/{org}/{repo}/status.svg base GetBadge
//
// Get a badge for the repo
//
// ---
// produces:
// - image/svg+xml
// parameters:
// - in: path
//   name: org
//   description: Name of the org the repo belongs to
//   required: true
//   type: string
// - in: path
//   name: repo
//   description: Name of the repo to get the badge for
//   required: true
//   type: string
// responses:
//   '200':
//     description: Successfully retrieved a status Badge
//     schema:
//       type: string

// GetBadge represents the API handler to
// return a build status badge.
func GetBadge(c *gin.Context) {
	// capture middleware values
	r := repo.Retrieve(c)
	branch := c.DefaultQuery("branch", r.GetBranch())

	logrus.Infof("Creating badge for latest build on %s for branch %s", r.GetFullName(), branch)

	// send API call to capture the last build for the repo and branch
	b, err := database.FromContext(c).GetLastBuildByBranch(r, branch)
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
