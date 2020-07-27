// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/token"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/source"
	"github.com/go-vela/server/util"

	"github.com/go-vela/types/library"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// swagger:operation POST /api/v1/users users CreateUser
//
// Create a user for the configured backend
//
// ---
// x-success_http_code: '201'
// produces:
// - application/json
// parameters:
// - in: body
//   name: body
//   description: Payload containing the user to create
//   required: true
//   schema:
//     "$ref": "#/definitions/User"
// - in: header
//   name: Authorization
//   description: Vela bearer token
//   required: true
//   type: string
// responses:
//   '201':
//     description: Successfully created the user
//     type: json
//     schema:
//       "$ref": "#/definitions/User"
//   '400':
//     description: Unable to create the user
//     schema:
//       type: string
//   '500':
//     description: Unable to create the user
//     schema:
//       type: string

// CreateUser represents the API handler to create
// a user in the configured backend.
func CreateUser(c *gin.Context) {
	logrus.Info("Creating new user")

	// capture body from API request
	input := new(library.User)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for new user: %w", err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// send API call to create the user
	err = database.FromContext(c).CreateUser(input)
	if err != nil {
		retErr := fmt.Errorf("unable to create user: %w", err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// send API call to capture the created user
	u, _ := database.FromContext(c).GetUserName(input.GetName())

	c.JSON(http.StatusCreated, u)
}

// swagger:operation GET /api/v1/users users GetUsers
//
// Retrieve a list of users for the configured backend
//
// ---
// x-success_http_code: '200'
// produces:
// - application/json
// parameters:
// - in: header
//   name: Authorization
//   description: Vela bearer token
//   required: true
//   type: string
// responses:
//   '200':
//     description: Successfully retrieved the list of users
//     type: json
//     schema:
//       "$ref": "#/definitions/User"
//   '400':
//     description: Unable to retrieve the list of users
//     schema:
//       type: string
//   '500':
//     description: Unable to retrieve the list of users
//     schema:
//       type: string

// GetUsers represents the API handler to capture a list
// of users from the configured backend.
func GetUsers(c *gin.Context) {
	logrus.Info("Reading lite users")

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

	// send API call to capture the total number of users
	t, err := database.FromContext(c).GetUserCount()
	if err != nil {
		retErr := fmt.Errorf("unable to get users count: %w", err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// send API call to capture the list of users
	u, err := database.FromContext(c).GetUserLiteList(page, perPage)
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

	c.JSON(http.StatusOK, u)
}

// swagger:operation GET /api/v1/user users GetCurrentUser
//
// Retrieve the current authenticated user from the configured backend
//
// ---
// x-success_http_code: '200'
// produces:
// - application/json
// parameters:
// - in: header
//   name: Authorization
//   description: Vela bearer token
//   required: true
//   type: string
// responses:
//   '200':
//     description: Successfully retrieved the current user
//     type: json
//     schema:
//       "$ref": "#/definitions/User"

// GetCurrentUser represents the API handler to capture the
// currently authenticated user from the configured backend.
func GetCurrentUser(c *gin.Context) {
	logrus.Infof("Reading current user")

	// retrieve user from context
	u := user.Retrieve(c)

	c.JSON(http.StatusOK, u)
}

// swagger:operation PUT /api/v1/user users UpdateCurrentUser
//
// Update the current authenticated user in the configured backend
//
// ---
// x-success_http_code: '200'
// produces:
// - application/json
// parameters:
// - in: body
//   name: body
//   description: Payload containing the user to update
//   required: true
//   schema:
//     "$ref": "#/definitions/User"
// - in: header
//   name: Authorization
//   description: Vela bearer token
//   required: true
//   type: string
// responses:
//   '200':
//     description: Successfully updated the current user
//     type: json
//     schema:
//       "$ref": "#/definitions/User"
//   '400':
//     description: Unable to update the current user
//     schema:
//       type: string
//   '404':
//     description: Unable to update the current user
//     schema:
//       type: string
//   '500':
//     description: Unable to update the current user
//     schema:
//       type: string

// UpdateCurrentUser represents the API handler to capture and
// update the currently authenticated user from the configured backend.
func UpdateCurrentUser(c *gin.Context) {
	// retrieve user from context
	user := user.Retrieve(c)

	logrus.Infof("Updating current user %s", user)

	// capture body from API request
	input := new(library.User)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for user %s: %w", user, err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// update user fields if provided
	if input.Favorites != nil {
		// update favorites if set
		user.SetFavorites(input.GetFavorites())
	}

	// send API call to update the user
	err = database.FromContext(c).UpdateUser(user)
	if err != nil {
		retErr := fmt.Errorf("unable to update user %s: %w", user, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// send API call to capture the updated user
	user, err = database.FromContext(c).GetUserName(user.GetName())
	if err != nil {
		retErr := fmt.Errorf("unable to get updated user %s: %w", user, err)

		util.HandleError(c, http.StatusNotFound, retErr)

		return
	}

	c.JSON(http.StatusOK, user)
}

// swagger:operation GET /api/v1/users/{user} users GetUser
//
// Retrieve a user for the configured backend
//
// ---
// x-success_http_code: '200'
// produces:
// - application/json
// parameters:
// - in: path
//   name: user
//   description: Name of the user
//   required: true
//   type: string
// - in: header
//   name: Authorization
//   description: Vela bearer token
//   required: true
//   type: string
// responses:
//   '200':
//     description: Successfully retrieved the user
//     type: json
//     schema:
//       "$ref": "#/definitions/User"
//   '400':
//     description: Unable to retrieve the user
//     schema:
//       type: string

// GetUser represents the API handler to capture a
// user from the configured backend.
func GetUser(c *gin.Context) {
	// capture middleware values
	user := c.Param("user")

	logrus.Infof("Reading user %s", user)

	// send API call to capture the user
	u, err := database.FromContext(c).GetUserName(user)
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
// x-success_http_code: '200'
// produces:
// - application/json
// parameters:
// - in: header
//   name: Authorization
//   description: Vela bearer token
//   required: true
//   type: string
// responses:
//   '200':
//     description: Successfully retrieved a list of repos for the current user
//     type: json
//     schema:
//       "$ref": "#/definitions/Repo"
//   '404':
//     description: Unable to retrieve a list of repos for the current user
//     schema:
//       type: string

// GetUserSourceRepos represents the API handler to capture
// the list of repos for a user from the configured backend.
func GetUserSourceRepos(c *gin.Context) {
	// capture middleware values
	u := user.Retrieve(c)

	logrus.Infof("Getting list of available repos for user %s", u.GetName())

	// variables to capture requested data
	srcRepos := []*library.Repo{}
	dbRepos := []*library.Repo{}
	output := make(map[string][]library.Repo)

	// send API call to capture the list of repos for the user
	srcRepos, err := source.FromContext(c).ListUserRepos(u)
	if err != nil {
		retErr := fmt.Errorf("unable to get source repos for user %s: %w", u.GetName(), err)

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
		for page > 0 {
			// send API call to capture the list of repos for the org
			dbReposPart, err := database.FromContext(c).GetOrgRepoList(org, page, 100)
			if err != nil {
				retErr := fmt.Errorf("unable to get repos for org %s: %w", org, err)

				util.HandleError(c, http.StatusNotFound, retErr)

				return
			}

			// add repos to list of database org repos
			dbRepos = append(dbRepos, dbReposPart...)

			// making an assumption that 50 means there is another page
			if len(dbReposPart) == 50 {
				page++
			} else {
				page = 0
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
// x-success_http_code: '200'
// produces:
// - application/json
// parameters:
// - in: body
//   name: body
//   description: Payload containing the user to update
//   required: true
//   schema:
//     "$ref": "#/definitions/User"
// - in: path
//   name: user
//   description: Name of the user
//   required: true
//   type: string
// - in: header
//   name: Authorization
//   description: Vela bearer token
//   required: true
//   type: string
// responses:
//   '200':
//     description: Successfully updated the user
//     type: json
//     schema:
//       "$ref": "#/definitions/User"
//   '400':
//     description: Unable to update the user
//     schema:
//       type: string
//   '404':
//     description: Unable to update the user
//     schema:
//       type: string
//   '500':
//     description: Unable to update the user
//     schema:
//       type: string

// UpdateUser represents the API handler to update
// a user in the configured backend.
func UpdateUser(c *gin.Context) {
	// capture middleware values
	user := c.Param("user")

	logrus.Infof("Updating user %s", user)

	// capture body from API request
	input := new(library.User)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for user %s: %w", user, err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// send API call to capture the user
	u, err := database.FromContext(c).GetUserName(user)
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
	u, _ = database.FromContext(c).GetUserName(user)

	c.JSON(http.StatusOK, u)
}

// swagger:operation DELETE /api/v1/users/{user} users DeleteUser
//
// Delete a user for the configured backend
//
// ---
// x-success_http_code: '200'
// produces:
// - application/json
// parameters:
// - in: path
//   name: user
//   description: Name of the user
//   required: true
//   type: string
// - in: header
//   name: Authorization
//   description: Vela bearer token
//   required: true
//   type: string
// responses:
//   '200':
//     description: Successfully deleted of user
//     schema:
//       type: string
//   '404':
//     description: Unable to delete user
//     schema:
//       type: string
//   '500':
//     description: Unable to delete user
//     schema:
//       type: string

// DeleteUser represents the API handler to remove
// a user from the configured backend.
func DeleteUser(c *gin.Context) {
	// capture middleware values
	user := c.Param("user")

	logrus.Infof("Deleting user %s", user)

	// send API call to capture the user
	u, err := database.FromContext(c).GetUserName(user)
	if err != nil {
		retErr := fmt.Errorf("unable to get user %s: %w", user, err)

		util.HandleError(c, http.StatusNotFound, retErr)

		return
	}

	// send API call to remove the user
	err = database.FromContext(c).DeleteUser(u.GetID())
	if err != nil {
		retErr := fmt.Errorf("unable to delete user %s: %w", u.GetName(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("User %s deleted", u.GetName()))
}

// swagger:operation POST /api/v1/user/token users CreateToken
//
// Create a token for the current authenticated user
//
// ---
// x-success_http_code: '201'
// produces:
// - application/json
// parameters:
// - in: header
//   name: Authorization
//   description: Vela bearer token
//   required: true
//   type: string
// responses:
//   '201':
//     description: Successfully created a token for the current user
//     type: json
//     schema:
//       "$ref": "#/definitions/Login"
//   '500':
//     description: Unable to create a token for the current user
//     schema:
//       type: string

// CreateToken represents the API handler to create
// a user token in the configured backend.
func CreateToken(c *gin.Context) {
	// capture middleware values
	u := user.Retrieve(c)

	logrus.Infof("Composing token for user %s", u.GetName())

	// compose JWT token for user
	rt, at, err := token.Compose(c, u)
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

	c.JSON(http.StatusOK, library.Login{Token: &at})
}

// swagger:operation DELETE /api/v1/user/token users DeleteToken
//
// Delete a token for the current authenticated user
//
// ---
// x-success_http_code: '200'
// produces:
// - application/json
// parameters:
// - in: header
//   name: Authorization
//   description: Vela bearer token
//   required: true
//   type: string
// responses:
//   '200':
//     description: Successfully delete a token for the current user
//     schema:
//       type: string
//   '500':
//     description: Unable to delete a token for the current user
//     schema:
//       type: string

// DeleteToken represents the API handler to revoke
// and recreate a user token in the configured backend.
func DeleteToken(c *gin.Context) {
	// capture middleware values
	u := user.Retrieve(c)

	logrus.Infof("Revoking token for user %s", u.GetName())

	// create unique id for the user
	uid, err := uuid.NewRandom()
	if err != nil {
		retErr := fmt.Errorf("unable to create UID for user %s: %w", u.GetName(), err)

		util.HandleError(c, http.StatusServiceUnavailable, retErr)

		return
	}

	u.SetHash(
		base64.StdEncoding.EncodeToString(
			[]byte(uid.String()),
		),
	)

	// compose JWT token for user
	rt, at, err := token.Compose(c, u)
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

	c.JSON(http.StatusOK, library.Login{Token: &at})
}
