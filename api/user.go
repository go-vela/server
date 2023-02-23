// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/internal/token"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// swagger:operation POST /api/v1/users users CreateUser
//
// Create a user for the configured backend
//
// ---
// produces:
// - application/json
// parameters:
// - in: body
//   name: body
//   description: Payload containing the user to create
//   required: true
//   schema:
//     "$ref": "#/definitions/User"
// security:
//   - ApiKeyAuth: []
// responses:
//   '201':
//     description: Successfully created the user
//     schema:
//       "$ref": "#/definitions/User"
//   '400':
//     description: Unable to create the user
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to create the user
//     schema:
//       "$ref": "#/definitions/Error"

// CreateUser represents the API handler to create
// a user in the configured backend.
func CreateUser(c *gin.Context) {
	// capture middleware values
	u := user.Retrieve(c)

	// capture body from API request
	input := new(library.User)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for new user: %w", err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"user": u.GetName(),
	}).Infof("creating new user %s", input.GetName())

	// send API call to create the user
	err = database.FromContext(c).CreateUser(input)
	if err != nil {
		retErr := fmt.Errorf("unable to create user: %w", err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// send API call to capture the created user
	user, _ := database.FromContext(c).GetUserForName(input.GetName())

	c.JSON(http.StatusCreated, user)
}

// swagger:operation GET /api/v1/users users GetUsers
//
// Retrieve a list of users for the configured backend
//
// ---
// produces:
// - application/json
// security:
//   - ApiKeyAuth: []
// parameters:
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
// responses:
//   '200':
//     description: Successfully retrieved the list of users
//     schema:
//       type: array
//       items:
//         "$ref": "#/definitions/User"
//     headers:
//       X-Total-Count:
//         description: Total number of results
//         type: integer
//       Link:
//         description: see https://tools.ietf.org/html/rfc5988
//         type: string
//   '400':
//     description: Unable to retrieve the list of users
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to retrieve the list of users
//     schema:
//       "$ref": "#/definitions/Error"

// GetUsers represents the API handler to capture a list
// of users from the configured backend.
func GetUsers(c *gin.Context) {
	// capture middleware values
	u := user.Retrieve(c)

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"user": u.GetName(),
	}).Info("reading lite users")

	// capture page query parameter if present
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		retErr := fmt.Errorf("unable to convert page query parameter for users: %w", err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// capture per_page query parameter if present
	perPage, err := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	if err != nil {
		retErr := fmt.Errorf("unable to convert per_page query parameter for users: %w", err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// ensure per_page isn't above or below allowed values
	perPage = util.MaxInt(1, util.MinInt(100, perPage))

	// send API call to capture the list of users
	users, t, err := database.FromContext(c).ListLiteUsers(page, perPage)
	if err != nil {
		retErr := fmt.Errorf("unable to get users: %w", err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// create pagination object
	pagination := Pagination{
		Page:    page,
		PerPage: perPage,
		Total:   t,
	}
	// set pagination headers
	pagination.SetHeaderLink(c)

	c.JSON(http.StatusOK, users)
}

// swagger:operation GET /api/v1/user users GetCurrentUser
//
// Retrieve the current authenticated user from the configured backend
//
// ---
// produces:
// - application/json
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved the current user
//     schema:
//       "$ref": "#/definitions/User"

// GetCurrentUser represents the API handler to capture the
// currently authenticated user from the configured backend.
func GetCurrentUser(c *gin.Context) {
	// capture middleware values
	u := user.Retrieve(c)

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"user": u.GetName(),
	}).Infof("reading current user %s", u.GetName())

	c.JSON(http.StatusOK, u)
}

// swagger:operation PUT /api/v1/user users UpdateCurrentUser
//
// Update the current authenticated user in the configured backend
//
// ---
// produces:
// - application/json
// parameters:
// - in: body
//   name: body
//   description: Payload containing the user to update
//   required: true
//   schema:
//     "$ref": "#/definitions/User"
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully updated the current user
//     schema:
//       "$ref": "#/definitions/User"
//   '400':
//     description: Unable to update the current user
//     schema:
//       "$ref": "#/definitions/Error"
//   '404':
//     description: Unable to update the current user
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to update the current user
//     schema:
//       "$ref": "#/definitions/Error"

// UpdateCurrentUser represents the API handler to capture and
// update the currently authenticated user from the configured backend.
func UpdateCurrentUser(c *gin.Context) {
	// capture middleware values
	u := user.Retrieve(c)

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"user": u.GetName(),
	}).Infof("updating current user %s", u.GetName())

	// capture body from API request
	input := new(library.User)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for user %s: %w", u.GetName(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// update user fields if provided
	if input.Favorites != nil {
		// update favorites if set
		u.SetFavorites(input.GetFavorites())
	}

	// send API call to update the user
	err = database.FromContext(c).UpdateUser(u)
	if err != nil {
		retErr := fmt.Errorf("unable to update user %s: %w", u.GetName(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// send API call to capture the updated user
	u, err = database.FromContext(c).GetUserForName(u.GetName())
	if err != nil {
		retErr := fmt.Errorf("unable to get updated user %s: %w", u.GetName(), err)

		util.HandleError(c, http.StatusNotFound, retErr)

		return
	}

	c.JSON(http.StatusOK, u)
}

// swagger:operation GET /api/v1/users/{user} users GetUser
//
// Retrieve a user for the configured backend
//
// ---
// produces:
// - application/json
// parameters:
// - in: path
//   name: user
//   description: Name of the user
//   required: true
//   type: string
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved the user
//     schema:
//       "$ref": "#/definitions/User"
//   '404':
//     description: Unable to retrieve the user
//     schema:
//       "$ref": "#/definitions/Error"

// GetUser represents the API handler to capture a
// user from the configured backend.
func GetUser(c *gin.Context) {
	// capture middleware values
	u := user.Retrieve(c)
	user := util.PathParameter(c, "user")

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"user": u.GetName(),
	}).Infof("reading user %s", user)

	// send API call to capture the user
	u, err := database.FromContext(c).GetUserForName(user)
	if err != nil {
		retErr := fmt.Errorf("unable to get user %s: %w", user, err)

		util.HandleError(c, http.StatusNotFound, retErr)

		return
	}

	c.JSON(http.StatusOK, u)
}

// swagger:operation GET /api/v1/user/source/repos users GetUserSourceRepos
//
// Retrieve a list of repos for the current authenticated user
//
// ---
// produces:
// - application/json
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved a list of repos for the current user
//     schema:
//       "$ref": "#/definitions/Repo"
//   '404':
//     description: Unable to retrieve a list of repos for the current user
//     schema:
//       "$ref": "#/definitions/Error"

// GetUserSourceRepos represents the API handler to capture
// the list of repos for a user from the configured backend.
func GetUserSourceRepos(c *gin.Context) {
	// capture middleware values
	u := user.Retrieve(c)

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"user": u.GetName(),
	}).Infof("reading available SCM repos for user %s", u.GetName())

	// variables to capture requested data
	dbRepos := []*library.Repo{}
	output := make(map[string][]library.Repo)

	// send API call to capture the list of repos for the user
	srcRepos, err := scm.FromContext(c).ListUserRepos(u)
	if err != nil {
		retErr := fmt.Errorf("unable to get SCM repos for user %s: %w", u.GetName(), err)

		util.HandleError(c, http.StatusNotFound, retErr)

		return
	}

	// create a map
	// TODO: clean this up
	for _, srepo := range srcRepos {
		// local variables to avoid bad memory address de-referencing
		// initialize active to false
		org := srepo.Org
		name := srepo.Name
		active := false

		// library struct to omit optional fields
		repo := library.Repo{
			Org:    org,
			Name:   name,
			Active: &active,
		}
		output[srepo.GetOrg()] = append(output[srepo.GetOrg()], repo)
	}

	for org := range output {
		// capture source repos from the database backend, grouped by org
		page := 1
		filters := map[string]interface{}{}

		for page > 0 {
			// send API call to capture the list of repos for the org
			dbReposPart, _, err := database.FromContext(c).ListReposForOrg(org, "name", filters, page, 100)
			if err != nil {
				retErr := fmt.Errorf("unable to get repos for org %s: %w", org, err)

				util.HandleError(c, http.StatusNotFound, retErr)

				return
			}

			// add repos to list of database org repos
			dbRepos = append(dbRepos, dbReposPart...)

			// assume no more pages exist if under 100 results are returned
			if len(dbReposPart) < 100 {
				page = 0
			} else {
				page++
			}
		}

		// apply org repos active status to output map
		for _, dbRepo := range dbRepos {
			if orgRepos, ok := output[dbRepo.GetOrg()]; ok {
				for i := range orgRepos {
					if orgRepos[i].GetName() == dbRepo.GetName() {
						active := dbRepo.GetActive()
						(&orgRepos[i]).Active = &active
					}
				}
			}
		}
	}

	c.JSON(http.StatusOK, output)
}

// swagger:operation PUT /api/v1/users/{user} users UpdateUser
//
// Update a user for the configured backend
//
// ---
// produces:
// - application/json
// parameters:
// - in: path
//   name: user
//   description: Name of the user
//   required: true
//   type: string
// - in: body
//   name: body
//   description: Payload containing the user to update
//   required: true
//   schema:
//     "$ref": "#/definitions/User"
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully updated the user
//     schema:
//       "$ref": "#/definitions/User"
//   '400':
//     description: Unable to update the user
//     schema:
//       "$ref": "#/definitions/Error"
//   '404':
//     description: Unable to update the user
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to update the user
//     schema:
//       "$ref": "#/definitions/Error"

// UpdateUser represents the API handler to update
// a user in the configured backend.
func UpdateUser(c *gin.Context) {
	// capture middleware values
	u := user.Retrieve(c)
	user := util.PathParameter(c, "user")

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"user": u.GetName(),
	}).Infof("updating user %s", user)

	// capture body from API request
	input := new(library.User)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for user %s: %w", user, err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// send API call to capture the user
	u, err = database.FromContext(c).GetUserForName(user)
	if err != nil {
		retErr := fmt.Errorf("unable to get user %s: %w", user, err)

		util.HandleError(c, http.StatusNotFound, retErr)

		return
	}

	// update user fields if provided
	if input.GetActive() {
		// update active if set to true
		u.SetActive(input.GetActive())
	}

	if input.GetAdmin() {
		// update admin if set to true
		u.SetAdmin(input.GetAdmin())
	}

	if input.Favorites != nil {
		// update favorites if set
		u.SetFavorites(input.GetFavorites())
	}

	// send API call to update the user
	err = database.FromContext(c).UpdateUser(u)
	if err != nil {
		retErr := fmt.Errorf("unable to update user %s: %w", user, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// send API call to capture the updated user
	u, _ = database.FromContext(c).GetUserForName(user)

	c.JSON(http.StatusOK, u)
}

// swagger:operation DELETE /api/v1/users/{user} users DeleteUser
//
// Delete a user for the configured backend
//
// ---
// produces:
// - application/json
// parameters:
// - in: path
//   name: user
//   description: Name of the user
//   required: true
//   type: string
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully deleted of user
//     schema:
//       type: string
//   '404':
//     description: Unable to delete user
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to delete user
//     schema:
//       "$ref": "#/definitions/Error"

// DeleteUser represents the API handler to remove
// a user from the configured backend.
func DeleteUser(c *gin.Context) {
	// capture middleware values
	u := user.Retrieve(c)
	user := util.PathParameter(c, "user")

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"user": u.GetName(),
	}).Infof("deleting user %s", user)

	// send API call to capture the user
	u, err := database.FromContext(c).GetUserForName(user)
	if err != nil {
		retErr := fmt.Errorf("unable to get user %s: %w", user, err)

		util.HandleError(c, http.StatusNotFound, retErr)

		return
	}

	// send API call to remove the user
	err = database.FromContext(c).DeleteUser(u)
	if err != nil {
		retErr := fmt.Errorf("unable to delete user %s: %w", u.GetName(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("user %s deleted", u.GetName()))
}

// swagger:operation POST /api/v1/user/token users CreateToken
//
// Create a token for the current authenticated user
//
// ---
// produces:
// - application/json
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully created a token for the current user
//     schema:
//       "$ref": "#/definitions/Login"
//   '503':
//     description: Unable to create a token for the current user
//     schema:
//       "$ref": "#/definitions/Error"

// CreateToken represents the API handler to create
// a user token in the configured backend.
//
//nolint:dupl // ignore duplicate flag with delete token
func CreateToken(c *gin.Context) {
	// capture middleware values
	u := user.Retrieve(c)

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"user": u.GetName(),
	}).Infof("composing token for user %s", u.GetName())

	tm := c.MustGet("token-manager").(*token.Manager)

	// compose JWT token for user
	rt, at, err := tm.Compose(c, u)
	if err != nil {
		retErr := fmt.Errorf("unable to compose token for user %s: %w", u.GetName(), err)

		util.HandleError(c, http.StatusServiceUnavailable, retErr)

		return
	}

	u.SetRefreshToken(rt)

	// send API call to update the user
	err = database.FromContext(c).UpdateUser(u)
	if err != nil {
		retErr := fmt.Errorf("unable to update user %s: %w", u.GetName(), err)

		util.HandleError(c, http.StatusServiceUnavailable, retErr)

		return
	}

	c.JSON(http.StatusOK, library.Token{Token: &at})
}

// swagger:operation DELETE /api/v1/user/token users DeleteToken
//
// Delete a token for the current authenticated user
//
// ---
// produces:
// - application/json
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully delete a token for the current user
//     schema:
//       type: string
//   '500':
//     description: Unable to delete a token for the current user
//     schema:
//       "$ref": "#/definitions/Error"

// DeleteToken represents the API handler to revoke
// and recreate a user token in the configured backend.
//
//nolint:dupl // ignore duplicate flag with create token
func DeleteToken(c *gin.Context) {
	// capture middleware values
	u := user.Retrieve(c)

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"user": u.GetName(),
	}).Infof("revoking token for user %s", u.GetName())

	tm := c.MustGet("token-manager").(*token.Manager)

	// compose JWT token for user
	rt, at, err := tm.Compose(c, u)
	if err != nil {
		retErr := fmt.Errorf("unable to compose token for user %s: %w", u.GetName(), err)

		util.HandleError(c, http.StatusServiceUnavailable, retErr)

		return
	}

	u.SetRefreshToken(rt)

	// send API call to update the user
	err = database.FromContext(c).UpdateUser(u)
	if err != nil {
		retErr := fmt.Errorf("unable to update user %s: %w", u.GetName(), err)

		util.HandleError(c, http.StatusServiceUnavailable, retErr)

		return
	}

	c.JSON(http.StatusOK, library.Token{Token: &at})
}
