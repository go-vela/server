// SPDX-License-Identifier: Apache-2.0

package scm

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/util"
)

// swagger:operation PATCH /api/v1/scm/repos/{org}/{repo}/sync scm SyncRepo
//
// Sync a repository with the scm service
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
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully synchronized repo
//     schema:
//       "$ref": "#/definitions/Repo"
//   '204':
//     description: Successful request resulting in no change
//   '301':
//     description: Repo has moved permanently (from SCM)
//     schema:
//       "$ref": "#/definitions/Error"
//   '401':
//     description: Unauthorized
//     schema:
//       "$ref": "#/definitions/Error"
//   '403':
//     description: User has been forbidden access to repository (from SCM)
//     schema:
//       "$ref": "#/definitions/Error"
//   '404':
//     description: Not found
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unexpected server error
//     schema:
//       "$ref": "#/definitions/Error"

// SyncRepo represents the API handler to
// synchronize a single repository between
// SCM service and the database should a discrepancy
// exist. Primarily used for deleted repos or to align
// subscribed events with allowed events.
func SyncRepo(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
	o := org.Retrieve(c)
	r := repo.Retrieve(c)
	u := user.Retrieve(c)
	ctx := c.Request.Context()

	l.Debugf("syncing repo %s", r.GetFullName())

	// retrieve repo from source code manager service
	_, respCode, err := scm.FromContext(c).GetRepo(ctx, u, r)
	// if there is an error retrieving repo, we know it is deleted: set to inactive
	if err != nil {
		if respCode == http.StatusNotFound {
			// set repo to inactive - do not delete
			r.SetActive(false)

			// update repo in database
			r, err = database.FromContext(c).UpdateRepo(ctx, r)
			if err != nil {
				retErr := fmt.Errorf("unable to update repo for org %s: %w", o, err)

				util.HandleError(c, http.StatusInternalServerError, retErr)

				return
			}

			l.Infof("repo %s has been updated - set to inactive", r.GetFullName())

			// exit with success as hook sync will be unnecessary
			c.JSON(http.StatusOK, r)

			return
		}

		retErr := fmt.Errorf("error while retrieving repo %s from %s: %w", r.GetFullName(), scm.FromContext(c).Driver(), err)

		util.HandleError(c, respCode, retErr)

		return
	}

	// verify the user is an admin of the repo
	// we cannot use our normal permissions check due to the possibility the repo was deleted
	perm, err := scm.FromContext(c).RepoAccess(ctx, u.GetName(), u.GetToken(), o, r.GetName())
	if err != nil {
		l.Errorf("unable to get user %s access level for org %s", u.GetName(), o)
	}

	if perm != constants.PermissionAdmin {
		retErr := fmt.Errorf("user %s does not have 'admin' permissions for the repo %s", u.GetName(), r.GetFullName())

		util.HandleError(c, http.StatusUnauthorized, retErr)

		return
	}

	// if we have webhook validation and the repo is active in the database,
	// update the repo hook in the SCM
	if c.Value("webhookvalidation").(bool) && r.GetActive() {
		// grab last hook from repo to fetch the webhook ID
		lastHook, err := database.FromContext(c).GetHookForRepo(ctx, r, r.GetHookCounter())
		if err != nil {
			retErr := fmt.Errorf("unable to retrieve last hook for repo %s: %w", r.GetFullName(), err)

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}

		// update webhook
		webhookExists, err := scm.FromContext(c).Update(ctx, u, r, lastHook.GetWebhookID())
		if err != nil {
			// if webhook has been manually deleted from GitHub,
			// set to inactive in database
			if !webhookExists {
				r.SetActive(false)

				r, err = database.FromContext(c).UpdateRepo(ctx, r)
				if err != nil {
					retErr := fmt.Errorf("unable to update repo for org %s: %w", o, err)

					util.HandleError(c, http.StatusInternalServerError, retErr)

					return
				}

				l.Infof("repo %s has been updated - set to inactive", r.GetFullName())

				c.JSON(http.StatusOK, r)

				return
			}

			retErr := fmt.Errorf("unable to update repo webhook for %s: %w", r.GetFullName(), err)

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}
	}

	// map this repo to an installation, if necessary
	installID := r.GetInstallID()

	r, err = scm.FromContext(c).SyncRepoWithInstallation(ctx, r)
	if err != nil {
		retErr := fmt.Errorf("unable to sync repo %s with installation: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// install_id was synced
	if r.GetInstallID() != installID {
		_, err := database.FromContext(c).UpdateRepo(ctx, r)
		if err != nil {
			retErr := fmt.Errorf("unable to update repo %s during repair: %w", r.GetFullName(), err)

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}

		l.Tracef("repo %s install_id synced to %d", r.GetFullName(), r.GetInstallID())
	}

	c.Status(http.StatusNoContent)
}
