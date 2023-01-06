// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package repo

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/util"
	"github.com/sirupsen/logrus"
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

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"org":  o,
		"repo": r.GetName(),
		"user": u.GetName(),
	}).Infof("repairing repo %s", r.GetFullName())

	// send API call to remove the webhook
	err := scm.FromContext(c).Disable(u, r.GetOrg(), r.GetName())
	if err != nil {
		retErr := fmt.Errorf("unable to delete webhook for %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// send API call to create the webhook
	_, err = scm.FromContext(c).Enable(u, r.GetOrg(), r.GetName(), r.GetHash())
	if err != nil {
		retErr := fmt.Errorf("unable to create webhook for %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// if the repo was previously inactive, mark it as active
	if !r.GetActive() {
		r.SetActive(true)

		// send API call to update the repo
		err = database.FromContext(c).UpdateRepo(r)
		if err != nil {
			retErr := fmt.Errorf("unable to set repo %s to active: %w", r.GetFullName(), err)

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}
	}

	c.JSON(http.StatusOK, fmt.Sprintf("repo %s repaired", r.GetFullName()))
}
