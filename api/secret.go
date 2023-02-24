// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/router/middleware/claims"
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
		input.SetImages(unique(input.GetImages()))
	}

	if len(input.GetEvents()) > 0 {
		input.SetEvents(unique(input.GetEvents()))
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
	err = secret.FromContext(c, e).Create(t, o, n, input)
	if err != nil {
		retErr := fmt.Errorf("unable to create secret %s for %s service: %w", entry, e, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	s, _ := secret.FromContext(c, e).Get(t, o, n, input.GetName())

	c.JSON(http.StatusOK, s.Sanitize())
}

//
// swagger:operation GET /api/v1/secrets/{engine}/{type}/{org}/{name} secrets GetSecrets
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

// GetSecrets represents the API handler to capture
// a list of secrets from the configured backend.
func GetSecrets(c *gin.Context) {
	// capture middleware values
	u := user.Retrieve(c)
	e := util.PathParameter(c, "engine")
	t := util.PathParameter(c, "type")
	o := util.PathParameter(c, "org")
	n := util.PathParameter(c, "name")

	var teams []string
	// get list of user's teams if type is shared secret and team is '*'
	if t == constants.SecretShared && n == "*" {
		var err error

		teams, err = scm.FromContext(c).ListUsersTeamsForOrg(u, o)
		if err != nil {
			retErr := fmt.Errorf("unable to get users %s teams for org %s: %w", u.GetName(), o, err)

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
	logrus.WithFields(fields).Infof("reading secrets %s from %s service", entry, e)

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
	total, err := secret.FromContext(c, e).Count(t, o, n, teams)
	if err != nil {
		retErr := fmt.Errorf("unable to get secret count for %s from %s service: %w", entry, e, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// ensure per_page isn't above or below allowed values
	perPage = util.MaxInt(1, util.MinInt(100, perPage))

	// send API call to capture the list of secrets
	s, err := secret.FromContext(c, e).List(t, o, n, page, perPage, teams)
	if err != nil {
		retErr := fmt.Errorf("unable to get secrets for %s from %s service: %w", entry, e, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// create pagination object
	pagination := Pagination{
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

//
// swagger:operation GET /api/v1/secrets/{engine}/{type}/{org}/{name}/{secret} secrets GetSecret
//
// Retrieve a secret from the configured backend
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
//   '500':
//     description: Unable to retrieve the secret
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
	logrus.WithFields(fields).Infof("reading secret %s from %s service", entry, e)

	// send API call to capture the secret
	secret, err := secret.FromContext(c, e).Get(t, o, n, s)
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

//
// swagger:operation PUT /api/v1/secrets/{engine}/{type}/{org}/{name}/{secret} secrets UpdateSecrets
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
		input.SetImages(unique(input.GetImages()))
	}

	if len(input.GetEvents()) > 0 {
		input.SetEvents(unique(input.GetEvents()))
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
	err = secret.FromContext(c, e).Update(t, o, n, input)
	if err != nil {
		retErr := fmt.Errorf("unable to update secret %s for %s service: %w", entry, e, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// send API call to capture the updated secret
	secret, _ := secret.FromContext(c, e).Get(t, o, n, input.GetName())

	c.JSON(http.StatusOK, secret.Sanitize())
}

//
// swagger:operation DELETE /api/v1/secrets/{engine}/{type}/{org}/{name}/{secret} secrets DeleteSecret
//
// Delete a secret from the configured backend
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
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully deleted the secret
//     schema:
//       type: string
//   '500':
//     description: Unable to delete the secret
//     schema:
//       "$ref": "#/definitions/Error"

// DeleteSecret deletes a secret from the provided secrets service.
func DeleteSecret(c *gin.Context) {
	// capture middleware values
	u := user.Retrieve(c)
	e := util.PathParameter(c, "engine")
	t := util.PathParameter(c, "type")
	o := util.PathParameter(c, "org")
	n := util.PathParameter(c, "name")
	s := strings.TrimPrefix(util.PathParameter(c, "secret"), "/")

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
	logrus.WithFields(fields).Infof("deleting secret %s from %s service", entry, e)

	// send API call to remove the secret
	err := secret.FromContext(c, e).Delete(t, o, n, s)
	if err != nil {
		retErr := fmt.Errorf("unable to delete secret %s from %s service: %w", entry, e, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("secret %s deleted from %s service", entry, e))
}

// unique is a helper function that takes a slice and
// validates that there are no duplicate entries.
func unique(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}

	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true

			list = append(list, entry)
		}
	}

	return list
}
