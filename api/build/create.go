// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package build

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/compiler"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/queue"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"
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
//   description: Payload containing the build to create
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
	builds, err := database.FromContext(c).CountBuildsForRepo(r, filters)
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
	skip := SkipEmptyBuild(p)
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
	err = PlanBuild(database.FromContext(c), p, input, r)
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
	input, _ = database.FromContext(c).GetBuildForRepo(r, input.GetNumber())

	c.JSON(http.StatusCreated, input)

	// send API call to set the status on the commit
	err = scm.FromContext(c).Status(u, input, r.GetOrg(), r.GetName())
	if err != nil {
		logger.Errorf("unable to set commit status for build %s/%d: %v", r.GetFullName(), input.GetNumber(), err)
	}

	// publish the build to the queue
	go PublishToQueue(
		queue.FromGinContext(c),
		database.FromContext(c),
		p,
		input,
		r,
		u,
	)
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
