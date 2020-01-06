// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"

	"golang.org/x/sync/errgroup"

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

// GetCurrentUser represents the API handler to capture the
// currently authenticated user from the configured backend.
func GetCurrentUser(c *gin.Context) {
	logrus.Infof("Reading current user")

	// retrieve user from context
	u := user.Retrieve(c)

	c.JSON(http.StatusOK, u)
}

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

// GetUserSourceRepos represents the API handler to capture
// the list of repos for a user from the configured backend.
func GetUserSourceRepos(c *gin.Context) {
	// capture middleware values
	u := user.Retrieve(c)

	logrus.Infof("Getting list of available repos for user %s", u.GetName())

	// variables to capture requested data
	srcRepos := []*library.Repo{}
	output := make(map[string][]library.Repo)
	threads := new(errgroup.Group)

	// capture user's repos from the source backend
	threads.Go(func() error {
		// send API call to capture the list of repos for the user
		r, err := source.FromContext(c).ListUserRepos(u)
		if err != nil {
			return fmt.Errorf("unable to get source repos for user %s: %w", u.GetName(), err)
		}

		// add repos to list of source repos
		srcRepos = append(srcRepos, r...)
		return nil
	})

	// wait for all threads to complete
	err := threads.Wait()
	if err != nil {
		util.HandleError(c, http.StatusInternalServerError, err)

		return
	}

	// create a map and if the repo exists in Vela attach the active status
	// TODO: clean this up
	for _, srepo := range srcRepos {
		// local variables to avoid bad memory address de-referencing
		org := srepo.Org
		name := srepo.Name

		// send API call to capture the source repo from the database, if it exists
		// if record does not exist, repo is considered inactive
		dbRepo, err := database.FromContext(c).GetRepo(srepo.GetOrg(), srepo.GetName())
		if err != nil && err.Error() != "record not found" {
			util.HandleError(c, http.StatusInternalServerError, err)

			return
		}

		active := dbRepo.GetActive()

		// library struct to omit optional fields
		repo := library.Repo{
			Org:    org,
			Name:   name,
			Active: &active,
		}
		output[srepo.GetOrg()] = append(output[srepo.GetOrg()], repo)
	}

	c.JSON(http.StatusOK, output)
}

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
		// update active if set
		u.SetActive(input.GetActive())
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

// CreateToken represents the API handler to create
// a user token in the configured backend.
func CreateToken(c *gin.Context) {
	// capture middleware values
	u := user.Retrieve(c)

	logrus.Infof("Composing token for user %s", u.GetName())

	// compose JWT token for user
	t, err := token.Compose(u)
	if err != nil {
		retErr := fmt.Errorf("unable to compose token for user %s: %w", u.GetName(), err)

		util.HandleError(c, http.StatusServiceUnavailable, retErr)

		return
	}

	c.JSON(http.StatusOK, library.Login{Username: u.Name, Token: &t})
}

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

	// send API call to update the user
	err = database.FromContext(c).UpdateUser(u)
	if err != nil {
		retErr := fmt.Errorf("unable to update user %s: %w", u.GetName(), err)

		util.HandleError(c, http.StatusServiceUnavailable, retErr)

		return
	}

	// compose JWT token for user
	t, err := token.Compose(u)
	if err != nil {
		retErr := fmt.Errorf("unable to compose token for user %s: %w", u.GetName(), err)

		util.HandleError(c, http.StatusServiceUnavailable, retErr)

		return
	}

	c.JSON(http.StatusOK, library.Login{Username: u.Name, Token: &t})
}
