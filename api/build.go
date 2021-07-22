// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-vela/compiler/compiler"

	"github.com/go-vela/pkg-queue/queue"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/executors"
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
// produces:
// - application/json
// parameters:
// - in: path
//   name: org
//   description: Name of the org
//   required: true
//   type: string
// - in: path
//   name: repo
//   description: Name of the repo
//   required: true
//   type: string
// - in: body
//   name: body
//   description: Payload containing the build to update
//   required: true
//   schema:
//     "$ref": "#/definitions/Build"
// security:
//   - ApiKeyAuth: []
// responses:
//   '201':
//     description: Successfully created the build
//     type: json
//     schema:
//       "$ref": "#/definitions/Build"
//   '400':
//     description: Unable to create the build
//     schema:
//       "$ref": "#/definitions/Error"
//   '404':
//     description: Unable to create the build
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to create the build
//     schema:
//       "$ref": "#/definitions/Error"

// CreateBuild represents the API handler to
// create a build in the configured backend.
//
// nolint: funlen // ignore function length due to comments
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
		// nolint: lll // ignore long line length due to error message
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
			// nolint: lll // ignore long line length due to error message
			retErr := fmt.Errorf("unable to process webhook: failed to get changeset for %s: %w", r.GetFullName(), err)

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}
	}

	// handle getting changeset from a pull_request
	if strings.EqualFold(input.GetEvent(), constants.EventPull) {
		// capture number from build
		number, err := getPRNumberFromBuild(input)
		if err != nil {
			// nolint: lll // ignore long line length due to error message
			retErr := fmt.Errorf("unable to create build: failed to get pull_request number for %s: %w", r.GetFullName(), err)

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}

		// send API call to capture list of files changed for the pull request
		files, err = source.FromContext(c).ChangesetPR(u, r, number)
		if err != nil {
			// nolint: lll // ignore long line length due to error message
			retErr := fmt.Errorf("unable to process webhook: failed to get changeset for %s: %w", r.GetFullName(), err)

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}
	}

	// send API call to capture the pipeline configuration file
	config, err := source.FromContext(c).ConfigBackoff(u, r.GetOrg(), r.GetName(), input.GetCommit())
	if err != nil {
		// nolint: lll // ignore long line length due to error message
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
		// nolint: lll // ignore long line length due to error message
		retErr := fmt.Errorf("unable to compile pipeline configuration for %s/%d: %w", r.GetFullName(), input.GetNumber(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// skip the build if only the init or clone steps are found
	skip := skipEmptyBuild(p)
	if skip != "" {
		c.JSON(http.StatusOK, skip)
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
		// nolint: lll // ignore long line length due to error message
		logrus.Errorf("unable to set commit status for build %s/%d: %v", r.GetFullName(), input.GetNumber(), err)
	}

	// publish the build to the queue
	go publishToQueue(
		queue.FromGinContext(c),
		p,
		input,
		r,
		u,
	)
}

// skipEmptyBuild checks if the build should be skipped due to it
// not containing any steps besides init or clone
func skipEmptyBuild(p *pipeline.Build) string {
	if len(p.Stages) == 1 {
		if p.Stages[0].Name == "init" {
			return "skipping build since only init stage found"
		}
	}

	if len(p.Stages) == 2 {
		if p.Stages[0].Name == "init" && p.Stages[1].Name == "clone" {
			return "skipping build since only init and clone stages found"
		}
	}

	if len(p.Steps) == 1 {
		if p.Steps[0].Name == "init" {
			return "skipping build since only init step found"
		}
	}

	if len(p.Steps) == 2 {
		if p.Steps[0].Name == "init" && p.Steps[1].Name == "clone" {
			return "skipping build since only init and clone steps found"
		}
	}

	return ""
}

// swagger:operation GET /api/v1/repos/{org}/{repo}/builds builds GetBuilds
//
// Get builds from the configured backend
//
// ---
// produces:
// - application/json
// parameters:
// - in: path
//   name: org
//   description: Name of the org
//   required: true
//   type: string
// - in: path
//   name: repo
//   description: Name of the repo
//   required: true
//   type: string
// - in: query
//   name: event
//   description: Filter by build event
//   type: string
//   enum:
//   - push
//   - pull_request
//   - tag
//   - deployment
//   - comment
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
//     description: Successfully retrieved the build
//     schema:
//       type: array
//       items:
//         "$ref": "#/definitions/Build"
//     headers:
//       X-Total-Count:
//         description: Total number of results
//         type: integer
//       Link:
//         description: see https://tools.ietf.org/html/rfc5988
//         type: string
//   '400':
//     description: Unable to retrieve the list of builds
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to retrieve the list of builds
//     schema:
//       "$ref": "#/definitions/Error"

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
		// nolint: lll // ignore long line length due to error message
		retErr := fmt.Errorf("unable to convert page query parameter for repo %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// capture per_page query parameter if present
	perPage, err := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	if err != nil {
		// nolint: lll // ignore long line length due to error message
		retErr := fmt.Errorf("unable to convert per_page query parameter for repo %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// ensure per_page isn't above or below allowed values
	//
	// nolint: gomnd // ignore magic number
	perPage = util.MaxInt(1, util.MinInt(100, perPage))

	// send API call to capture the list of builds for the repo (and event type if passed in)
	if len(event) > 0 {
		b, t, err = database.FromContext(c).GetRepoBuildListByEvent(r, event, page, perPage)
	} else {
		b, t, err = database.FromContext(c).GetRepoBuildList(r, page, perPage)
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

	c.JSON(http.StatusOK, b)
}

// swagger:operation GET /api/v1/repos/{org} builds GetOrgBuilds
//
// Get a list of builds by org in the configured backend
//
// ---
// produces:
// - application/json
// parameters:
// - in: path
//   name: org
//   description: Name of the org
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
//     description: Successfully retrieved build list
//     schema:
//       type: array
//       items:
//         "$ref": "#/definitions/Build"
//     headers:
//       X-Total-Count:
//         description: Total number of results
//         type: integer
//       Link:
//         description: see https://tools.ietf.org/html/rfc5988
//         type: string
//   '400':
//     description: Unable to retrieve the list of builds
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to retrieve the list of builds
//     schema:
//       "$ref": "#/definitions/Error"

// GetOrgBuilds represents the API handler to capture a
// list of builds associated with an org from the configured backend.
func GetOrgBuilds(c *gin.Context) {
	// variables that will hold the build list and total count
	var (
		b []*library.Build
		t int64
	)

	// capture middleware values
	o := c.Param("org")
	// capture the event type parameter
	event := c.Query("event")

	logrus.Infof("Reading builds for org %s", o)

	// capture page query parameter if present
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		retErr := fmt.Errorf("unable to convert page query parameter for org %s: %w", o, err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// capture per_page query parameter if present
	perPage, err := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	if err != nil {
		retErr := fmt.Errorf("unable to convert per_page query parameter for Org %s: %w", o, err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// ensure per_page isn't above or below allowed values
	//
	// nolint: gomnd // ignore magic number
	perPage = util.MaxInt(1, util.MinInt(100, perPage))

	// send API call to capture the list of builds for the org (and event type if passed in)
	if len(event) > 0 {
		b, t, err = database.FromContext(c).GetOrgBuildListByEvent(o, event, page, perPage)
	} else {
		b, t, err = database.FromContext(c).GetOrgBuildList(o, page, perPage)
	}

	if err != nil {
		retErr := fmt.Errorf("unable to get builds for org %s: %w", o, err)

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

	c.JSON(http.StatusOK, b)
}

// swagger:operation GET /api/v1/repos/{org}/{repo}/builds/{build} builds GetBuild
//
// Get a build in the configured backend
//
// ---
// produces:
// - application/json
// parameters:
// - in: path
//   name: org
//   description: Name of the org
//   required: true
//   type: string
// - in: path
//   name: repo
//   description: Name of the repo
//   required: true
//   type: string
// - in: path
//   name: build
//   description: Build number to restart
//   required: true
//   type: integer
// security:
//   - ApiKeyAuth: []
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
// produces:
// - application/json
// parameters:
// - in: path
//   name: org
//   description: Name of the org
//   required: true
//   type: string
// - in: path
//   name: repo
//   description: Name of the repo
//   required: true
//   type: string
// - in: path
//   name: build
//   description: Build number to restart
//   required: true
//   type: integer
// security:
//   - ApiKeyAuth: []
// responses:
//   '201':
//     description: Successfully restarted the build
//     schema:
//       "$ref": "#/definitions/Build"
//   '400':
//     description: Unable to restart the build
//     schema:
//       "$ref": "#/definitions/Error"
//   '404':
//     description: Unable to restart the build
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to restart the build
//     schema:
//       "$ref": "#/definitions/Error"

// RestartBuild represents the API handler to
// restart an existing build in the configured backend.
//
// nolint: funlen // ignore function length due to comments
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
		// nolint: lll // ignore long line length due to error message
		retErr := fmt.Errorf("unable to get last build for %s/%d: %w", r.GetFullName(), b.GetNumber(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// update fields in build object
	b.SetID(0)
	b.SetNumber(r.GetCounter())
	b.SetParent(lastBuild.GetNumber())
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
			// nolint: lll // ignore long line length due to error message
			retErr := fmt.Errorf("unable to process webhook: failed to get changeset for %s: %w", r.GetFullName(), err)

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}
	}

	// handle getting changeset from a pull_request
	if strings.EqualFold(b.GetEvent(), constants.EventPull) {
		// capture number from build
		number, err := getPRNumberFromBuild(b)
		if err != nil {
			// nolint: lll // ignore long line length due to error message
			retErr := fmt.Errorf("unable to restart build: failed to get pull_request number for %s: %w", r.GetFullName(), err)

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}

		// send API call to capture list of files changed for the pull request
		files, err = source.FromContext(c).ChangesetPR(u, r, number)
		if err != nil {
			// nolint: lll // ignore long line length due to error message
			retErr := fmt.Errorf("unable to process webhook: failed to get changeset for %s: %w", r.GetFullName(), err)

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}
	}

	// send API call to capture the pipeline configuration file
	config, err := source.FromContext(c).ConfigBackoff(u, r.GetOrg(), r.GetName(), b.GetCommit())
	if err != nil {
		// nolint: lll // ignore long line length due to error message
		retErr := fmt.Errorf("unable to get pipeline configuration for %s/%d: %w", r.GetFullName(), b.GetNumber(), err)

		util.HandleError(c, http.StatusNotFound, retErr)

		return
	}

	// update the build numbers based off repo counter
	inc := r.GetCounter() + 1

	r.SetCounter(inc)
	b.SetNumber(inc)

	// parse and compile the pipeline configuration file
	p, err := compiler.FromContext(c).
		WithBuild(b).
		WithFiles(files).
		WithMetadata(m).
		WithRepo(r).
		WithUser(u).
		Compile(config)
	if err != nil {
		// nolint: lll // ignore long line length due to error message
		retErr := fmt.Errorf("unable to compile pipeline configuration for %s/%d: %w", r.GetFullName(), b.GetNumber(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// skip the build if only the init or clone steps are found
	skip := skipEmptyBuild(p)
	if skip != "" {
		c.JSON(http.StatusOK, skip)
		return
	}

	// create the objects from the pipeline in the database
	err = planBuild(database.FromContext(c), p, b, r)
	if err != nil {
		util.HandleError(c, http.StatusInternalServerError, err)

		return
	}

	// send API call to update repo for ensuring counter is incremented
	err = database.FromContext(c).UpdateRepo(r)
	if err != nil {
		retErr := fmt.Errorf("%s: failed to update repo %s: %v", baseErr, r.GetFullName(), err)
		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// send API call to capture the restarted build
	b, _ = database.FromContext(c).GetBuild(b.GetNumber(), r)

	c.JSON(http.StatusCreated, b)

	// send API call to set the status on the commit
	err = source.FromContext(c).Status(u, b, r.GetOrg(), r.GetName())
	if err != nil {
		// nolint: lll // ignore long line length due to error message
		logrus.Errorf("unable to set commit status for build %s/%d: %v", r.GetFullName(), b.GetNumber(), err)
	}

	// publish the build to the queue
	go publishToQueue(
		queue.FromGinContext(c),
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
// produces:
// - application/json
// parameters:
// - in: path
//   name: org
//   description: Name of the org
//   required: true
//   type: string
// - in: path
//   name: repo
//   description: Name of the repo
//   required: true
//   type: string
// - in: path
//   name: build
//   description: Build number to update
//   required: true
//   type: integer
// - in: body
//   name: body
//   description: Payload containing the build to update
//   required: true
//   schema:
//     "$ref": "#/definitions/Build"
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully updated the build
//     schema:
//       "$ref": "#/definitions/Build"
//   '404':
//     description: Unable to update the build
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to update the build
//     schema:
//       "$ref": "#/definitions/Error"

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
		// nolint: lll // ignore long line length due to error message
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

	if input.GetEnqueued() > 0 {
		// update enqueued if set
		b.SetEnqueued(input.GetEnqueued())
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
		b.GetStatus() == constants.StatusCanceled ||
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
			// nolint: lll // ignore long line length due to error message
			logrus.Errorf("unable to set commit status for build %s/%d: %v", r.GetFullName(), b.GetNumber(), err)
		}
	}
}

// swagger:operation DELETE /api/v1/repos/{org}/{repo}/builds/{build} builds DeleteBuild
//
// Delete a build in the configured backend
//
// ---
// produces:
// - application/json
// parameters:
// - in: path
//   name: org
//   description: Name of the org
//   required: true
//   type: string
// - in: path
//   name: repo
//   description: Name of the repo
//   required: true
//   type: string
// - in: path
//   name: build
//   description: Build number to delete
//   required: true
//   type: integer
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully deleted the build
//     schema:
//       type: string
//   '400':
//     description: Unable to delete the build
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to delete the build
//     schema:
//       "$ref": "#/definitions/Error"

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

// getPRNumberFromBuild is a helper function to
// extract the pull request number from a Build.
func getPRNumberFromBuild(b *library.Build) (int, error) {
	// parse out pull request number from base ref
	//
	// pattern: refs/pull/1/head
	var parts []string
	if strings.HasPrefix(b.GetRef(), "refs/pull/") {
		parts = strings.Split(b.GetRef(), "/")
	}

	// just being safe to avoid out of range index errors
	//
	// nolint:gomnd // magic number of 3 used once
	if len(parts) < 3 {
		return 0, fmt.Errorf("invalid ref: %s", b.GetRef())
	}

	// return the results of converting number to string
	return strconv.Atoi(parts[2])
}

// planBuild is a helper function to plan the build for
// execution. This creates all resources, like steps
// and services, for the build in the configured backend.
//
// nolint: lll // ignore long line length due to variable names
func planBuild(database database.Service, p *pipeline.Build, b *library.Build, r *library.Repo) error {
	// update fields in build object
	b.SetCreated(time.Now().UTC().Unix())

	// send API call to create the build
	err := database.CreateBuild(b)
	if err != nil {
		// clean up the objects from the pipeline in the database
		cleanBuild(database, b, nil, nil)

		return fmt.Errorf("unable to create new build for %s: %v", r.GetFullName(), err)
	}

	// send API call to capture the created build
	b, _ = database.GetBuild(b.GetNumber(), r)

	// plan all services for the build
	services, err := planServices(database, p, b)
	if err != nil {
		// clean up the objects from the pipeline in the database
		cleanBuild(database, b, services, nil)

		return err
	}

	// plan all steps for the build
	steps, err := planSteps(database, p, b)
	if err != nil {
		// clean up the objects from the pipeline in the database
		cleanBuild(database, b, services, steps)

		return err
	}

	return nil
}

// cleanBuild is a helper function to kill the build
// without execution. This will kill all resources,
// like steps and services, for the build in the
// configured backend.
//
// nolint: lll // ignore long line length due to variable names
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

// swagger:operation DELETE /api/v1/repos/{org}/{repo}/builds/{build}/cancel builds CancelBuild
//
// Cancel a running build
//
// ---
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
//   description: Build number to cancel
//   required: true
//   type: integer
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully canceled the build
//     schema:
//       type: string
//   '400':
//     description: Unable to cancel build
//     schema:
//       "$ref": "#/definitions/Error"
//   '404':
//     description: Unable to cancel build
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to cancel build
//     schema:
//       "$ref": "#/definitions/Error"

// CancelBuild represents the API handler to
// cancel a running build.
func CancelBuild(c *gin.Context) {
	r := repo.Retrieve(c)
	b := build.Retrieve(c)
	e := executors.Retrieve(c)

	// check to see if build is pending
	// todo: remove builds from the queue
	if strings.EqualFold(b.GetStatus(), constants.StatusPending) {
		retErr := fmt.Errorf("found build %s/%d but its status was %s",
			r.GetFullName(),
			b.GetNumber(),
			b.GetStatus(),
		)
		util.HandleError(c, http.StatusBadRequest, retErr)
		return
	}

	// check to see if build is not running
	// https://github.com/go-vela/types/blob/master/constants/status.go
	if !strings.EqualFold(b.GetStatus(), constants.StatusRunning) {
		retErr := fmt.Errorf("found build %s/%d but its status was %s",
			r.GetFullName(),
			b.GetNumber(),
			b.GetStatus(),
		)
		util.HandleError(c, http.StatusBadRequest, retErr)
		return
	}

	// retrieve the worker info
	w, err := database.FromContext(c).GetWorker(b.GetHost())
	if err != nil {
		retErr := fmt.Errorf("unable to get worker: %w", err)
		util.HandleError(c, http.StatusNotFound, retErr)
		return
	}

	for _, executor := range e {
		// check each executor on the worker running the build
		// to see if it's running the build we want to cancel
		//
		// nolint:whitespace // ignore leading newline to improve readability
		if strings.EqualFold(executor.Repo.GetFullName(), r.GetFullName()) &&
			*executor.GetBuild().Number == b.GetNumber() {

			// prepare the request to the worker
			client := http.DefaultClient
			client.Timeout = 30 * time.Second

			// set the API endpoint path we send the request to
			u := fmt.Sprintf("%s/api/v1/executors/%d/build/cancel", w.GetAddress(), executor.GetID())
			req, err := http.NewRequest("DELETE", u, nil)
			if err != nil {
				retErr := fmt.Errorf("unable to form a request to %s: %w", u, err)
				util.HandleError(c, http.StatusBadRequest, retErr)
				return
			}

			// add the token to authenticate to the worker
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.MustGet("secret").(string)))

			// perform the request to the worker
			resp, err := client.Do(req)
			if err != nil {
				retErr := fmt.Errorf("unable to connect to %s: %w", u, err)
				util.HandleError(c, http.StatusBadRequest, retErr)
				return
			}
			defer resp.Body.Close()

			// Read Response Body
			respBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				retErr := fmt.Errorf("unable to read response from %s: %w", u, err)
				util.HandleError(c, http.StatusBadRequest, retErr)
				return
			}

			err = json.Unmarshal(respBody, b)
			if err != nil {
				retErr := fmt.Errorf("unable to parse response from %s: %w", u, err)
				util.HandleError(c, http.StatusBadRequest, retErr)
				return
			}

			c.JSON(resp.StatusCode, b)
			return
		}
	}

	// build has been abandoned
	retErr := fmt.Errorf("unable to find a running build for %s/%d", r.GetFullName(), b.GetNumber())
	util.HandleError(c, http.StatusInternalServerError, retErr)
}
