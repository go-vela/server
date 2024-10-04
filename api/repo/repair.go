// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	wh "github.com/go-vela/server/api/webhook"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/internal"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/util"
)

// swagger:operation PATCH /api/v1/repos/{org}/{repo}/repair repos RepairRepo
//
// Repair a hook for a repository in Vela and the configured SCM
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
//     description: Successfully repaired the repo
//     schema:
//       type: string
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
//   '500':
//     description: Unexpected server error
//     schema:
//       "$ref": "#/definitions/Error"

// RepairRepo represents the API handler to remove
// and then create a webhook for a repo.
func RepairRepo(c *gin.Context) {
	// capture middleware values
	m := c.MustGet("metadata").(*internal.Metadata)
	l := c.MustGet("logger").(*logrus.Entry)
	r := repo.Retrieve(c)
	u := user.Retrieve(c)
	ctx := c.Request.Context()

	l.Debugf("repairing repo %s", r.GetFullName())

	// todo: get org app installation
	// doesnt exist? redirect them and wait...

	// todo: from org installation, check if this repo is visible/enabled
	// no? use scm api to add the repo to the org

	// check if we should create the webhook
	if c.Value("webhookvalidation").(bool) {
		// send API call to remove the webhook
		err := scm.FromContext(c).Disable(ctx, u, r.GetOrg(), r.GetName())
		if err != nil {
			retErr := fmt.Errorf("unable to delete webhook for %s: %w", r.GetFullName(), err)

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}

		hook, err := database.FromContext(c).LastHookForRepo(ctx, r)
		if err != nil {
			retErr := fmt.Errorf("unable to get last hook for %s: %w", r.GetFullName(), err)

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}

		// send API call to create the webhook
		hook, _, err = scm.FromContext(c).Enable(ctx, u, r, hook)
		if err != nil {
			retErr := fmt.Errorf("unable to create webhook for %s: %w", r.GetFullName(), err)

			switch err.Error() {
			case "repo already enabled":
				util.HandleError(c, http.StatusConflict, retErr)
				return
			case "repo not found":
				util.HandleError(c, http.StatusNotFound, retErr)
				return
			}

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}

		hook.SetRepo(r)

		_, err = database.FromContext(c).CreateHook(ctx, hook)
		if err != nil {
			retErr := fmt.Errorf("unable to create initialization webhook for %s: %w", r.GetFullName(), err)

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}

		l.WithFields(logrus.Fields{
			"hook": hook.GetID(),
		}).Info("new webhook created")
	}

	// get repo information from the source
	sourceRepo, _, err := scm.FromContext(c).GetRepo(ctx, u, r)
	if err != nil {
		retErr := fmt.Errorf("unable to retrieve repo info for %s from source: %w", sourceRepo.GetFullName(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// if repo has a name change, then update DB with new name
	// if repo has an org change, update org as well
	if sourceRepo.GetName() != r.GetName() || sourceRepo.GetOrg() != r.GetOrg() {
		h, err := database.FromContext(c).LastHookForRepo(ctx, r)
		if err != nil {
			retErr := fmt.Errorf("unable to get last hook for %s: %w", r.GetFullName(), err)

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}

		// set sourceRepo PreviousName to old name if name is changed
		// ignore if repo is transferred and name is unchanged
		if sourceRepo.GetName() != r.GetName() {
			sourceRepo.SetPreviousName(r.GetName())
		}

		r, err = wh.RenameRepository(ctx, h, sourceRepo, c, m)
		if err != nil {
			util.HandleError(c, http.StatusInternalServerError, err)
			return
		}
	}

	// if the repo was previously inactive, mark it as active
	if !r.GetActive() {
		r.SetActive(true)

		// send API call to update the repo
		_, err := database.FromContext(c).UpdateRepo(ctx, r)
		if err != nil {
			retErr := fmt.Errorf("unable to set repo %s to active: %w", r.GetFullName(), err)

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}

		l.Infof("repo %s updated - set to active", r.GetFullName())
	}

	c.JSON(http.StatusOK, fmt.Sprintf("repo %s repaired", r.GetFullName()))
}
