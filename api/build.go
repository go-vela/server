// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-vela/compiler/compiler"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/queue"
	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/source"
	"github.com/go-vela/server/util"

	"github.com/go-vela/types"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// swagger:operation POST /api/v1/repos/{org}/{repo}/builds builds CreateBuild
//
// Create a build in the configured backend
//
// ---
// x-success_http_code: '201'
// produces:
// - application/json
// parameters:
// - in: body
//   name: body
//   description: Payload containing the build to update
//   required: true
//   schema:
//     "$ref": "#/definitions/Build"
// - in: path
//   name: repo
//   description: Name of the repo
//   required: true
//   type: string
// - in: path
//   name: org
//   description: Name of the org
//   required: true
//   type: string
// - in: header
//   name: Authorization
//   description: Vela bearer token
//   required: true
//   type: string
// responses:
//   '201':
//     description: Successfully created the build
//     type: json
//     schema:
//       "$ref": "#/definitions/Build"
//   '400':
//     description: Unable to create the build
//     schema:
//       type: string
//   '404':
//     description: Unable to create the build
//     schema:
//       type: string
//   '500':
//     description: Unable to create the build
//     schema:
//       type: string

// CreateBuild represents the API handler to
// create a build in the configured backend.
func CreateBuild(c *gin.Context) {
	// capture middleware values
	m := c.MustGet("metadata").(*types.Metadata)
	r := repo.Retrieve(c)

	logrus.Infof("Creating new build for repo %s", r.GetFullName())

	// capture body from API request
	input := new(library.Build)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for new build for repo %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// verify the build has a valid event and the repo allows that event type
	if (input.GetEvent() == constants.EventPush && !r.GetAllowPush()) ||
		(input.GetEvent() == constants.EventPull && !r.GetAllowPull()) ||
		(input.GetEvent() == constants.EventTag && !r.GetAllowTag()) ||
		(input.GetEvent() == constants.EventDeploy && !r.GetAllowDeploy()) {
		retErr := fmt.Errorf("unable to create new build: %s does not have %s events enabled", r.GetFullName(), input.GetEvent())

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// send API call to capture the repo owner
	u, err := database.FromContext(c).GetUser(r.GetUserID())
	if err != nil {
		retErr := fmt.Errorf("unable to get owner for %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// send API call to capture the last build for the repo
	lastBuild, err := database.FromContext(c).GetLastBuild(r)
	if err != nil {
		retErr := fmt.Errorf("unable to get last build for %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// update fields in build object
	input.SetRepoID(r.GetID())
	input.SetStatus(constants.StatusPending)
	input.SetCreated(time.Now().UTC().Unix())
	input.SetNumber(1)
	input.SetParent(input.GetNumber())

	if lastBuild != nil {
		input.SetNumber(
			lastBuild.GetNumber() + 1,
		)
		input.SetParent(lastBuild.GetNumber())
	}

	// populate the build link if a web address is provided
	if len(m.Vela.WebAddress) > 0 {
		input.SetLink(
			fmt.Sprintf("%s/%s/%d", m.Vela.WebAddress, r.GetFullName(), input.GetNumber()),
		)
	}

	// variable to store changeset files
	var files []string
	// check if the build event is not pull_request
	if !strings.EqualFold(input.GetEvent(), constants.EventPull) {
		// send API call to capture list of files changed for the commit
		files, err = source.FromContext(c).Changeset(u, r, input.GetCommit())
		if err != nil {
			retErr := fmt.Errorf("unable to process webhook: failed to get changeset for %s: %w", r.GetFullName(), err)

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}
	}

	// files is empty if the build event is pull_request
	if len(files) == 0 {
		// parse out pull request number from base ref
		//
		// pattern: refs/pull/1/head
		var parts []string
		if strings.HasPrefix(input.GetRef(), "refs/pull/") {
			parts = strings.Split(input.GetRef(), "/")
		}

		// capture number by converting from string
		number, err := strconv.Atoi(parts[2])
		if err != nil {
			// capture number by scanning from string
			_, err := fmt.Sscanf(input.GetRef(), "%s/%s/%d/%s", nil, nil, &number, nil)
			if err != nil {
				retErr := fmt.Errorf("unable to process webhook: failed to get pull_request number for %s: %w", r.GetFullName(), err)

				util.HandleError(c, http.StatusInternalServerError, retErr)

				return
			}
		}

		// send API call to capture list of files changed for the pull request
		files, err = source.FromContext(c).ChangesetPR(u, r, number)
		if err != nil {
			retErr := fmt.Errorf("unable to process webhook: failed to get changeset for %s: %w", r.GetFullName(), err)

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}
	}

	// send API call to capture the pipeline configuration file
	config, err := source.FromContext(c).ConfigBackoff(u, r.GetOrg(), r.GetName(), input.GetCommit())
	if err != nil {
		retErr := fmt.Errorf("unable to get pipeline configuration for %s/%d: %w", r.GetFullName(), input.GetNumber(), err)

		util.HandleError(c, http.StatusNotFound, retErr)

		return
	}

	// parse and compile the pipeline configuration file
	p, err := compiler.FromContext(c).
		WithBuild(input).
		WithFiles(files).
		WithMetadata(m).
		WithRepo(r).
		WithUser(u).
		Compile(config)
	if err != nil {
		retErr := fmt.Errorf("unable to compile pipeline configuration for %s/%d: %w", r.GetFullName(), input.GetNumber(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// create the objects from the pipeline in the database
	err = planBuild(database.FromContext(c), p, input, r)
	if err != nil {
		util.HandleError(c, http.StatusInternalServerError, err)

		return
	}

	// send API call to capture the created build
	input, _ = database.FromContext(c).GetBuild(input.GetNumber(), r)

	c.JSON(http.StatusCreated, input)

	// send API call to set the status on the commit
	err = source.FromContext(c).Status(u, input, r.GetOrg(), r.GetName())
	if err != nil {
		logrus.Errorf("unable to set commit status for build %s/%d: %v", r.GetFullName(), input.GetNumber(), err)
	}

	// publish the build to the queue
	go publishToQueue(
		queue.FromContext(c),
		p,
		input,
		r,
		u,
	)
}

// swagger:operation GET /api/v1/repos/{org}/{repo}/builds builds GetBuilds
//
// Create a build in the configured backend
//
// ---
// x-success_http_code: '200'
// produces:
// - application/json
// parameters:
// - in: path
//   name: repo
//   description: Name of the repo
//   required: true
//   type: string
// - in: path
//   name: org
//   description: Name of the org
//   required: true
//   type: string
// - in: header
//   name: Authorization
//   description: Vela bearer token
//   required: true
//   type: string
// responses:
//   '200':
//     description: Successfully retrieved the build
//     type: json
//     schema:
//       "$ref": "#/definitions/Build"
//   '400':
//     description: Unable to retrieve the list of builds
//     schema:
//       type: string
//   '500':
//     description: Unable to retrieve the list of builds
//     schema:
//       type: string

// GetBuilds represents the API handler to capture a
// list of builds for a repo from the configured backend.
func GetBuilds(c *gin.Context) {
	// variables that will hold the build list and total count
	var (
		b []*library.Build
		t int64
	)

	// capture middleware values
	r := repo.Retrieve(c)
	// capture the event type parameter
	event := c.Query("event")

	logrus.Infof("Reading builds for repo %s", r.GetFullName())

	// capture page query parameter if present
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		retErr := fmt.Errorf("unable to convert page query parameter for repo %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// capture per_page query parameter if present
	perPage, err := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	if err != nil {
		retErr := fmt.Errorf("unable to convert per_page query parameter for repo %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// ensure per_page isn't above or below allowed values
	perPage = util.MaxInt(1, util.MinInt(100, perPage))

	// send API call to capture the list of builds for the repo (and event type if passed in)
	if len(event) > 0 {
		b, t, err = database.FromContext(c).GetRepoBuildListByEvent(r, page, perPage, event)
	} else {
		b, t, err = database.FromContext(c).GetRepoBuildList(r, page, perPage) //[here] step 4
	}

	if err != nil {
		retErr := fmt.Errorf("unable to get builds for repo %s: %w", r.GetFullName(), err)

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

	c.JSON(http.StatusOK, b) //[here] this is where post is made.
}

func GetBuildsOrg(c *gin.Context) { //[here] Note: Function names are subject to change.
	// variables that will hold the build list and total count
	// var (
	// 	b []*library.Build
	// 	t int64
	// )

	// // capture middleware values
	// r := repo.Retrieve(c)
	// // capture the event type parameter
	// event := c.Query("event")

	// logrus.Infof("Reading builds for repo %s", r.GetFullName())

	// // capture page query parameter if present
	// page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	// if err != nil {
	// 	retErr := fmt.Errorf("unable to convert page query parameter for repo %s: %w", r.GetFullName(), err)

	// 	util.HandleError(c, http.StatusBadRequest, retErr)

	// 	return
	// }

	// // capture per_page query parameter if present
	// perPage, err := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	// if err != nil {
	// 	retErr := fmt.Errorf("unable to convert per_page query parameter for repo %s: %w", r.GetFullName(), err)

	// 	util.HandleError(c, http.StatusBadRequest, retErr)

	// 	return
	// }

	// // ensure per_page isn't above or below allowed values
	// perPage = util.MaxInt(1, util.MinInt(100, perPage))

	// // send API call to capture the list of builds for the repo (and event type if passed in)
	// if len(event) > 0 {
	// 	b, t, err = database.FromContext(c).GetRepoBuildListByEvent(r, page, perPage, event)
	// } else {
	// 	b, t, err = database.FromContext(c).GetOrgBuildList(r, page, perPage) //[here] step 4.5
	// }

	// if err != nil {
	// 	retErr := fmt.Errorf("unable to get builds for repo %s: %w", r.GetFullName(), err)

	// 	util.HandleError(c, http.StatusInternalServerError, retErr)

	// 	return
	// }

	// // create pagination object
	// pagination := Pagination{
	// 	Page:    page,
	// 	PerPage: perPage,
	// 	Total:   t,
	// }
	// // set pagination headers
	// pagination.SetHeaderLink(c)

	c.JSON(http.StatusOK, "You did it!") //[here] this is where post is made.
}

// swagger:operation POST /api/v1/repos/{org}/{repo}/builds/{build} builds GetBuild
//
// Get a build in the configured backend
//
// ---
// x-success_http_code: '200'
// produces:
// - application/json
// parameters:
// - in: path
//   name: repo
//   description: Name of the repo
//   required: true
//   type: string
// - in: path
//   name: org
//   description: Name of the org
//   required: true
//   type: string
// - in: path
//   name: build
//   description: Build number to restart
//   required: true
//   type: integer
// - in: header
//   name: Authorization
//   description: Vela bearer token
//   required: true
//   type: string
// responses:
//   '200':
//     description: Successfully restarted the build
//     type: json
//     schema:
//       "$ref": "#/definitions/Build"

// GetBuild represents the API handler to capture
// a build for a repo from the configured backend.
func GetBuild(c *gin.Context) {
	// capture middleware values
	r := repo.Retrieve(c)

	logrus.Infof("Reading build %s/%s", r.GetFullName(), c.Param("build"))

	// retrieve build from context
	b := build.Retrieve(c)

	c.JSON(http.StatusOK, b)
}

// swagger:operation POST /api/v1/repos/{org}/{repo}/builds/{build} builds RestartBuild
//
// Restart a build in the configured backend
//
// ---
// x-success_http_code: '201'
// produces:
// - application/json
// parameters:
// - in: path
//   name: repo
//   description: Name of the repo
//   required: true
//   type: string
// - in: path
//   name: org
//   description: Name of the org
//   required: true
//   type: string
// - in: path
//   name: build
//   description: Build number to restart
//   required: true
//   type: integer
// - in: header
//   name: Authorization
//   description: Vela bearer token
//   required: true
//   type: string
// responses:
//   '201':
//     description: Successfully restarted the build
//     type: json
//     schema:
//       "$ref": "#/definitions/Build"
//   '400':
//     description: Unable to restart the build
//     schema:
//       type: string
//   '404':
//     description: Unable to restart the build
//     schema:
//       type: string
//   '500':
//     description: Unable to restart the build
//     schema:
//       type: string

// RestartBuild represents the API handler to
// restart an existing build in the configured backend.
func RestartBuild(c *gin.Context) {
	// capture middleware values
	m := c.MustGet("metadata").(*types.Metadata)
	b := build.Retrieve(c)
	r := repo.Retrieve(c)

	logrus.Infof("Restarting build %s/%d", r.GetFullName(), b.GetNumber())

	// send API call to capture the repo owner
	u, err := database.FromContext(c).GetUser(r.GetUserID())
	if err != nil {
		retErr := fmt.Errorf("unable to get owner for %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// send API call to capture the last build for the repo
	lastBuild, err := database.FromContext(c).GetLastBuild(r)
	if err != nil {
		retErr := fmt.Errorf("unable to get last build for %s/%d: %w", r.GetFullName(), b.GetNumber(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// update fields in build object
	b.SetID(0)
	b.SetParent(b.GetNumber())
	b.SetNumber((lastBuild.GetNumber() + 1))
	b.SetCreated(time.Now().UTC().Unix())
	b.SetEnqueued(0)
	b.SetStarted(0)
	b.SetFinished(0)
	b.SetStatus(constants.StatusPending)
	b.SetHost("")
	b.SetRuntime("")
	b.SetDistribution("")

	// populate the build link if a web address is provided
	if len(m.Vela.WebAddress) > 0 {
		b.SetLink(
			fmt.Sprintf("%s/%s/%d", m.Vela.WebAddress, r.GetFullName(), b.GetNumber()),
		)
	}

	// variable to store changeset files
	var files []string
	// check if the build event is not pull_request
	if !strings.EqualFold(b.GetEvent(), constants.EventPull) {
		// send API call to capture list of files changed for the commit
		files, err = source.FromContext(c).Changeset(u, r, b.GetCommit())
		if err != nil {
			retErr := fmt.Errorf("unable to process webhook: failed to get changeset for %s: %w", r.GetFullName(), err)

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}
	}

	// files is empty if the build event is pull_request
	if len(files) == 0 {
		// parse out pull request number from base ref
		//
		// pattern: refs/pull/1/head
		var parts []string
		if strings.HasPrefix(b.GetRef(), "refs/pull/") {
			parts = strings.Split(b.GetRef(), "/")
		}

		// capture number by converting from string
		number, err := strconv.Atoi(parts[2])
		if err != nil {
			// capture number by scanning from string
			_, err := fmt.Sscanf(b.GetRef(), "%s/%s/%d/%s", nil, nil, &number, nil)
			if err != nil {
				retErr := fmt.Errorf("unable to process webhook: failed to get pull_request number for %s: %w", r.GetFullName(), err)

				util.HandleError(c, http.StatusInternalServerError, retErr)

				return
			}
		}

		// send API call to capture list of files changed for the pull request
		files, err = source.FromContext(c).ChangesetPR(u, r, number)
		if err != nil {
			retErr := fmt.Errorf("unable to process webhook: failed to get changeset for %s: %w", r.GetFullName(), err)

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}
	}

	// send API call to capture the pipeline configuration file
	config, err := source.FromContext(c).ConfigBackoff(u, r.GetOrg(), r.GetName(), b.GetCommit())
	if err != nil {
		retErr := fmt.Errorf("unable to get pipeline configuration for %s/%d: %w", r.GetFullName(), b.GetNumber(), err)

		util.HandleError(c, http.StatusNotFound, retErr)

		return
	}

	// parse and compile the pipeline configuration file
	p, err := compiler.FromContext(c).
		WithBuild(b).
		WithFiles(files).
		WithMetadata(m).
		WithRepo(r).
		WithUser(u).
		Compile(config)
	if err != nil {
		retErr := fmt.Errorf("unable to compile pipeline configuration for %s/%d: %w", r.GetFullName(), b.GetNumber(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// create the objects from the pipeline in the database
	err = planBuild(database.FromContext(c), p, b, r)
	if err != nil {
		util.HandleError(c, http.StatusInternalServerError, err)

		return
	}

	// send API call to capture the restarted build
	b, _ = database.FromContext(c).GetBuild(b.GetNumber(), r)

	c.JSON(http.StatusCreated, b)

	// send API call to set the status on the commit
	err = source.FromContext(c).Status(u, b, r.GetOrg(), r.GetName())
	if err != nil {
		logrus.Errorf("unable to set commit status for build %s/%d: %v", r.GetFullName(), b.GetNumber(), err)
	}

	// publish the build to the queue
	go publishToQueue(
		queue.FromContext(c),
		p,
		b,
		r,
		u,
	)
}

// swagger:operation PUT /api/v1/repos/{org}/{repo}/builds/{build} builds UpdateBuild
//
// Updates a build in the configured backend
//
// ---
// x-success_http_code: '200'
// produces:
// - application/json
// parameters:
// - in: body
//   name: body
//   description: Payload containing the build to update
//   required: true
//   schema:
//     "$ref": "#/definitions/Build"
// - in: path
//   name: repo
//   description: Name of the repo
//   required: true
//   type: string
// - in: path
//   name: org
//   description: Name of the org
//   required: true
//   type: string
// - in: path
//   name: build
//   description: Build number to restart
//   required: true
//   type: integer
// - in: header
//   name: Authorization
//   description: Vela bearer token
//   required: true
//   type: string
// responses:
//   '200':
//     description: Successfully restarted the build
//     type: json
//     schema:
//       "$ref": "#/definitions/Build"
//   '404':
//     description: Unable to restart the build
//     schema:
//       type: string
//   '500':
//     description: Unable to restart the build
//     schema:
//       type: string

// UpdateBuild represents the API handler to update
// a build for a repo in the configured backend.
func UpdateBuild(c *gin.Context) {
	// capture middleware values
	b := build.Retrieve(c)
	r := repo.Retrieve(c)

	logrus.Infof("Updating build %s/%d", r.GetFullName(), b.GetNumber())

	// capture body from API request
	input := new(library.Build)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for build %s/%d: %w", r.GetFullName(), b.GetNumber(), err)

		util.HandleError(c, http.StatusNotFound, retErr)

		return
	}

	// update build fields if provided
	if len(input.GetStatus()) > 0 {
		// update status if set
		b.SetStatus(input.GetStatus())
	}

	if len(input.GetError()) > 0 {
		// update error if set
		b.SetError(input.GetError())
	}

	if input.GetStarted() > 0 {
		// update started if set
		b.SetStarted(input.GetStarted())
	}

	if input.GetFinished() > 0 {
		// update finished if set
		b.SetFinished(input.GetFinished())
	}

	if len(input.GetTitle()) > 0 {
		// update title if set
		b.SetTitle(input.GetTitle())
	}

	if len(input.GetMessage()) > 0 {
		// update message if set
		b.SetMessage(input.GetMessage())
	}

	if len(input.GetHost()) > 0 {
		// update host if set
		b.SetHost(input.GetHost())
	}

	if len(input.GetRuntime()) > 0 {
		// update runtime if set
		b.SetRuntime(input.GetRuntime())
	}

	if len(input.GetDistribution()) > 0 {
		// update distribution if set
		b.SetDistribution(input.GetDistribution())
	}

	// send API call to update the build
	err = database.FromContext(c).UpdateBuild(b)
	if err != nil {
		retErr := fmt.Errorf("unable to update build %s/%d: %w", r.GetFullName(), b.GetNumber(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// send API call to capture the updated build
	b, _ = database.FromContext(c).GetBuild(b.GetNumber(), r)

	c.JSON(http.StatusOK, b)

	// check if the build is in a "final" state
	if b.GetStatus() == constants.StatusSuccess ||
		b.GetStatus() == constants.StatusFailure ||
		b.GetStatus() == constants.StatusKilled ||
		b.GetStatus() == constants.StatusError {
		// send API call to capture the repo owner
		u, err := database.FromContext(c).GetUser(r.GetUserID())
		if err != nil {
			logrus.Errorf("unable to get owner for %s/%d: %v", r.GetFullName(), b.GetNumber(), err)
		}

		// send API call to set the status on the commit
		err = source.FromContext(c).Status(u, b, r.GetOrg(), r.GetName())
		if err != nil {
			logrus.Errorf("unable to set commit status for build %s/%d: %v", r.GetFullName(), b.GetNumber(), err)
		}
	}
}

// swagger:operation DELETE /api/v1/repos/{org}/{repo}/builds/{build} builds DeleteBuild
//
// Delete a build in the configured backend
//
// ---
// x-success_http_code: '200'
// produces:
// - application/json
// parameters:
// - in: path
//   name: repo
//   description: Name of the repo
//   required: true
//   type: string
// - in: path
//   name: org
//   description: Name of the org
//   required: true
//   type: string
// - in: path
//   name: build
//   description: Build number to restart
//   required: true
//   type: integer
// - in: header
//   name: Authorization
//   description: Vela bearer token
//   required: true
//   type: string
// responses:
//   '200':
//     description: Successfully restarted the build
//     schema:
//       type: string
//   '400':
//     description: Unable to restart the build
//     schema:
//       type: string

// DeleteBuild represents the API handler to remove
// a build for a repo from the configured backend.
func DeleteBuild(c *gin.Context) {
	// capture middleware values
	b := build.Retrieve(c)
	r := repo.Retrieve(c)

	logrus.Infof("Deleting build %s/%d", r.GetFullName(), b.GetNumber())

	// send API call to remove the build
	err := database.FromContext(c).DeleteBuild(b.GetID())
	if err != nil {
		retErr := fmt.Errorf("unable to delete build %s/%d: %v", r.GetFullName(), b.GetNumber(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("Build %s/%d deleted", r.GetFullName(), b.GetNumber()))
}

// planBuild is a helper function to plan the build for
// execution. This creates all resources, like steps
// and services, for the build in the configured backend.
func planBuild(database database.Service, p *pipeline.Build, b *library.Build, r *library.Repo) error {
	// update fields in build object
	b.SetEnqueued(time.Now().UTC().Unix())

	// send API call to create the build
	err := database.CreateBuild(b)
	if err != nil {
		cleanBuild(database, b, nil, nil)

		return fmt.Errorf("unable to create new build for %s: %v", r.GetFullName(), err)
	}

	// send API call to capture the created build
	b, _ = database.GetBuild(int(b.GetNumber()), r)

	// plan all services for the build
	services, err := planServices(database, p, b)
	if err != nil {
		cleanBuild(database, b, services, nil)

		return err
	}

	// plan all steps for the build
	steps, err := planSteps(database, p, b)
	if err != nil {
		cleanBuild(database, b, services, steps)

		return err
	}

	return nil
}

// cleanBuild is a helper function to kill the build
// without execution. This will kill all resources,
// like steps and services, for the build in the
// configured backend.
func cleanBuild(database database.Service, b *library.Build, services []*library.Service, steps []*library.Step) {
	// update fields in build object
	b.SetError("unable to publish build to queue")
	b.SetStatus(constants.StatusError)
	b.SetFinished(time.Now().UTC().Unix())

	// send API call to update the build
	err := database.UpdateBuild(b)
	if err != nil {
		logrus.Errorf("unable to kill build %d: %v", b.GetNumber(), err)
	}

	for _, s := range services {
		// update fields in service object
		s.SetStatus(constants.StatusKilled)
		s.SetFinished(time.Now().UTC().Unix())

		// send API call to update the service
		err := database.UpdateService(s)
		if err != nil {
			logrus.Errorf("unable to kill service %s for build %d: %v", s.GetName(), b.GetNumber(), err)
		}
	}

	for _, s := range steps {
		// update fields in step object
		s.SetStatus(constants.StatusKilled)
		s.SetFinished(time.Now().UTC().Unix())

		// send API call to update the step
		err := database.UpdateStep(s)
		if err != nil {
			logrus.Errorf("unable to kill step %s for build %d: %v", s.GetName(), b.GetNumber(), err)
		}
	}
}
