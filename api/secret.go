// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/secret"
	"github.com/go-vela/server/util"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// swagger:operation POST /api/v1/secrets/{engine}/{type}/{org}/{name} secrets CreateSecret
//
// Create a secret
//
// ---
// x-success_http_code: '201'
// produces:
// - application/json
// parameters:
// - in: path
//   name: engine
//   description: Secret engine to create a secret in
//   required: true
//   type: string
// - in: path
//   name: type
//   description: Secret type to create. Options 'org', 'repo', or 'shared'
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
//     type: json
//     schema:
//       "$ref": "#/definitions/Secret"
//   '400':
//     description: Unable to create the secret
//     schema:
//       type: string
//   '500':
//     description: Unable to create the secret
//     schema:
//       type: string

// CreateSecret represents the API handler to
// create a secret in the configured backend.
func CreateSecret(c *gin.Context) {
	// define the default events to enable for a secret
	defaultEvents := []string{constants.EventPush, constants.EventTag, constants.EventDeploy}

	// capture middleware values
	e := c.Param("engine")
	t := c.Param("type")
	o := c.Param("org")
	n := c.Param("name")

	logrus.Infof("Creating secret %s/%s/%s for %s service", t, o, n, e)

	// capture body from API request
	input := new(library.Secret)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for secret %s/%s/%s for %s service: %w", t, o, n, e, err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// update fields in secret object
	input.SetOrg(o)
	input.SetRepo(n)
	input.SetType(t)

	if len(input.GetImages()) > 0 {
		input.SetImages(unique(input.GetImages()))
	}

	if len(input.GetEvents()) > 0 {
		input.SetEvents(unique(input.GetEvents()))
	}

	if len(input.GetEvents()) == 0 {
		// set default events to enable for the secret
		input.SetEvents(defaultEvents)
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
		retErr := fmt.Errorf("unable to create secret %s/%s/%s for %s service: %w", t, o, n, e, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	s, _ := secret.FromContext(c, e).Get(t, o, n, input.GetName())
	c.JSON(http.StatusOK, s.Sanitize())
}

// swagger:operation GET /api/v1/secrets/{engine}/{type}/{org}/{name} secrets GetSecrets
//
// Retrieve a list of secrets from the configured backend.
//
// ---
// x-success_http_code: '200'
// produces:
// - application/json
// parameters:
// - in: path
//   name: engine
//   description: Secret engine to create a secret in
//   required: true
//   type: string
// - in: path
//   name: type
//   description: Secret type to create. Options 'org', 'repo', or 'shared'
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
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved the list of secrets
//     type: json
//     schema:
//       type: array
//       items:
//         "$ref": "#/definitions/Secret"
//   '400':
//     description: Unable to retrieve the list of secrets
//     schema:
//       type: string
//   '500':
//     description: Unable to retrieve the list of secrets
//     schema:
//       type: string

// GetSecrets represents the API handler to capture
// a list of secrets from the configured backend.
func GetSecrets(c *gin.Context) {
	// capture middleware values
	e := c.Param("engine")
	t := c.Param("type")
	o := c.Param("org")
	n := c.Param("name")

	logrus.Infof("Reading secrets %s/%s/%s from %s service", t, o, n, e)

	// capture page query parameter if present
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		retErr := fmt.Errorf("unable to convert page query parameter for %s/%s/%s from %s service: %w", t, o, n, e, err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// capture per_page query parameter if present
	perPage, err := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	if err != nil {
		retErr := fmt.Errorf("unable to convert per_page query parameter for %s/%s/%s from %s service: %w", t, o, n, e, err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// send API call to capture the total number of secrets
	total, err := secret.FromContext(c, e).Count(t, o, n)
	if err != nil {
		retErr := fmt.Errorf("unable to get secret count for %s/%s/%s from %s service: %w", t, o, n, e, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// ensure per_page isn't above or below allowed values
	perPage = util.MaxInt(1, util.MinInt(100, perPage))

	// send API call to capture the list of secrets
	s, err := secret.FromContext(c, e).List(t, o, n, page, perPage)
	if err != nil {
		retErr := fmt.Errorf("unable to get secrets for %s/%s/%s from %s service: %w", t, o, n, e, err)

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

// swagger:operation GET /api/v1/secrets/{engine}/{type}/{org}/{name}/{secret} secrets GetSecret
//
// Retrieve a secret from the configured backend.
//
// ---
// x-success_http_code: '200'
// produces:
// - application/json
// parameters:
// - in: path
//   name: engine
//   description: Secret engine to create a secret in
//   required: true
//   type: string
// - in: path
//   name: type
//   description: Secret type to create. Options 'org', 'repo', or 'shared'
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
//     type: json
//     schema:
//       "$ref": "#/definitions/Secret"
//   '500':
//     description: Unable to retrieve the secret
//     schema:
//       type: string

// GetSecret gets a secret from the provided secrets service.
func GetSecret(c *gin.Context) {
	// capture middleware values
	u := user.Retrieve(c)
	e := c.Param("engine")
	t := c.Param("type")
	o := c.Param("org")
	n := c.Param("name")
	s := strings.TrimPrefix(c.Param("secret"), "/")

	logrus.Infof("Reading secret %s/%s/%s/%s from %s service", t, o, n, s, e)

	// send API call to capture the secret
	secret, err := secret.FromContext(c, e).Get(t, o, n, s)
	if err != nil {
		retErr := fmt.Errorf("unable to get secret %s/%s/%s/%s from %s service: %w", t, o, n, s, e, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// only allow agents to access the full secret with the value
	if u.GetAdmin() && u.GetName() == "vela-worker" {
		c.JSON(http.StatusOK, secret)

		return
	}

	c.JSON(http.StatusOK, secret.Sanitize())
}

// swagger:operation PUT /api/v1/secrets/{engine}/{type}/{org}/{name}/{secret} secrets UpdateSecrets
//
// Update a secret from the configured backend.
//
// ---
// x-success_http_code: '200'
// produces:
// - application/json
// parameters:
// - in: path
//   name: engine
//   description: Secret engine to create a secret in
//   required: true
//   type: string
// - in: path
//   name: type
//   description: Secret type to create. Options 'org', 'repo', or 'shared'
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
//     type: json
//     schema:
//       "$ref": "#/definitions/Secret"
//   '400':
//     description: Unable to update the secret
//     schema:
//       type: string
//   '500':
//     description: Unable to update the secret
//     schema:
//       type: string

// UpdateSecret updates a secret for the provided secrets service.
func UpdateSecret(c *gin.Context) {
	// capture middleware values
	e := c.Param("engine")
	t := c.Param("type")
	o := c.Param("org")
	n := c.Param("name")
	s := strings.TrimPrefix(c.Param("secret"), "/")

	logrus.Infof("Updating secret %s/%s/%s/%s for %s service", t, o, n, s, e)

	// capture body from API request
	input := new(library.Secret)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for secret %s/%s/%s/%s for %s service: %v", t, o, n, s, e, err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// update secret fields if provided
	input.SetName(s)
	input.SetOrg(o)
	input.SetRepo(n)
	input.SetType(t)

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
		retErr := fmt.Errorf("unable to update secret %s/%s/%s/%s for %s service: %w", t, o, n, s, e, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// send API call to capture the updated secret
	secret, _ := secret.FromContext(c, e).Get(t, o, n, input.GetName())

	c.JSON(http.StatusOK, secret.Sanitize())
}

// swagger:operation DELETE /api/v1/secrets/{engine}/{type}/{org}/{name}/{secret} secrets DeleteSecret
//
// Delete a secret from the configured backend.
//
// ---
// x-success_http_code: '200'
// produces:
// - application/json
// parameters:
// - in: path
//   name: engine
//   description: Secret engine to create a secret in
//   required: true
//   type: string
// - in: path
//   name: type
//   description: Secret type to create. Options 'org', 'repo', or 'shared'
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
//   description: Payload containing secret to update
//   required: true
//   schema:
//     "$ref": "#/definitions/Secret"
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully updated the secret
//     schema:
//       type: string
//   '500':
//     description: Unable to update the secret
//     schema:
//       type: string

// DeleteSecret deletes a secret from the provided secrets service.
func DeleteSecret(c *gin.Context) {
	// capture middleware values
	e := c.Param("engine")
	t := c.Param("type")
	o := c.Param("org")
	n := c.Param("name")
	s := strings.TrimPrefix(c.Param("secret"), "/")

	logrus.Infof("Deleting secret %s/%s/%s/%s from %s service", t, o, n, s, e)

	// send API call to remove the secret
	err := secret.FromContext(c, e).Delete(t, o, n, s)
	if err != nil {
		retErr := fmt.Errorf("unable to delete secret %s/%s/%s/%s from %s service: %w", t, o, n, s, e, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("Secret %s/%s/%s/%s deleted from %s service", t, o, n, s, e))
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
