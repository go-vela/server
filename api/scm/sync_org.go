// SPDX-License-Identifier: Apache-2.0

package scm

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/util"
)

// swagger:operation PATCH /api/v1/scm/orgs/{org}/sync scm SyncReposForOrg
//
// Sync repositories from scm service and database in a specified organization
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
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully synchronized repos
//     schema:
//       type: array
//       items:
//         "$ref": "#/definitions/Repo"
//   '204':
//     description: Successful request resulting in no change
//   '301':
//     description: One repository in the organiation has moved permanently (from SCM)
//     schema:
//       "$ref": "#/definitions/Error"
//   '400':
//     description: Invalid request payload or path
//     schema:
//       "$ref": "#/definitions/Error"
//   '401':
//     description: Unauthorized
//     schema:
//       "$ref": "#/definitions/Error"
//   '403':
//     description: User has been forbidden access to at least one repository in organiation (from SCM)
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unexpected server error
//     schema:
//       "$ref": "#/definitions/Error"

// SyncReposForOrg represents the API handler to
// synchronize organization repositories between
// SCM Service and the database should a discrepancy
// exist. Primarily used for deleted repos or to align
// subscribed events with allowed events.
func SyncReposForOrg(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
	o := org.Retrieve(c)
	u := user.Retrieve(c)
	ctx := c.Request.Context()

	l.Debugf("syncing repos for org %s", o)

	// see if the user is an org admin
	perm, err := scm.FromContext(c).OrgAccess(ctx, u, o)
	if err != nil {
		l.Errorf("unable to get user %s access level for org %s", u.GetName(), o)
	}

	// only allow org-wide syncing if user is admin of org
	if perm != constants.PermissionAdmin {
		retErr := fmt.Errorf("unable to sync repos in org %s: must be an org admin", o)

		util.HandleError(c, http.StatusUnauthorized, retErr)

		return
	}

	// send API call to capture the total number of repos for the org
	t, err := database.FromContext(c).CountReposForOrg(ctx, o, map[string]interface{}{})
	if err != nil {
		retErr := fmt.Errorf("unable to get repo count for org %s: %w", o, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	repos := []*types.Repo{}
	page := 0
	// capture all repos belonging to a certain org in database
	for orgRepos := int64(0); orgRepos < t; orgRepos += 100 {
		r, err := database.FromContext(c).ListReposForOrg(ctx, o, "name", map[string]interface{}{}, page, 100)
		if err != nil {
			retErr := fmt.Errorf("unable to get repo count for org %s: %w", o, err)

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}

		repos = append(repos, r...)

		page++
	}

	var results []*types.Repo

	// iterate through captured repos and check if they are in GitHub
	for _, repo := range repos {
		_, respCode, err := scm.FromContext(c).GetRepo(ctx, u, repo)
		// if repo cannot be captured from GitHub, set to inactive in database
		if err != nil {
			if respCode == http.StatusNotFound {
				_, err := database.FromContext(c).UpdateRepo(ctx, repo)
				if err != nil {
					retErr := fmt.Errorf("unable to update repo for org %s: %w", o, err)

					util.HandleError(c, http.StatusInternalServerError, retErr)

					return
				}

				l.Infof("repo %s has been updated - set to inactive", repo.GetFullName())

				results = append(results, repo)
			} else {
				retErr := fmt.Errorf("error while retrieving repo %s from %s: %w", repo.GetFullName(), scm.FromContext(c).Driver(), err)

				util.HandleError(c, respCode, retErr)

				return
			}
		}

		// if we have webhook validation and the repo is active in the database,
		// update the repo hook in the SCM
		if c.Value("webhookvalidation").(bool) && repo.GetActive() {
			// grab last hook from repo to fetch the webhook ID
			lastHook, err := database.FromContext(c).LastHookForRepo(ctx, repo)
			if err != nil {
				retErr := fmt.Errorf("unable to retrieve last hook for repo %s: %w", repo.GetFullName(), err)

				util.HandleError(c, http.StatusInternalServerError, retErr)

				return
			}

			// update webhook
			webhookExists, err := scm.FromContext(c).Update(ctx, u, repo, lastHook.GetWebhookID())
			if err != nil {
				// if webhook has been manually deleted from GitHub,
				// set to inactive in database
				if !webhookExists {
					repo.SetActive(false)

					_, err := database.FromContext(c).UpdateRepo(ctx, repo)
					if err != nil {
						retErr := fmt.Errorf("unable to update repo for org %s: %w", o, err)

						util.HandleError(c, http.StatusInternalServerError, retErr)

						return
					}

					l.Infof("repo %s has been updated - set to inactive", repo.GetFullName())

					results = append(results, repo)

					continue
				}

				retErr := fmt.Errorf("unable to update repo webhook for %s: %w", repo.GetFullName(), err)

				util.HandleError(c, http.StatusInternalServerError, retErr)

				return
			}
		}
	}

	// if no repo was changed, return no content status
	if len(results) == 0 {
		c.Status(http.StatusNoContent)

		return
	}

	c.JSON(http.StatusOK, results)
}
