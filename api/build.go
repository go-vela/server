// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-vela/server/internal/token"
	"github.com/go-vela/server/router/middleware/claims"
	"github.com/go-vela/server/router/middleware/org"

	"github.com/go-vela/server/compiler"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/queue"
	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/executors"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/scm"
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
//   '200':
//     description: Request processed but build was skipped
//     schema:
//       type: string
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

// CreateBuild represents the API handler to create a build in the configured backend.
//
//nolint:funlen,gocyclo // ignore function length and cyclomatic complexity
func CreateBuild(c *gin.Context) {
	// capture middleware values
	m := c.MustGet("metadata").(*types.Metadata)
	o := org.Retrieve(c)
	r := repo.Retrieve(c)
	u := user.Retrieve(c)

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logger := logrus.WithFields(logrus.Fields{
		"org":  o,
		"repo": r.GetName(),
		"user": u.GetName(),
	})

	logger.Infof("creating new build for repo %s", r.GetFullName())

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
	u, err = database.FromContext(c).GetUser(r.GetUserID())
	if err != nil {
		retErr := fmt.Errorf("unable to get owner for %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// create SQL filters for querying pending and running builds for repo
	filters := map[string]interface{}{
		"status": []string{constants.StatusPending, constants.StatusRunning},
	}

	// send API call to capture the number of pending or running builds for the repo
	builds, err := database.FromContext(c).GetRepoBuildCount(r, filters)
	if err != nil {
		retErr := fmt.Errorf("unable to create new build: unable to get count of builds for repo %s", r.GetFullName())

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// check if the number of pending and running builds exceeds the limit for the repo
	if builds >= r.GetBuildLimit() {
		retErr := fmt.Errorf("unable to create new build: repo %s has exceeded the concurrent build limit of %d", r.GetFullName(), r.GetBuildLimit())

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// update fields in build object
	input.SetRepoID(r.GetID())
	input.SetStatus(constants.StatusPending)
	input.SetCreated(time.Now().UTC().Unix())

	// set the parent equal to the current repo counter
	input.SetParent(r.GetCounter())
	// check if the parent is set to 0
	if input.GetParent() == 0 {
		// parent should be "1" if it's the first build ran
		input.SetParent(1)
	}

	// update the build numbers based off repo counter
	inc := r.GetCounter() + 1
	r.SetCounter(inc)
	input.SetNumber(inc)

	// populate the build link if a web address is provided
	if len(m.Vela.WebAddress) > 0 {
		input.SetLink(
			fmt.Sprintf("%s/%s/%d", m.Vela.WebAddress, r.GetFullName(), input.GetNumber()),
		)
	}

	// variable to store changeset files
	var files []string
	// check if the build event is not issue_comment or pull_request
	if !strings.EqualFold(input.GetEvent(), constants.EventComment) &&
		!strings.EqualFold(input.GetEvent(), constants.EventPull) {
		// send API call to capture list of files changed for the commit
		files, err = scm.FromContext(c).Changeset(u, r, input.GetCommit())
		if err != nil {
			retErr := fmt.Errorf("unable to create new build: failed to get changeset for %s: %w", r.GetFullName(), err)

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}
	}

	// check if the build event is a pull_request
	if strings.EqualFold(input.GetEvent(), constants.EventPull) {
		// capture number from build
		number, err := getPRNumberFromBuild(input)
		if err != nil {
			retErr := fmt.Errorf("unable to create new build: failed to get pull_request number for %s: %w", r.GetFullName(), err)

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}

		// send API call to capture list of files changed for the pull request
		files, err = scm.FromContext(c).ChangesetPR(u, r, number)
		if err != nil {
			retErr := fmt.Errorf("unable to create new build: failed to get changeset for %s: %w", r.GetFullName(), err)

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}
	}

	var (
		// variable to store the raw pipeline configuration
		config []byte
		// variable to store executable pipeline
		p *pipeline.Build
		// variable to store pipeline configuration
		pipeline *library.Pipeline
		// variable to store the pipeline type for the repository
		pipelineType = r.GetPipelineType()
	)

	// send API call to attempt to capture the pipeline
	pipeline, err = database.FromContext(c).GetPipelineForRepo(input.GetCommit(), r)
	if err != nil { // assume the pipeline doesn't exist in the database yet
		// send API call to capture the pipeline configuration file
		config, err = scm.FromContext(c).ConfigBackoff(u, r, input.GetCommit())
		if err != nil {
			retErr := fmt.Errorf("unable to create new build: failed to get pipeline configuration for %s: %w", r.GetFullName(), err)

			util.HandleError(c, http.StatusNotFound, retErr)

			return
		}
	} else {
		config = pipeline.GetData()
	}

	// ensure we use the expected pipeline type when compiling
	//
	// The pipeline type for a repo can change at any time which can break compiling
	// existing pipelines in the system for that repo. To account for this, we update
	// the repo pipeline type to match what was defined for the existing pipeline
	// before compiling. After we're done compiling, we reset the pipeline type.
	if len(pipeline.GetType()) > 0 {
		r.SetPipelineType(pipeline.GetType())
	}

	var compiled *library.Pipeline
	// parse and compile the pipeline configuration file
	p, compiled, err = compiler.FromContext(c).
		Duplicate().
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
	// reset the pipeline type for the repo
	//
	// The pipeline type for a repo can change at any time which can break compiling
	// existing pipelines in the system for that repo. To account for this, we update
	// the repo pipeline type to match what was defined for the existing pipeline
	// before compiling. After we're done compiling, we reset the pipeline type.
	r.SetPipelineType(pipelineType)

	// skip the build if only the init or clone steps are found
	skip := skipEmptyBuild(p)
	if skip != "" {
		// set build to successful status
		input.SetStatus(constants.StatusSuccess)

		// send API call to set the status on the commit
		err = scm.FromContext(c).Status(u, input, r.GetOrg(), r.GetName())
		if err != nil {
			logger.Errorf("unable to set commit status for %s/%d: %v", r.GetFullName(), input.GetNumber(), err)
		}

		c.JSON(http.StatusOK, skip)

		return
	}

	// check if the pipeline did not already exist in the database
	//
	//nolint:dupl // ignore duplicate code
	if pipeline == nil {
		pipeline = compiled
		pipeline.SetRepoID(r.GetID())
		pipeline.SetCommit(input.GetCommit())
		pipeline.SetRef(input.GetRef())

		// send API call to create the pipeline
		err = database.FromContext(c).CreatePipeline(pipeline)
		if err != nil {
			retErr := fmt.Errorf("unable to create new build: failed to create pipeline for %s: %w", r.GetFullName(), err)

			util.HandleError(c, http.StatusBadRequest, retErr)

			return
		}

		// send API call to capture the created pipeline
		pipeline, err = database.FromContext(c).GetPipelineForRepo(pipeline.GetCommit(), r)
		if err != nil {
			//nolint:lll // ignore long line length due to error message
			retErr := fmt.Errorf("unable to create new build: failed to get new pipeline %s/%s: %w", r.GetFullName(), pipeline.GetCommit(), err)

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}
	}

	input.SetPipelineID(pipeline.GetID())

	// create the objects from the pipeline in the database
	err = planBuild(database.FromContext(c), p, input, r)
	if err != nil {
		util.HandleError(c, http.StatusInternalServerError, err)

		return
	}

	// send API call to update repo for ensuring counter is incremented
	err = database.FromContext(c).UpdateRepo(r)
	if err != nil {
		retErr := fmt.Errorf("unable to create new build: failed to update repo %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// send API call to capture the created build
	input, _ = database.FromContext(c).GetBuild(input.GetNumber(), r)

	c.JSON(http.StatusCreated, input)

	// send API call to set the status on the commit
	err = scm.FromContext(c).Status(u, input, r.GetOrg(), r.GetName())
	if err != nil {
		logger.Errorf("unable to set commit status for build %s/%d: %v", r.GetFullName(), input.GetNumber(), err)
	}

	// publish the build to the queue
	go publishToQueue(
		queue.FromGinContext(c),
		database.FromContext(c),
		p,
		input,
		r,
		u,
	)
}

// skipEmptyBuild checks if the build should be skipped due to it
// not containing any steps besides init or clone.
//
//nolint:goconst // ignore init and clone constants
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

// swagger:operation GET /api/v1/search/builds/{id} builds GetBuildByID
//
// Get a single build by its id in the configured backend
//
// ---
// produces:
// - application/json
// parameters:
// - in: path
//   name: id
//   description: build id
//   required: true
//   type: number
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved build
//     schema:
//       "$ref": "#/definitions/Build"
//   '400':
//     description: Unable to retrieve the build
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to retrieve the build
//     schema:
//       "$ref": "#/definitions/Error"

// GetBuildByID represents the API handler to capture a
// build by its id from the configured backend.
func GetBuildByID(c *gin.Context) {
	// Variables that will hold the library types of the build and repo
	var (
		b *library.Build
		r *library.Repo
	)

	// Capture user from middleware
	u := user.Retrieve(c)

	// Parse build ID from path
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		retErr := fmt.Errorf("unable to parse build id: %w", err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"build": id,
		"user":  u.GetName(),
	}).Infof("reading build %d", id)

	// Get build from database
	b, err = database.FromContext(c).GetBuildByID(id)
	if err != nil {
		retErr := fmt.Errorf("unable to get build: %w", err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// Get repo from database using repo ID field from build
	r, err = database.FromContext(c).GetRepo(b.GetRepoID())
	if err != nil {
		retErr := fmt.Errorf("unable to get repo: %w", err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// Capture user access from SCM. We do this in order to ensure user has access and is not
	// just retrieving any build using a random id number.
	perm, err := scm.FromContext(c).RepoAccess(u, u.GetToken(), r.GetOrg(), r.GetName())
	if err != nil {
		logrus.Errorf("unable to get user %s access level for repo %s", u.GetName(), r.GetFullName())
	}

	// Ensure that user has at least read access to repo to return the build
	if perm == "none" && !u.GetAdmin() {
		retErr := fmt.Errorf("unable to retrieve build %d: user does not have read access to repo %s", id, r.GetFullName())

		util.HandleError(c, http.StatusUnauthorized, retErr)

		return
	}

	c.JSON(http.StatusOK, b)
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
//   name: commit
//   description: Filter builds based on the commit hash
//   type: string
// - in: query
//   name: branch
//   description: Filter builds by branch
//   type: string
// - in: query
//   name: status
//   description: Filter by build status
//   type: string
//   enum:
//   - canceled
//   - error
//   - failure
//   - killed
//   - pending
//   - running
//   - success
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
// - in: query
//   name: before
//   description: filter builds created before a certain time
//   type: integer
//   default: 1
// - in: query
//   name: after
//   description: filter builds created after a certain time
//   type: integer
//   default: 0
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved the builds
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
	// variables that will hold the build list, build list filters and total count
	var (
		filters = map[string]interface{}{}
		b       []*library.Build
		t       int64
	)

	// capture middleware values
	o := org.Retrieve(c)
	r := repo.Retrieve(c)
	u := user.Retrieve(c)

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"org":  o,
		"repo": r.GetName(),
		"user": u.GetName(),
	}).Infof("reading builds for repo %s", r.GetFullName())

	// capture the branch name parameter
	branch := c.Query("branch")
	// capture the event type parameter
	event := c.Query("event")
	// capture the status type parameter
	status := c.Query("status")
	// capture the commit hash parameter
	commit := c.Query("commit")

	// check if branch filter was provided
	if len(branch) > 0 {
		// add branch to filters map
		filters["branch"] = branch
	}
	// check if event filter was provided
	if len(event) > 0 {
		// verify the event provided is a valid event type
		if event != constants.EventComment && event != constants.EventDeploy &&
			event != constants.EventPush && event != constants.EventPull &&
			event != constants.EventTag {
			retErr := fmt.Errorf("unable to process event %s: invalid event type provided", event)

			util.HandleError(c, http.StatusBadRequest, retErr)

			return
		}

		// add event to filters map
		filters["event"] = event
	}
	// check if status filter was provided
	if len(status) > 0 {
		// verify the status provided is a valid status type
		if status != constants.StatusCanceled && status != constants.StatusError &&
			status != constants.StatusFailure && status != constants.StatusKilled &&
			status != constants.StatusPending && status != constants.StatusRunning &&
			status != constants.StatusSuccess {
			retErr := fmt.Errorf("unable to process status %s: invalid status type provided", status)

			util.HandleError(c, http.StatusBadRequest, retErr)

			return
		}

		// add status to filters map
		filters["status"] = status
	}

	// check if commit hash filter was provided
	if len(commit) > 0 {
		// add commit to filters map
		filters["commit"] = commit
	}

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

	// capture before query parameter if present, default to now
	before, err := strconv.ParseInt(c.DefaultQuery("before", strconv.FormatInt(time.Now().UTC().Unix(), 10)), 10, 64)
	if err != nil {
		retErr := fmt.Errorf("unable to convert before query parameter for repo %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// capture after query parameter if present, default to 0
	after, err := strconv.ParseInt(c.DefaultQuery("after", "0"), 10, 64)
	if err != nil {
		retErr := fmt.Errorf("unable to convert after query parameter for repo %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	b, t, err = database.FromContext(c).GetRepoBuildList(r, filters, before, after, page, perPage)
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
	// variables that will hold the build list, build list filters and total count
	var (
		filters = map[string]interface{}{}
		b       []*library.Build
		t       int64
	)

	// capture middleware values
	o := org.Retrieve(c)
	u := user.Retrieve(c)

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"org":  o,
		"user": u.GetName(),
	}).Infof("reading builds for org %s", o)

	// capture the branch name parameter
	branch := c.Query("branch")
	// capture the event type parameter
	event := c.Query("event")
	// capture the status type parameter
	status := c.Query("status")

	// check if branch filter was provided
	if len(branch) > 0 {
		// add branch to filters map
		filters["branch"] = branch
	}
	// check if event filter was provided
	if len(event) > 0 {
		// verify the event provided is a valid event type
		if event != constants.EventComment && event != constants.EventDeploy &&
			event != constants.EventPush && event != constants.EventPull &&
			event != constants.EventTag {
			retErr := fmt.Errorf("unable to process event %s: invalid event type provided", event)

			util.HandleError(c, http.StatusBadRequest, retErr)

			return
		}

		// add event to filters map
		filters["event"] = event
	}
	// check if status filter was provided
	if len(status) > 0 {
		// verify the status provided is a valid status type
		if status != constants.StatusCanceled && status != constants.StatusError &&
			status != constants.StatusFailure && status != constants.StatusKilled &&
			status != constants.StatusPending && status != constants.StatusRunning &&
			status != constants.StatusSuccess {
			retErr := fmt.Errorf("unable to process status %s: invalid status type provided", status)

			util.HandleError(c, http.StatusBadRequest, retErr)

			return
		}

		// add status to filters map
		filters["status"] = status
	}

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
	perPage = util.MaxInt(1, util.MinInt(100, perPage))

	// See if the user is an org admin to bypass individual permission checks
	perm, err := scm.FromContext(c).OrgAccess(u, o)
	if err != nil {
		logrus.Errorf("unable to get user %s access level for org %s", u.GetName(), o)
	}
	// Only show public repos to non-admins
	//nolint:goconst // ignore need for constant
	if perm != "admin" {
		filters["visibility"] = constants.VisibilityPublic
	}

	// send API call to capture the list of builds for the org (and event type if passed in)
	b, t, err = database.FromContext(c).GetOrgBuildList(o, filters, page, perPage)

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
//   description: Build number to retrieve
//   required: true
//   type: integer
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved the build
//     type: json
//     schema:
//       "$ref": "#/definitions/Build"

// GetBuild represents the API handler to capture
// a build for a repo from the configured backend.
func GetBuild(c *gin.Context) {
	// capture middleware values
	b := build.Retrieve(c)
	o := org.Retrieve(c)
	r := repo.Retrieve(c)
	u := user.Retrieve(c)

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"build": b.GetNumber(),
		"org":   o,
		"repo":  r.GetName(),
		"user":  u.GetName(),
	}).Infof("reading build %s/%d", r.GetFullName(), b.GetNumber())

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
//   '200':
//     description: Request processed but build was skipped
//     schema:
//       type: string
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

// RestartBuild represents the API handler to restart an existing build in the configured backend.
//
//nolint:funlen // ignore statement count
func RestartBuild(c *gin.Context) {
	// capture middleware values
	m := c.MustGet("metadata").(*types.Metadata)
	b := build.Retrieve(c)
	o := org.Retrieve(c)
	r := repo.Retrieve(c)
	u := user.Retrieve(c)

	entry := fmt.Sprintf("%s/%d", r.GetFullName(), b.GetNumber())

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logger := logrus.WithFields(logrus.Fields{
		"build": b.GetNumber(),
		"org":   o,
		"repo":  r.GetName(),
		"user":  u.GetName(),
	})

	logger.Infof("restarting build %s", entry)

	// send API call to capture the repo owner
	u, err := database.FromContext(c).GetUser(r.GetUserID())
	if err != nil {
		retErr := fmt.Errorf("unable to get owner for %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// create SQL filters for querying pending and running builds for repo
	filters := map[string]interface{}{
		"status": []string{constants.StatusPending, constants.StatusRunning},
	}

	// send API call to capture the number of pending or running builds for the repo
	builds, err := database.FromContext(c).GetRepoBuildCount(r, filters)
	if err != nil {
		retErr := fmt.Errorf("unable to restart build: unable to get count of builds for repo %s", r.GetFullName())

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// check if the number of pending and running builds exceeds the limit for the repo
	if builds >= r.GetBuildLimit() {
		retErr := fmt.Errorf("unable to restart build: repo %s has exceeded the concurrent build limit of %d", r.GetFullName(), r.GetBuildLimit())

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// update fields in build object
	b.SetID(0)
	b.SetCreated(time.Now().UTC().Unix())
	b.SetEnqueued(0)
	b.SetStarted(0)
	b.SetFinished(0)
	b.SetStatus(constants.StatusPending)
	b.SetHost("")
	b.SetRuntime("")
	b.SetDistribution("")

	// update the PR event action if action was never set
	// for backwards compatibility with pre-0.14 releases.
	if b.GetEvent() == constants.EventPull && b.GetEventAction() == "" {
		// technically, the action could have been opened or synchronize.
		// will not affect behavior of the pipeline since we did not
		// support actions for builds where this would be the case.
		b.SetEventAction(constants.ActionOpened)
	}

	// set the parent equal to the restarted build number
	b.SetParent(b.GetNumber())
	// update the build numbers based off repo counter
	inc := r.GetCounter() + 1
	r.SetCounter(inc)
	b.SetNumber(inc)

	// populate the build link if a web address is provided
	if len(m.Vela.WebAddress) > 0 {
		b.SetLink(
			fmt.Sprintf("%s/%s/%d", m.Vela.WebAddress, r.GetFullName(), b.GetNumber()),
		)
	}

	// variable to store changeset files
	var files []string
	// check if the build event is not issue_comment or pull_request
	if !strings.EqualFold(b.GetEvent(), constants.EventComment) &&
		!strings.EqualFold(b.GetEvent(), constants.EventPull) {
		// send API call to capture list of files changed for the commit
		files, err = scm.FromContext(c).Changeset(u, r, b.GetCommit())
		if err != nil {
			retErr := fmt.Errorf("unable to restart build: failed to get changeset for %s: %w", r.GetFullName(), err)

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}
	}

	// check if the build event is a pull_request
	if strings.EqualFold(b.GetEvent(), constants.EventPull) {
		// capture number from build
		number, err := getPRNumberFromBuild(b)
		if err != nil {
			retErr := fmt.Errorf("unable to restart build: failed to get pull_request number for %s: %w", r.GetFullName(), err)

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}

		// send API call to capture list of files changed for the pull request
		files, err = scm.FromContext(c).ChangesetPR(u, r, number)
		if err != nil {
			retErr := fmt.Errorf("unable to restart build: failed to get changeset for %s: %w", r.GetFullName(), err)

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}
	}

	// variables to store pipeline configuration
	var (
		// variable to store the raw pipeline configuration
		config []byte
		// variable to store executable pipeline
		p *pipeline.Build
		// variable to store pipeline configuration
		pipeline *library.Pipeline
		// variable to store the pipeline type for the repository
		pipelineType = r.GetPipelineType()
	)

	// send API call to attempt to capture the pipeline
	pipeline, err = database.FromContext(c).GetPipelineForRepo(b.GetCommit(), r)
	if err != nil { // assume the pipeline doesn't exist in the database yet (before pipeline support was added)
		// send API call to capture the pipeline configuration file
		config, err = scm.FromContext(c).ConfigBackoff(u, r, b.GetCommit())
		if err != nil {
			retErr := fmt.Errorf("unable to get pipeline configuration for %s: %w", r.GetFullName(), err)

			util.HandleError(c, http.StatusNotFound, retErr)

			return
		}
	} else {
		config = pipeline.GetData()
	}

	// ensure we use the expected pipeline type when compiling
	//
	// The pipeline type for a repo can change at any time which can break compiling
	// existing pipelines in the system for that repo. To account for this, we update
	// the repo pipeline type to match what was defined for the existing pipeline
	// before compiling. After we're done compiling, we reset the pipeline type.
	if len(pipeline.GetType()) > 0 {
		r.SetPipelineType(pipeline.GetType())
	}

	var compiled *library.Pipeline
	// parse and compile the pipeline configuration file
	p, compiled, err = compiler.FromContext(c).
		Duplicate().
		WithBuild(b).
		WithFiles(files).
		WithMetadata(m).
		WithRepo(r).
		WithUser(u).
		Compile(config)
	if err != nil {
		retErr := fmt.Errorf("unable to compile pipeline configuration for %s: %w", entry, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}
	// reset the pipeline type for the repo
	//
	// The pipeline type for a repo can change at any time which can break compiling
	// existing pipelines in the system for that repo. To account for this, we update
	// the repo pipeline type to match what was defined for the existing pipeline
	// before compiling. After we're done compiling, we reset the pipeline type.
	r.SetPipelineType(pipelineType)

	// skip the build if only the init or clone steps are found
	skip := skipEmptyBuild(p)
	if skip != "" {
		// set build to successful status
		b.SetStatus(constants.StatusSkipped)

		// send API call to set the status on the commit
		err = scm.FromContext(c).Status(u, b, r.GetOrg(), r.GetName())
		if err != nil {
			logrus.Errorf("unable to set commit status for %s/%d: %v", r.GetFullName(), b.GetNumber(), err)
		}

		c.JSON(http.StatusOK, skip)

		return
	}

	// check if the pipeline did not already exist in the database
	//
	//nolint:dupl // ignore duplicate code
	if pipeline == nil {
		pipeline = compiled
		pipeline.SetRepoID(r.GetID())
		pipeline.SetCommit(b.GetCommit())
		pipeline.SetRef(b.GetRef())

		// send API call to create the pipeline
		err = database.FromContext(c).CreatePipeline(pipeline)
		if err != nil {
			retErr := fmt.Errorf("unable to create pipeline for %s: %w", r.GetFullName(), err)

			util.HandleError(c, http.StatusBadRequest, retErr)

			return
		}

		// send API call to capture the created pipeline
		pipeline, err = database.FromContext(c).GetPipelineForRepo(pipeline.GetCommit(), r)
		if err != nil {
			//nolint:lll // ignore long line length due to error message
			retErr := fmt.Errorf("unable to get new pipeline %s/%s: %w", r.GetFullName(), pipeline.GetCommit(), err)

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}
	}

	b.SetPipelineID(pipeline.GetID())

	// create the objects from the pipeline in the database
	err = planBuild(database.FromContext(c), p, b, r)
	if err != nil {
		util.HandleError(c, http.StatusInternalServerError, err)

		return
	}

	// send API call to update repo for ensuring counter is incremented
	err = database.FromContext(c).UpdateRepo(r)
	if err != nil {
		retErr := fmt.Errorf("unable to restart build: failed to update repo %s: %w", r.GetFullName(), err)
		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// send API call to capture the restarted build
	b, _ = database.FromContext(c).GetBuild(b.GetNumber(), r)

	c.JSON(http.StatusCreated, b)

	// send API call to set the status on the commit
	err = scm.FromContext(c).Status(u, b, r.GetOrg(), r.GetName())
	if err != nil {
		logger.Errorf("unable to set commit status for build %s: %v", entry, err)
	}

	// publish the build to the queue
	go publishToQueue(
		queue.FromGinContext(c),
		database.FromContext(c),
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
	o := org.Retrieve(c)
	r := repo.Retrieve(c)
	u := user.Retrieve(c)

	entry := fmt.Sprintf("%s/%d", r.GetFullName(), b.GetNumber())

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"build": b.GetNumber(),
		"org":   o,
		"repo":  r.GetName(),
		"user":  u.GetName(),
	}).Infof("updating build %s", entry)

	// capture body from API request
	input := new(library.Build)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for build %s: %w", entry, err)

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
	b, err = database.FromContext(c).UpdateBuild(b)
	if err != nil {
		retErr := fmt.Errorf("unable to update build %s: %w", entry, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

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
			logrus.Errorf("unable to get owner for build %s: %v", entry, err)
		}

		// send API call to set the status on the commit
		err = scm.FromContext(c).Status(u, b, r.GetOrg(), r.GetName())
		if err != nil {
			logrus.Errorf("unable to set commit status for build %s: %v", entry, err)
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
	o := org.Retrieve(c)
	r := repo.Retrieve(c)
	u := user.Retrieve(c)

	entry := fmt.Sprintf("%s/%d", r.GetFullName(), b.GetNumber())

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"build": b.GetNumber(),
		"org":   o,
		"repo":  r.GetName(),
		"user":  u.GetName(),
	}).Infof("deleting build %s", entry)

	// send API call to remove the build
	err := database.FromContext(c).DeleteBuild(b.GetID())
	if err != nil {
		retErr := fmt.Errorf("unable to delete build %s: %w", entry, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("build %s deleted", entry))
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
	if len(parts) < 3 {
		return 0, fmt.Errorf("invalid ref: %s", b.GetRef())
	}

	// return the results of converting number to string
	return strconv.Atoi(parts[2])
}

// planBuild is a helper function to plan the build for
// execution. This creates all resources, like steps
// and services, for the build in the configured backend.
// TODO:
// - return build and error.
func planBuild(database database.Service, p *pipeline.Build, b *library.Build, r *library.Repo) error {
	// update fields in build object
	b.SetCreated(time.Now().UTC().Unix())

	// send API call to create the build
	b, err := database.CreateBuild(b)
	if err != nil {
		// clean up the objects from the pipeline in the database
		// TODO:
		// - return build in CreateBuild
		// - even if it was created, we need to get the new build id
		//   otherwise it will be 0, which attempts to INSERT instead
		//   of UPDATE-ing the existing build - which results in
		//   a constraint error (repo_id, number)
		// - do we want to update the build or just delete it?
		cleanBuild(database, b, nil, nil, err)

		return fmt.Errorf("unable to create new build for %s: %w", r.GetFullName(), err)
	}

	// plan all services for the build
	services, err := planServices(database, p, b)
	if err != nil {
		// clean up the objects from the pipeline in the database
		cleanBuild(database, b, services, nil, err)

		return err
	}

	// plan all steps for the build
	steps, err := planSteps(database, p, b)
	if err != nil {
		// clean up the objects from the pipeline in the database
		cleanBuild(database, b, services, steps, err)

		return err
	}

	return nil
}

// cleanBuild is a helper function to kill the build
// without execution. This will kill all resources,
// like steps and services, for the build in the
// configured backend.
func cleanBuild(database database.Service, b *library.Build, services []*library.Service, steps []*library.Step, e error) {
	// update fields in build object
	b.SetError(fmt.Sprintf("unable to publish to queue: %s", e.Error()))
	b.SetStatus(constants.StatusError)
	b.SetFinished(time.Now().UTC().Unix())

	// send API call to update the build
	_, err := database.UpdateBuild(b)
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

// CancelBuild represents the API handler to cancel a running build.
//
//nolint:funlen // ignore statement count
func CancelBuild(c *gin.Context) {
	// capture middleware values
	b := build.Retrieve(c)
	e := executors.Retrieve(c)
	o := org.Retrieve(c)
	r := repo.Retrieve(c)
	u := user.Retrieve(c)

	entry := fmt.Sprintf("%s/%d", r.GetFullName(), b.GetNumber())

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"build": b.GetNumber(),
		"org":   o,
		"repo":  r.GetName(),
		"user":  u.GetName(),
	}).Infof("canceling build %s", entry)

	// TODO: add support for removing builds from the queue
	//
	// check to see if build is not running
	if !strings.EqualFold(b.GetStatus(), constants.StatusRunning) {
		retErr := fmt.Errorf("found build %s but its status was %s", entry, b.GetStatus())

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// retrieve the worker info
	w, err := database.FromContext(c).GetWorkerForHostname(b.GetHost())
	if err != nil {
		retErr := fmt.Errorf("unable to get worker for build %s: %w", entry, err)
		util.HandleError(c, http.StatusNotFound, retErr)

		return
	}

	for _, executor := range e {
		// check each executor on the worker running the build to see if it's running the build we want to cancel
		if strings.EqualFold(executor.Repo.GetFullName(), r.GetFullName()) && *executor.GetBuild().Number == b.GetNumber() {
			// prepare the request to the worker
			client := http.DefaultClient
			client.Timeout = 30 * time.Second

			// set the API endpoint path we send the request to
			u := fmt.Sprintf("%s/api/v1/executors/%d/build/cancel", w.GetAddress(), executor.GetID())

			req, err := http.NewRequestWithContext(context.Background(), "DELETE", u, nil)
			if err != nil {
				retErr := fmt.Errorf("unable to form a request to %s: %w", u, err)
				util.HandleError(c, http.StatusBadRequest, retErr)

				return
			}

			tm := c.MustGet("token-manager").(*token.Manager)

			// set mint token options
			mto := &token.MintTokenOpts{
				Hostname:      "vela-server",
				TokenType:     constants.WorkerAuthTokenType,
				TokenDuration: time.Minute * 1,
			}

			// mint token
			tkn, err := tm.MintToken(mto)
			if err != nil {
				retErr := fmt.Errorf("unable to generate auth token: %w", err)
				util.HandleError(c, http.StatusInternalServerError, retErr)

				return
			}

			// add the token to authenticate to the worker
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tkn))

			// perform the request to the worker
			resp, err := client.Do(req)
			if err != nil {
				retErr := fmt.Errorf("unable to connect to %s: %w", u, err)
				util.HandleError(c, http.StatusBadRequest, retErr)

				return
			}
			defer resp.Body.Close()

			// Read Response Body
			respBody, err := io.ReadAll(resp.Body)
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
	// update the status in the build table
	b.SetStatus(constants.StatusCanceled)

	_, err = database.FromContext(c).UpdateBuild(b)
	if err != nil {
		retErr := fmt.Errorf("unable to update status for build %s: %w", entry, err)
		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// retrieve the steps for the build from the step table
	steps := []*library.Step{}
	page := 1
	perPage := 100

	for page > 0 {
		// retrieve build steps (per page) from the database
		stepsPart, err := database.FromContext(c).GetBuildStepList(b, page, perPage)
		if err != nil {
			retErr := fmt.Errorf("unable to retrieve steps for build %s: %w", entry, err)
			util.HandleError(c, http.StatusNotFound, retErr)

			return
		}

		// add page of steps to list steps
		steps = append(steps, stepsPart...)

		// assume no more pages exist if under 100 results are returned
		if len(stepsPart) < 100 {
			page = 0
		} else {
			page++
		}
	}

	// iterate over each step for the build
	// setting anything running or pending to canceled
	for _, step := range steps {
		if step.GetStatus() == constants.StatusRunning || step.GetStatus() == constants.StatusPending {
			step.SetStatus(constants.StatusCanceled)

			err = database.FromContext(c).UpdateStep(step)
			if err != nil {
				retErr := fmt.Errorf("unable to update step %s for build %s: %w", step.GetName(), entry, err)
				util.HandleError(c, http.StatusNotFound, retErr)

				return
			}
		}
	}

	// retrieve the services for the build from the service table
	services := []*library.Service{}
	page = 1

	for page > 0 {
		// retrieve build services (per page) from the database
		servicesPart, err := database.FromContext(c).GetBuildServiceList(b, page, perPage)
		if err != nil {
			retErr := fmt.Errorf("unable to retrieve services for build %s: %w", entry, err)
			util.HandleError(c, http.StatusNotFound, retErr)

			return
		}

		// add page of services to the list of services
		services = append(services, servicesPart...)

		// assume no more pages exist if under 100 results are returned
		if len(servicesPart) < 100 {
			page = 0
		} else {
			page++
		}
	}

	// iterate over each service for the build
	// setting anything running or pending to canceled
	for _, service := range services {
		if service.GetStatus() == constants.StatusRunning || service.GetStatus() == constants.StatusPending {
			service.SetStatus(constants.StatusCanceled)

			err = database.FromContext(c).UpdateService(service)
			if err != nil {
				retErr := fmt.Errorf("unable to update service %s for build %s: %w",
					service.GetName(),
					entry,
					err,
				)
				util.HandleError(c, http.StatusNotFound, retErr)

				return
			}
		}
	}

	c.JSON(http.StatusOK, b)
}

// swagger:operation GET /api/v1/repos/{org}/{repo}/builds/{build}/token builds GetBuildToken
//
// Get a build token
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
//   description: Build number
//   required: true
//   type: integer
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved build token
//     schema:
//       "$ref": "#/definitions/Token"
//   '400':
//     description: Bad request
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to generate build token
//     schema:
//       "$ref": "#/definitions/Error"

// GetBuildToken represents the API handler to generate a build token.
func GetBuildToken(c *gin.Context) {
	// capture middleware values
	b := build.Retrieve(c)
	o := org.Retrieve(c)
	r := repo.Retrieve(c)
	cl := claims.Retrieve(c)

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"build": b.GetNumber(),
		"org":   o,
		"repo":  r.GetName(),
		"user":  cl.Subject,
	}).Infof("generating build token for build %s/%d", r.GetFullName(), b.GetNumber())

	// if build is not in a pending state, then a build token should not be needed - bad request
	if !strings.EqualFold(b.GetStatus(), constants.StatusPending) {
		retErr := fmt.Errorf("unable to mint build token: build is not in pending state")
		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// retrieve token manager from context
	tm := c.MustGet("token-manager").(*token.Manager)

	// set expiration to repo timeout plus configurable buffer
	exp := (time.Duration(r.GetTimeout()) * time.Minute) + tm.BuildTokenBufferDuration

	// set mint token options
	bmto := &token.MintTokenOpts{
		Hostname:      cl.Subject,
		BuildID:       b.GetID(),
		Repo:          r.GetFullName(),
		TokenType:     constants.WorkerBuildTokenType,
		TokenDuration: exp,
	}

	// mint token
	bt, err := tm.MintToken(bmto)
	if err != nil {
		retErr := fmt.Errorf("unable to generate build token: %w", err)
		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, library.Token{Token: &bt})
}
