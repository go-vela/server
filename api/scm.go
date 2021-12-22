// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"fmt"
	"net/http"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/util"

	"github.com/go-vela/types/library"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// swagger:operation GET /api/v1/scm/orgs/{org}/sync scm SyncRepos
//
// Sync up repos from scm service and database in a specified org
//
// ---
// produces:
// - application/json
// parameters:
// - in: path
//   name: org
//   description: Name of the org
//   required: true
//   type: string
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully synchronized repos
//     schema:
//       type: string
//   '500':
//     description: Unable to synchronize org repositories
//     schema:
//       "$ref": "#/definitions/Error"

// SyncRepos represents the API handler to
// synchronize organization repositories between
// SCM Service and the database should a discrepancy
// exist. Common after deleting SCM repos.
func SyncRepos(c *gin.Context) {
	// capture middleware values
	u := user.Retrieve(c)
	org := c.Param("org")
	logrus.Infof("Reading repos for org %s", org)

	// See if the user is an org admin to bypass individual permission checks
	perm, err := scm.FromContext(c).OrgAccess(u, org)
	if err != nil {
		logrus.Errorf("unable to get user %s access level for org %s", u.GetName(), org)
	}

	filters := map[string]string{}
	// Only show public repos to non-admins
	if perm != "admin" {
		filters["visibility"] = "public"
	}

	// send API call to capture the total number of repos for the org
	t, err := database.FromContext(c).GetOrgRepoCount(org, filters)
	if err != nil {
		retErr := fmt.Errorf("unable to get repo count for org %s: %w", org, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	repos := []*library.Repo{}
	page := 0
	// capture all repos belonging to a certain org in database
	// nolint: gomnd // ignore magic number
	for orgRepos := int64(0); orgRepos < t; orgRepos += 100 {
		r, err := database.FromContext(c).GetOrgRepoList(org, filters, page, 100)
		if err != nil {
			retErr := fmt.Errorf("unable to get repo count for org %s: %w", org, err)

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}
		repos = append(repos, r...)
		page++
	}

	// iterate through captured repos and check if they are in GitHub
	for _, repo := range repos {
		_, err := scm.FromContext(c).GetRepo(u, repo)
		// if repo cannot be captured from GitHub, set to inactive in database
		if err != nil {
			repo.SetActive(false)
			e := database.FromContext(c).UpdateRepo(repo)
			if e != nil {
				retErr := fmt.Errorf("unable to update repo for org %s: %w", org, err)

				util.HandleError(c, http.StatusInternalServerError, retErr)

				return
			}
		}
	}
	c.JSON(http.StatusOK, fmt.Sprintf("org %s repos synced", org))
}

// swagger:operation GET /api/v1/scm/repos/{org}/{repo}/sync scm SyncRepo
//
// Sync up scm service and database in the context of a specific repo
//
// ---
// produces:
// - application/json
// parameters:
// - in: path
//   name: org
//   description: Name of the org
//   required: true
//   type: string
// - in: path
//   name: repo
//   description: Name of the repo
//   required: true
//   type: string
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully synchronized repo
//     schema:
//     type: string
//   '500':
//     description: Unable to synchronize repo
//     schema:
//       "$ref": "#/definitions/Error"

// SyncRepo represents the API handler to
// synchronize a single repository between
// SCM service and the database should a discrepancy
// exist. Common after deleting SCM repos.
func SyncRepo(c *gin.Context) {
	logrus.Infof("Reading repo %s/%s", c.Param("org"), c.Param("repo"))
	// capture middleware values
	u := user.Retrieve(c)
	org := c.Param("org")
	repo := c.Param("repo")

	// retrieve repo from context
	r, _ := database.FromContext(c).GetRepo(org, repo)

	// retrieve repo from source code manager service
	_, err := scm.FromContext(c).GetRepo(u, r)

	// if there is an error retrieving repo, we know it is deleted: sync time
	if err != nil {
		// set repo to inactive - do not delete
		r.SetActive(false)
		// update repo in database
		e := database.FromContext(c).UpdateRepo(r)
		if e != nil {
			retErr := fmt.Errorf("unable to update repo for org %s: %w", org, err)

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}
	}
	c.JSON(http.StatusOK, fmt.Sprintf("repo %s synced", r.GetFullName()))
}
