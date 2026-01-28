// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/util"
)

// swagger:operation PUT /api/v1/repos/{org}/{repo} repos UpdateRepo
//
// Update a repository
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
// - in: body
//   name: body
//   description: The repository object with the fields to be updated
//   required: true
//   schema:
//     "$ref": "#/definitions/Repo"
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully updated the repo
//     schema:
//       "$ref": "#/definitions/Repo"
//   '400':
//     description: Invalid request payload or path
//     schema:
//       "$ref": "#/definitions/Error"
//   '401':
//     description: Unauthorized
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unexpected server error
//     schema:
//       "$ref": "#/definitions/Error"
//   '503':
//     description: Unable to update the repo
//     schema:
//       "$ref": "#/definitions/Error"

// UpdateRepo represents the API handler to update a repo.
//
//nolint:funlen,gocyclo // ignore function length
func UpdateRepo(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
	r := repo.Retrieve(c)
	u := user.Retrieve(c)
	maxBuildLimit := c.Value("maxBuildLimit").(int32)
	defaultRepoEvents := c.Value("defaultRepoEvents").([]string)
	defaultRepoEventsMask := c.Value("defaultRepoEventsMask").(int64)
	ctx := c.Request.Context()

	l.Debug("updating repo")

	// capture body from API request
	input := new(types.Repo)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for repo %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	eventsChanged := false

	// update repo fields if provided
	if len(input.GetBranch()) > 0 {
		// update branch if set
		r.SetBranch(input.GetBranch())
	}

	// update build limit if set
	if input.GetBuildLimit() > 0 {
		// allow build limit between 1 - value configured by server
		limit := max(constants.BuildLimitMin, min(input.GetBuildLimit(), maxBuildLimit))
		r.SetBuildLimit(limit)
	}

	if input.GetTimeout() > 0 {
		// update build timeout if set
		limit := max(constants.BuildTimeoutMin, min(input.GetTimeout(), constants.BuildTimeoutMax))
		r.SetTimeout(limit)
	}

	if input.GetApprovalTimeout() > 0 {
		// update build approval timeout if set
		limit := max(constants.ApprovalTimeoutMin, min(input.GetApprovalTimeout(), constants.ApprovalTimeoutMax))
		r.SetApprovalTimeout(limit)
	}

	if input.GetCounter() > 0 {
		if input.GetCounter() <= r.GetCounter() {
			retErr := fmt.Errorf("unable to set counter for repo %s: must be greater than current %d",
				r.GetFullName(), r.GetCounter())

			util.HandleError(c, http.StatusBadRequest, retErr)

			return
		}

		r.SetCounter(input.GetCounter())
	}

	if len(input.GetVisibility()) > 0 {
		// update visibility if set
		r.SetVisibility(input.GetVisibility())
	}

	if len(input.GetApproveBuild()) > 0 {
		// ensure the approve build setting matches one of the expected values
		if input.GetApproveBuild() != constants.ApproveForkAlways &&
			input.GetApproveBuild() != constants.ApproveForkNoWrite &&
			input.GetApproveBuild() != constants.ApproveNever &&
			input.GetApproveBuild() != constants.ApproveOnce {
			retErr := fmt.Errorf("approve_build of %s is invalid", input.GetApproveBuild())

			util.HandleError(c, http.StatusBadRequest, retErr)

			return
		}

		// update fork policy if set
		r.SetApproveBuild(input.GetApproveBuild())
	}

	if input.Private != nil {
		// update private if set
		r.SetPrivate(input.GetPrivate())
	}

	if input.Active != nil {
		// update active if set
		r.SetActive(input.GetActive())
	}

	// set allow events based on input if given
	if input.AllowEvents != nil {
		r.SetAllowEvents(input.GetAllowEvents())

		eventsChanged = true
	}

	// set merge queue events based on input if given
	if input.MergeQueueEvents != nil {
		for _, event := range input.GetMergeQueueEvents() {
			// only allow events possibly related to a PR merge queue
			if !slices.Contains([]string{constants.EventPush, constants.EventPull, constants.EventComment}, event) {
				retErr := fmt.Errorf("merge_queue_event of %s is invalid", event)

				util.HandleError(c, http.StatusBadRequest, retErr)

				return
			}
		}

		r.SetMergeQueueEvents(input.GetMergeQueueEvents())
	}

	// set default events if no events are enabled
	if r.GetAllowEvents().ToDatabase() == 0 {
		r.SetAllowEvents(defaultAllowedEvents(defaultRepoEvents, defaultRepoEventsMask))
	}

	if len(input.GetPipelineType()) != 0 {
		// ensure the pipeline type matches one of the expected values
		if input.GetPipelineType() != constants.PipelineTypeYAML &&
			input.GetPipelineType() != constants.PipelineTypeGo &&
			input.GetPipelineType() != constants.PipelineTypeStarlark {
			retErr := fmt.Errorf("pipeline_type of %s is invalid", input.GetPipelineType())

			util.HandleError(c, http.StatusBadRequest, retErr)

			return
		}

		r.SetPipelineType(input.GetPipelineType())
	}

	// set hash for repo if no hash is already set
	if len(r.GetHash()) == 0 {
		// create unique id for the repo
		uid, err := uuid.NewRandom()
		if err != nil {
			retErr := fmt.Errorf("unable to create UID for repo %s: %w", r.GetFullName(), err)

			util.HandleError(c, http.StatusServiceUnavailable, retErr)

			return
		}

		r.SetHash(
			base64.StdEncoding.EncodeToString(
				[]byte(strings.TrimSpace(uid.String())),
			),
		)
	}

	// fields restricted to platform admins
	if u.GetAdmin() {
		// trusted
		if input.GetTrusted() != r.GetTrusted() {
			r.SetTrusted(input.GetTrusted())
		}
	}

	// if webhook validation is not set or events didn't change, skip webhook update
	if c.Value("webhookvalidation").(bool) && eventsChanged {
		// grab last hook from repo to fetch the webhook ID
		lastHook, err := database.FromContext(c).GetHookForRepo(ctx, r, r.GetHookCounter())
		if err != nil {
			retErr := fmt.Errorf("unable to retrieve last hook for repo %s: %w", r.GetFullName(), err)

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}
		// if user is platform admin, use repo owner token to make changes to webhook
		if u.GetAdmin() {
			// capture admin name for logging
			admn := u.GetName()

			// log admin override update repo hook
			l.Debugf("platform admin %s updating repo webhook events for repo %s", admn, r.GetFullName())

			u = r.GetOwner()
		}
		// update webhook with new events
		_, err = scm.FromContext(c).Update(ctx, u, r, lastHook.GetWebhookID())
		if err != nil {
			retErr := fmt.Errorf("unable to update repo webhook for %s: %w", r.GetFullName(), err)

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}
	}

	// send API call to update the repo
	r, err = database.FromContext(c).UpdateRepo(ctx, r)
	if err != nil {
		retErr := fmt.Errorf("unable to update repo %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, r)
}
