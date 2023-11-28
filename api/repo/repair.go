// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"fmt"
	"github.com/gin-gonic/gin"
	wh "github.com/go-vela/server/api/webhook"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types"
	"github.com/sirupsen/logrus"
	"net/http"
)

// swagger:operation PATCH /api/v1/repos/{org}/{repo}/repair repos RepairRepo
//
// Remove and recreate the webhook for a repo
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
//     description: Successfully repaired the repo
//     schema:
//       type: string
//   '500':
//     description: Unable to repair the repo
//     schema:
//       "$ref": "#/definitions/Error"

// RepairRepo represents the API handler to remove
// and then create a webhook for a repo.
func RepairRepo(c *gin.Context) {
	// capture middleware values
	o := org.Retrieve(c)
	r := repo.Retrieve(c)
	u := user.Retrieve(c)
	ctx := c.Request.Context()
	// capture middleware values
	m := c.MustGet("metadata").(*types.Metadata)

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"org":  o,
		"repo": r.GetName(),
		"user": u.GetName(),
	}).Infof("repairing repo %s", r.GetFullName())

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

		hook.SetRepoID(r.GetID())

		_, err = database.FromContext(c).CreateHook(ctx, hook)
		if err != nil {
			retErr := fmt.Errorf("unable to create initialization webhook for %s: %w", r.GetFullName(), err)

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}
	}

	// get repo information from the source
	sourceRepo, err := scm.FromContext(c).GetRepo(ctx, u, r)
	if err != nil {
		retErr := fmt.Errorf("unable to retrieve repo info for %s from source: %w", sourceRepo.GetFullName(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// if repo has a name change, then update DB with new name
	// if repo has a org change, update org as well
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
	}

	c.JSON(http.StatusOK, fmt.Sprintf("repo %s repaired", r.GetFullName()))
}
