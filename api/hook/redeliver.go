// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package hook

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/util"
	"github.com/sirupsen/logrus"
)

// swagger:operation POST /api/v1/hooks/{org}/{repo}/{hook}/redeliver webhook RedeliverHook
//
// Redeliver a webhook from the SCM
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
// - in: path
//   name: hook
//   description: Number of the hook
//   required: true
//   type: integer
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully redelivered the webhook
//     schema:
//       "$ref": "#/definitions/Webhook"
//   '400':
//     description: The webhook was unable to be redelivered
//     schema:
//       "$ref": "#/definitions/Error"
//   '404':
//     description: The webhook was unable to be redelivered
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: The webhook was unable to be redelivered
//     schema:
//       "$ref": "#/definitions/Error"

// RedeliverHook represents the API handler to redeliver
// a webhook from the SCM.
func RedeliverHook(c *gin.Context) {
	// capture middleware values
	o := org.Retrieve(c)
	r := repo.Retrieve(c)
	u := user.Retrieve(c)
	hook := util.PathParameter(c, "hook")

	entry := fmt.Sprintf("%s/%s", r.GetFullName(), hook)

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"org":  o,
		"hook": hook,
		"repo": r.GetName(),
		"user": u.GetName(),
	}).Infof("redelivering hook %s", entry)

	number, err := strconv.Atoi(hook)
	if err != nil {
		retErr := fmt.Errorf("invalid hook parameter provided: %s", hook)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// send API call to capture the webhook
	h, err := database.FromContext(c).GetHookForRepo(r, number)
	if err != nil {
		retErr := fmt.Errorf("unable to get hook %s: %w", entry, err)

		util.HandleError(c, http.StatusNotFound, retErr)

		return
	}

	err = scm.FromContext(c).RedeliverWebhook(c, u, r, h)
	if err != nil {
		retErr := fmt.Errorf("unable to redeliver hook %s: %w", entry, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("hook %s redelivered", entry))
}
