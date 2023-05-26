// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package hook

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// swagger:operation POST /api/v1/hooks/{org}/{repo} webhook CreateHook
//
// Create a webhook for the configured backend
//
// ---
// produces:
// - application/json
// parameters:
// - in: body
//   name: body
//   description: Webhook payload that we expect from the user or VCS
//   required: true
//   schema:
//     "$ref": "#/definitions/Webhook"
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
//   '201':
//     description: The webhook has been created
//     schema:
//       "$ref": "#/definitions/Webhook"
//   '400':
//     description: The webhook was unable to be created
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: The webhook was unable to be created
//     schema:
//       "$ref": "#/definitions/Error"

// CreateHook represents the API handler to create
// a webhook in the configured backend.
func CreateHook(c *gin.Context) {
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
	}).Infof("creating new hook for repo %s", r.GetFullName())

	// capture body from API request
	input := new(library.Hook)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for new hook for repo %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// send API call to capture the last hook for the repo
	lastHook, err := database.FromContext(c).LastHookForRepo(r)
	if err != nil {
		retErr := fmt.Errorf("unable to get last hook for repo %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// update fields in webhook object
	input.SetRepoID(r.GetID())
	input.SetNumber(1)

	if input.GetCreated() == 0 {
		input.SetCreated(time.Now().UTC().Unix())
	}

	if lastHook != nil {
		input.SetNumber(
			lastHook.GetNumber() + 1,
		)
	}

	// send API call to create the webhook
	h, err := database.FromContext(c).CreateHook(input)
	if err != nil {
		retErr := fmt.Errorf("unable to create hook for repo %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusCreated, h)
}
