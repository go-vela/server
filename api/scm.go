// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/library"
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
// exist. Primarily used for deleted repos or to align
// subscribed events with allowed events.
func SyncRepos(c *gin.Context) {
	// capture middleware values
	o := org.Retrieve(c)
	u := user.Retrieve(c)

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logger := logrus.WithFields(logrus.Fields{
		"org":  o,
		"user": u.GetName(),
	})

	logger.Infof("syncing repos for org %s", o)

	// see if the user is an org admin
	perm, err := scm.FromContext(c).OrgAccess(u, o)
	if err != nil {
		logger.Errorf("unable to get user %s access level for org %s", u.GetName(), o)
	}

	// only allow org-wide syncing if user is admin of org
	if perm != "admin" {
		retErr := fmt.Errorf("unable to sync repos in org %s: must be an org admin", o)

		util.HandleError(c, http.StatusUnauthorized, retErr)

		return
	}

	// send API call to capture the total number of repos for the org
	t, err := database.FromContext(c).CountReposForOrg(o, map[string]interface{}{})
	if err != nil {
		retErr := fmt.Errorf("unable to get repo count for org %s: %w", o, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	repos := []*library.Repo{}
	page := 0
	// capture all repos belonging to a certain org in database
	for orgRepos := int64(0); orgRepos < t; orgRepos += 100 {
		r, _, err := database.FromContext(c).ListReposForOrg(o, "name", map[string]interface{}{}, page, 100)
		if err != nil {
			retErr := fmt.Errorf("unable to get repo count for org %s: %w", o, err)

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

			err := database.FromContext(c).UpdateRepo(repo)
			if err != nil {
				retErr := fmt.Errorf("unable to update repo for org %s: %w", o, err)

				util.HandleError(c, http.StatusInternalServerError, retErr)

				return
			}
		}

		// if we have webhook validation, update the repo hook in the SCM
		if c.Value("webhookvalidation").(bool) {
			// grab last hook from repo to fetch the webhook ID
			lastHook, err := database.FromContext(c).LastHookForRepo(repo)
			if err != nil {
				retErr := fmt.Errorf("unable to retrieve last hook for repo %s: %w", repo.GetFullName(), err)

				util.HandleError(c, http.StatusInternalServerError, retErr)

				return
			}

			// update webhook
			err = scm.FromContext(c).Update(u, repo, lastHook.GetWebhookID())
			if err != nil {
				retErr := fmt.Errorf("unable to update repo webhook for %s: %w", repo.GetFullName(), err)

				util.HandleError(c, http.StatusInternalServerError, retErr)

				return
			}
		}
	}

	c.JSON(http.StatusOK, fmt.Sprintf("org %s repos synced", o))
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
// exist. Primarily used for deleted repos or to align
// subscribed events with allowed events.
func SyncRepo(c *gin.Context) {
	// capture middleware values
	o := org.Retrieve(c)
	r := repo.Retrieve(c)
	u := user.Retrieve(c)

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logger := logrus.WithFields(logrus.Fields{
		"org":  o,
		"repo": r.GetName(),
		"user": u.GetName(),
	})

	logger.Infof("syncing repo %s", r.GetFullName())

	// retrieve repo from source code manager service
	_, err := scm.FromContext(c).GetRepo(u, r)

	// if there is an error retrieving repo, we know it is deleted: set to inactive
	if err != nil {
		// set repo to inactive - do not delete
		r.SetActive(false)

		// update repo in database
		err := database.FromContext(c).UpdateRepo(r)
		if err != nil {
			retErr := fmt.Errorf("unable to update repo for org %s: %w", o, err)

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}

		// exit with success as hook sync will be unnecessary
		c.JSON(http.StatusOK, fmt.Sprintf("repo %s synced", r.GetFullName()))

		return
	}

	// verify the user is an admin of the repo
	// we cannot use our normal permissions check due to the possibility the repo was deleted
	perm, err := scm.FromContext(c).RepoAccess(u, u.GetToken(), o, r.GetName())
	if err != nil {
		logger.Errorf("unable to get user %s access level for org %s", u.GetName(), o)
	}

	if !strings.EqualFold(perm, "admin") {
		retErr := fmt.Errorf("user %s does not have 'admin' permissions for the repo %s", u.GetName(), r.GetFullName())

		util.HandleError(c, http.StatusUnauthorized, retErr)

		return
	}

	// if we have webhook validation, update the repo hook in the SCM
	if c.Value("webhookvalidation").(bool) {
		// grab last hook from repo to fetch the webhook ID
		lastHook, err := database.FromContext(c).LastHookForRepo(r)
		if err != nil {
			retErr := fmt.Errorf("unable to retrieve last hook for repo %s: %w", r.GetFullName(), err)

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}

		// update webhook
		err = scm.FromContext(c).Update(u, r, lastHook.GetWebhookID())
		if err != nil {
			retErr := fmt.Errorf("unable to update repo webhook for %s: %w", r.GetFullName(), err)

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}
	}

	c.JSON(http.StatusOK, fmt.Sprintf("repo %s synced", r.GetFullName()))
}
