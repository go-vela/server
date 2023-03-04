// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package repo

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// swagger:operation PUT /api/v1/repos/{org}/{repo} repos UpdateRepo
//
// Update a repo in the configured backend
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
// - in: body
//   name: body
//   description: Payload containing the repo to update
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
//     description: Unable to update the repo
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to update the repo
//     schema:
//       "$ref": "#/definitions/Error"
//   '503':
//     description: Unable to update the repo
//     schema:
//       "$ref": "#/definitions/Error"

// UpdateRepo represents the API handler to update
// a repo in the configured backend.
//
//nolint:funlen // ignore line length
func UpdateRepo(c *gin.Context) {
	// capture middleware values
	o := org.Retrieve(c)
	r := repo.Retrieve(c)
	u := user.Retrieve(c)
	maxBuildLimit := c.Value("maxBuildLimit").(int64)

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"org":  o,
		"repo": r.GetName(),
		"user": u.GetName(),
	}).Infof("updating repo %s", r.GetFullName())

	// capture body from API request
	input := new(library.Repo)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for repo %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// update repo fields if provided
	if len(input.GetBranch()) > 0 {
		// update branch if set
		r.SetBranch(input.GetBranch())
	}

	// update build limit if set
	if input.GetBuildLimit() > 0 {
		// allow build limit between 1 - value configured by server
		r.SetBuildLimit(
			int64(
				util.MaxInt(
					constants.BuildLimitMin,
					util.MinInt(
						int(input.GetBuildLimit()),
						int(maxBuildLimit),
					), // clamp max
				), // clamp min
			),
		)
	}

	if input.GetTimeout() > 0 {
		// update build timeout if set
		r.SetTimeout(
			int64(
				util.MaxInt(
					constants.BuildTimeoutMin,
					util.MinInt(
						int(input.GetTimeout()),
						constants.BuildTimeoutMax,
					), // clamp max
				), // clamp min
			),
		)
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

	if input.Private != nil {
		// update private if set
		r.SetPrivate(input.GetPrivate())
	}

	if input.Active != nil {
		// update active if set
		r.SetActive(input.GetActive())
	}

	if input.AllowPull != nil {
		// update allow_pull if set
		r.SetAllowPull(input.GetAllowPull())
	}

	if input.AllowPush != nil {
		// update allow_push if set
		r.SetAllowPush(input.GetAllowPush())
	}

	if input.AllowDeploy != nil {
		// update allow_deploy if set
		r.SetAllowDeploy(input.GetAllowDeploy())
	}

	if input.AllowTag != nil {
		// update allow_tag if set
		r.SetAllowTag(input.GetAllowTag())
	}

	if input.AllowComment != nil {
		// update allow_comment if set
		r.SetAllowComment(input.GetAllowComment())
	}

	// set default events if no events are enabled
	if !r.GetAllowPull() && !r.GetAllowPush() &&
		!r.GetAllowDeploy() && !r.GetAllowTag() &&
		!r.GetAllowComment() {
		r.SetAllowPull(true)
		r.SetAllowPush(true)
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

	// send API call to update the repo
	err = database.FromContext(c).UpdateRepo(r)
	if err != nil {
		retErr := fmt.Errorf("unable to update repo %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// send API call to capture the updated repo
	r, _ = database.FromContext(c).GetRepoForOrg(r.GetOrg(), r.GetName())

	c.JSON(http.StatusOK, r)
}
