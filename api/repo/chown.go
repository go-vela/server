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
	"github.com/go-vela/server/util"
	"github.com/sirupsen/logrus"
)

// swagger:operation PATCH /api/v1/repos/{org}/{repo}/chown repos ChownRepo
//
// Change the owner of the webhook for a repo
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
//     description: Successfully changed the owner for the repo
//     schema:
//       type: string
//   '500':
//     description: Unable to change the owner for the repo
//     schema:
//       "$ref": "#/definitions/Error"

// ChownRepo represents the API handler to change
// the owner of a repo in the configured backend.
func ChownRepo(c *gin.Context) {
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
	}).Infof("changing owner of repo %s to %s", r.GetFullName(), u.GetName())

	// update repo owner
	r.SetUserID(u.GetID())

	// send API call to updated the repo
	err := database.FromContext(c).UpdateRepo(r)
	if err != nil {
		retErr := fmt.Errorf("unable to change owner of repo %s to %s: %w", r.GetFullName(), u.GetName(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("repo %s changed owner to %s", r.GetFullName(), u.GetName()))
}
