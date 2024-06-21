// SPDX-License-Identifier: Apache-2.0

package secret

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/secret"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/constants"
)

// swagger:operation DELETE /api/v1/secrets/{engine}/{type}/{org}/{name}/{secret} secrets DeleteSecret
//
// Delete a secret
//
// ---
// produces:
// - application/json
// parameters:
// - in: path
//   name: engine
//   description: Secret engine to delete the secret from, eg. "native"
//   required: true
//   type: string
// - in: path
//   name: type
//   description: Secret type to delete
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
//   description: Name of the repository if a repository secret, team name if a shared secret, or '*' if an organzation secret
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
//     description: Successfully deleted the secret
//     schema:
//       type: string
//   '401':
//     description: Unauthorized
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unexpected server error
//     schema:
//       "$ref": "#/definitions/Error"

// DeleteSecret deletes a secret from the provided secrets service.
func DeleteSecret(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
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

	// check if secret is a shared secret
	if strings.EqualFold(t, constants.SecretShared) {
		// update log fields from API metadata
		delete(fields, "secret_repo")
		fields["secret_team"] = n
	}

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logger := l.WithFields(fields)

	logger.Debugf("deleting secret %s from %s service", entry, e)

	// send API call to remove the secret
	err := secret.FromContext(c, e).Delete(ctx, t, o, n, s)
	if err != nil {
		retErr := fmt.Errorf("unable to delete secret %s from %s service: %w", entry, e, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	logger.Infof("secret %s deleted from %s service", entry, e)

	c.JSON(http.StatusOK, fmt.Sprintf("secret %s deleted from %s service", entry, e))
}
