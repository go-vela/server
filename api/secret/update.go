// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package secret

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/secret"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

//
// swagger:operation PUT /api/v1/secrets/{engine}/{type}/{org}/{name}/{secret} secrets UpdateSecret
//
// Update a secret on the configured backend
//
// ---
// produces:
// - application/json
// parameters:
// - in: path
//   name: engine
//   description: Secret engine to update the secret in, eg. "native"
//   required: true
//   type: string
// - in: path
//   name: type
//   description: Secret type to update
//   enum:
//   - org
//   - repo
//   - shared
//   required: true
//   type: string
// - in: path
//   name: org
//   description: Name of the org
//   required: true
//   type: string
// - in: path
//   name: name
//   description: Name of the repo if a repo secret, team name if a shared secret, or '*' if an org secret
//   required: true
//   type: string
// - in: path
//   name: secret
//   description: Name of the secret
//   required: true
//   type: string
// - in: body
//   name: body
//   description: Payload containing the secret to create
//   required: true
//   schema:
//     "$ref": "#/definitions/Secret"
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully updated the secret
//     schema:
//       "$ref": "#/definitions/Secret"
//   '400':
//     description: Unable to update the secret
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to update the secret
//     schema:
//       "$ref": "#/definitions/Error"

// UpdateSecret updates a secret for the provided secrets service.
func UpdateSecret(c *gin.Context) {
	// capture middleware values
	u := user.Retrieve(c)
	e := util.PathParameter(c, "engine")
	t := util.PathParameter(c, "type")
	o := util.PathParameter(c, "org")
	n := util.PathParameter(c, "name")
	s := strings.TrimPrefix(util.PathParameter(c, "secret"), "/")
	ctx := c.Request.Context()

	entry := fmt.Sprintf("%s/%s/%s/%s", t, o, n, s)

	// create log fields from API metadata
	fields := logrus.Fields{
		"engine": e,
		"org":    o,
		"repo":   n,
		"secret": s,
		"type":   t,
		"user":   u.GetName(),
	}

	// check if secret is a shared secret
	if strings.EqualFold(t, constants.SecretShared) {
		// update log fields from API metadata
		fields = logrus.Fields{
			"engine": e,
			"org":    o,
			"secret": s,
			"team":   n,
			"type":   t,
			"user":   u.GetName(),
		}
	}

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(fields).Infof("updating secret %s for %s service", entry, e)

	// capture body from API request
	input := new(library.Secret)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for secret %s for %s service: %w", entry, e, err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// update secret fields if provided
	input.SetName(s)
	input.SetOrg(o)
	input.SetRepo(n)
	input.SetType(t)
	input.SetUpdatedAt(time.Now().UTC().Unix())
	input.SetUpdatedBy(u.GetName())

	if input.Images != nil {
		// update images if set
		input.SetImages(util.Unique(input.GetImages()))
	}

	if len(input.GetEvents()) > 0 {
		input.SetEvents(util.Unique(input.GetEvents()))
	}

	if input.AllowCommand != nil {
		// update allow_command if set
		input.SetAllowCommand(input.GetAllowCommand())
	}

	// check if secret is a shared secret
	if strings.EqualFold(t, constants.SecretShared) {
		// update the team instead of repo
		input.SetTeam(n)
		input.Repo = nil
	}

	// send API call to update the secret
	secret, err := secret.FromContext(c, e).Update(ctx, t, o, n, input)
	if err != nil {
		retErr := fmt.Errorf("unable to update secret %s for %s service: %w", entry, e, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, secret.Sanitize())
}
