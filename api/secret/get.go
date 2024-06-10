// SPDX-License-Identifier: Apache-2.0

package secret

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/router/middleware/claims"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/secret"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/constants"
)

//
// swagger:operation GET /api/v1/secrets/{engine}/{type}/{org}/{name}/{secret} secrets GetSecret
//
// Get a secret
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
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved the secret
//     schema:
//       "$ref": "#/definitions/Secret"
//   '401':
//     description: Unauthorized
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unexpected server error
//     schema:
//       "$ref": "#/definitions/Error"

// GetSecret gets a secret from the provided secrets service.
func GetSecret(c *gin.Context) {
	// capture middleware values
	cl := claims.Retrieve(c)
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
	logrus.WithFields(fields).Debugf("reading secret %s from %s service", entry, e)

	// send API call to capture the secret
	secret, err := secret.FromContext(c, e).Get(ctx, t, o, n, s)
	if err != nil {
		retErr := fmt.Errorf("unable to get secret %s from %s service: %w", entry, e, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// only allow workers to access the full secret with the value
	if strings.EqualFold(cl.TokenType, constants.WorkerBuildTokenType) {
		c.JSON(http.StatusOK, secret)

		return
	}

	c.JSON(http.StatusOK, secret.Sanitize())
}
