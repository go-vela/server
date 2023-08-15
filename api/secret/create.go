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
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/secret"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

//
// swagger:operation POST /api/v1/secrets/{engine}/{type}/{org}/{name} secrets CreateSecret
//
// Create a secret
//
// ---
// produces:
// - application/json
// parameters:
// - in: path
//   name: engine
//   description: Secret engine to create a secret in, eg. "native"
//   required: true
//   type: string
// - in: path
//   name: type
//   description: Secret type to create
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
//     description: Successfully created the secret
//     schema:
//       "$ref": "#/definitions/Secret"
//   '400':
//     description: Unable to create the secret
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to create the secret
//     schema:
//       "$ref": "#/definitions/Error"

// CreateSecret represents the API handler to
// create a secret in the configured backend.
//
//nolint:funlen // suppress long function error
func CreateSecret(c *gin.Context) {
	// capture middleware values
	u := user.Retrieve(c)
	e := util.PathParameter(c, "engine")
	t := util.PathParameter(c, "type")
	o := util.PathParameter(c, "org")
	n := util.PathParameter(c, "name")

	entry := fmt.Sprintf("%s/%s/%s", t, o, n)

	// create log fields from API metadata
	fields := logrus.Fields{
		"engine": e,
		"org":    o,
		"repo":   n,
		"type":   t,
		"user":   u.GetName(),
	}

	// check if secret is a shared secret
	if strings.EqualFold(t, constants.SecretShared) {
		// update log fields from API metadata
		fields = logrus.Fields{
			"engine": e,
			"org":    o,
			"team":   n,
			"type":   t,
			"user":   u.GetName(),
		}
	}

	if strings.EqualFold(t, constants.SecretOrg) {
		// retrieve org name from SCM
		//
		// SCM can be case insensitive, causing access retrieval to work
		// but Org/Repo != org/repo in Vela. So this check ensures that
		// what a user inputs matches the casing we expect in Vela since
		// the SCM will have the source of truth for casing.
		org, err := scm.FromContext(c).GetOrgName(u, o)
		if err != nil {
			retErr := fmt.Errorf("unable to retrieve organization %s", o)

			util.HandleError(c, http.StatusNotFound, retErr)

			return
		}

		// check if casing is accurate
		if org != o {
			retErr := fmt.Errorf("unable to retrieve organization %s. Did you mean %s?", o, org)

			util.HandleError(c, http.StatusNotFound, retErr)

			return
		}
	}

	if strings.EqualFold(t, constants.SecretRepo) {
		// retrieve org and repo name from SCM
		//
		// same story as org secret. SCM has accurate casing.
		scmOrg, scmRepo, err := scm.FromContext(c).GetOrgAndRepoName(u, o, n)
		if err != nil {
			retErr := fmt.Errorf("unable to retrieve repository %s/%s", o, n)

			util.HandleError(c, http.StatusNotFound, retErr)

			return
		}

		// check if casing is accurate for org entry
		if scmOrg != o {
			retErr := fmt.Errorf("unable to retrieve org %s. Did you mean %s?", o, scmOrg)

			util.HandleError(c, http.StatusNotFound, retErr)

			return
		}

		// check if casing is accurate for repo entry
		if scmRepo != n {
			retErr := fmt.Errorf("unable to retrieve repository %s. Did you mean %s?", n, scmRepo)

			util.HandleError(c, http.StatusNotFound, retErr)

			return
		}
	}

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(fields).Infof("creating new secret %s for %s service", entry, e)

	// capture body from API request
	input := new(library.Secret)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for secret %s for %s service: %w", entry, e, err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// reject secrets with solely whitespace characters as its value
	trimmed := strings.TrimSpace(input.GetValue())
	if len(trimmed) == 0 {
		retErr := fmt.Errorf("secret value must contain non-whitespace characters")

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// update fields in secret object
	input.SetOrg(o)
	input.SetRepo(n)
	input.SetType(t)
	input.SetCreatedAt(time.Now().UTC().Unix())
	input.SetCreatedBy(u.GetName())
	input.SetUpdatedAt(time.Now().UTC().Unix())
	input.SetUpdatedBy(u.GetName())

	if len(input.GetImages()) > 0 {
		input.SetImages(util.Unique(input.GetImages()))
	}

	if len(input.GetEvents()) > 0 {
		input.SetEvents(util.Unique(input.GetEvents()))
	}

	if len(input.GetEvents()) == 0 {
		// set default events to enable for the secret
		input.SetEvents([]string{constants.EventPush, constants.EventTag, constants.EventDeploy})
	}

	if input.AllowCommand == nil {
		input.SetAllowCommand(true)
	}

	// check if secret is a shared secret
	if strings.EqualFold(t, constants.SecretShared) {
		// update the team instead of repo
		input.SetTeam(n)
		input.Repo = nil
	}

	// send API call to create the secret
	s, err := secret.FromContext(c, e).Create(t, o, n, input)
	if err != nil {
		retErr := fmt.Errorf("unable to create secret %s for %s service: %w", entry, e, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, s.Sanitize())
}
