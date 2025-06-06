// SPDX-License-Identifier: Apache-2.0

package secret

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/secret"
	"github.com/go-vela/server/util"
)

//
// swagger:operation PUT /api/v1/secrets/{engine}/{type}/{org}/{name}/{secret} secrets UpdateSecret
//
// Update a secret
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
//   description: Name of the organization
//   required: true
//   type: string
// - in: path
//   name: name
//   description: Name of the repository if a repository secret, team name if a shared secret, or '*' if an organization secret
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

// UpdateSecret updates a secret for the provided secrets service.
func UpdateSecret(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
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
		"secret_engine": e,
		"secret_org":    o,
		"secret_repo":   n,
		"secret_name":   s,
		"secret_type":   t,
	}

	dbSecret, err := secret.FromContext(c, e).Get(ctx, t, o, n, s)
	if err != nil {
		retErr := fmt.Errorf("unable to get secret from database: %w", err)

		util.HandleError(c, http.StatusNotFound, retErr)

		return
	}

	// check if secret is a shared secret
	if strings.EqualFold(t, constants.SecretShared) {
		// update log fields from API metadata
		delete(fields, "secret_repo")
		fields["secret_team"] = n
	}

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	l.WithFields(fields).Debugf("updating secret %s for %s service", entry, e)

	// capture body from API request
	input := new(types.Secret)

	err = c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for secret %s for %s service: %w", entry, e, err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	if input.RepoAllowlist != nil {
		input.SetRepoAllowlist(util.Unique(input.GetRepoAllowlist()))
	}

	err = validateAllowlist(ctx, database.FromContext(c), dbSecret.GetRepoAllowlist(), input)
	if err != nil {
		retErr := fmt.Errorf("invalid allowlist: %w", err)

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

	if input.AllowCommand != nil {
		// update allow_command if set
		input.SetAllowCommand(input.GetAllowCommand())
	}

	if input.AllowSubstitution != nil {
		input.SetAllowSubstitution(input.GetAllowSubstitution())
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

	l.WithFields(fields).Info("secret updated")

	c.JSON(http.StatusOK, secret.Sanitize())
}
