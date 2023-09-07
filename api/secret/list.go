// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package secret

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/api"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/secret"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

//
// swagger:operation GET /api/v1/secrets/{engine}/{type}/{org}/{name} secrets ListSecrets
//
// Retrieve a list of secrets from the configured backend
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
// - in: query
//   name: page
//   description: The page of results to retrieve
//   type: integer
//   default: 1
// - in: query
//   name: per_page
//   description: How many results per page to return
//   type: integer
//   maximum: 100
//   default: 10
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved the list of secrets
//     schema:
//       type: array
//       items:
//         "$ref": "#/definitions/Secret"
//     headers:
//       X-Total-Count:
//         description: Total number of results
//         type: integer
//       Link:
//         description: see https://tools.ietf.org/html/rfc5988
//         type: string
//   '400':
//     description: Unable to retrieve the list of secrets
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to retrieve the list of secrets
//     schema:
//       "$ref": "#/definitions/Error"

// ListSecrets represents the API handler to capture
// a list of secrets from the configured backend.
func ListSecrets(c *gin.Context) {
	// capture middleware values
	u := user.Retrieve(c)
	e := util.PathParameter(c, "engine")
	t := util.PathParameter(c, "type")
	o := util.PathParameter(c, "org")
	n := util.PathParameter(c, "name")
	ctx := c.Request.Context()

	var teams []string
	// get list of user's teams if type is shared secret and team is '*'
	if t == constants.SecretShared && n == "*" {
		var err error

		teams, err = scm.FromContext(c).ListUsersTeamsForOrg(u, o)
		if err != nil {
			retErr := fmt.Errorf("unable to list users %s teams for org %s: %w", u.GetName(), o, err)

			util.HandleError(c, http.StatusBadRequest, retErr)

			return
		}
	}

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

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(fields).Infof("listing secrets %s from %s service", entry, e)

	// capture page query parameter if present
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		retErr := fmt.Errorf("unable to convert page query parameter for %s from %s service: %w", entry, e, err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// capture per_page query parameter if present
	perPage, err := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	if err != nil {
		retErr := fmt.Errorf("unable to convert per_page query parameter for %s from %s service: %w", entry, e, err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// send API call to capture the total number of secrets
	total, err := secret.FromContext(c, e).Count(ctx, t, o, n, teams)
	if err != nil {
		retErr := fmt.Errorf("unable to get secret count for %s from %s service: %w", entry, e, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// ensure per_page isn't above or below allowed values
	perPage = util.MaxInt(1, util.MinInt(100, perPage))

	// send API call to capture the list of secrets
	s, err := secret.FromContext(c, e).List(ctx, t, o, n, page, perPage, teams)
	if err != nil {
		retErr := fmt.Errorf("unable to list secrets for %s from %s service: %w", entry, e, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// create pagination object
	pagination := api.Pagination{
		Page:    page,
		PerPage: perPage,
		Total:   total,
	}
	// set pagination headers
	pagination.SetHeaderLink(c)

	// variable we want to return
	secrets := []*library.Secret{}
	// iterate through all secrets
	for _, secret := range s {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := secret

		// sanitize secret to ensure no value is provided
		secrets = append(secrets, tmp.Sanitize())
	}

	c.JSON(http.StatusOK, secrets)
}
