// SPDX-License-Identifier: Apache-2.0

package build

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/compiler"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/queue"
	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/claims"
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
	cl := claims.Retrieve(c)
	b := build.Retrieve(c)
	o := org.Retrieve(c)
	r := repo.Retrieve(c)
	u := user.Retrieve(c)
	ctx := c.Request.Context()

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

	if strings.EqualFold(b.GetStatus(), constants.StatusPendingApproval) {
		retErr := fmt.Errorf("unable to restart build %s/%d: cannot restart a build pending approval", r.GetFullName(), b.GetNumber())

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	logger.Infof("restarting build %s", entry)

	// send API call to capture the repo owner
	u, err := database.FromContext(c).GetUser(ctx, r.GetUserID())
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
	builds, err := database.FromContext(c).CountBuildsForRepo(ctx, r, filters)
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
	b.SetSender(cl.Subject)

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
		files, err = scm.FromContext(c).Changeset(ctx, u, r, b.GetCommit())
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
		files, err = scm.FromContext(c).ChangesetPR(ctx, u, r, number)
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
	pipeline, err = database.FromContext(c).GetPipelineForRepo(ctx, b.GetCommit(), r)
	if err != nil { // assume the pipeline doesn't exist in the database yet (before pipeline support was added)
		// send API call to capture the pipeline configuration file
		config, err = scm.FromContext(c).ConfigBackoff(ctx, u, r, b.GetCommit())
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
		WithCommit(b.GetCommit()).
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
	skip := SkipEmptyBuild(p)
	if skip != "" {
		// set build to successful status
		b.SetStatus(constants.StatusSkipped)

		// send API call to set the status on the commit
		err = scm.FromContext(c).Status(ctx, u, b, r.GetOrg(), r.GetName())
		if err != nil {
			logrus.Errorf("unable to set commit status for %s/%d: %v", r.GetFullName(), b.GetNumber(), err)
		}

		c.JSON(http.StatusOK, skip)

		return
	}

	// check if the pipeline did not already exist in the database
	if pipeline == nil {
		pipeline = compiled
		pipeline.SetRepoID(r.GetID())
		pipeline.SetCommit(b.GetCommit())
		pipeline.SetRef(b.GetRef())

		// send API call to create the pipeline
		pipeline, err = database.FromContext(c).CreatePipeline(ctx, pipeline)
		if err != nil {
			retErr := fmt.Errorf("unable to create pipeline for %s: %w", r.GetFullName(), err)

			util.HandleError(c, http.StatusBadRequest, retErr)

			return
		}
	}

	b.SetPipelineID(pipeline.GetID())

	// create the objects from the pipeline in the database
	err = PlanBuild(ctx, database.FromContext(c), p, b, r)
	if err != nil {
		util.HandleError(c, http.StatusInternalServerError, err)

		return
	}

	// send API call to update repo for ensuring counter is incremented
	r, err = database.FromContext(c).UpdateRepo(ctx, r)
	if err != nil {
		retErr := fmt.Errorf("unable to restart build: failed to update repo %s: %w", r.GetFullName(), err)
		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// send API call to capture the restarted build
	b, _ = database.FromContext(c).GetBuildForRepo(ctx, r, b.GetNumber())

	c.JSON(http.StatusCreated, b)

	// if the event is a deployment, update the build list
	if strings.EqualFold(b.GetEvent(), constants.EventDeploy) {
		d, err := database.FromContext(c).GetDeploymentForRepo(c, r, b.GetDeployNumber())
		if err != nil {
			logger.Errorf("unable to set get deployment for build %s: %v", entry, err)
		}

		build := append(d.GetBuilds(), b)

		d.SetBuilds(build)

		_, err = database.FromContext(c).UpdateDeployment(d)
		if err != nil {
			logger.Errorf("unable to set update deployment for build %s: %v", entry, err)
		}
	}

	// send API call to set the status on the commit
	err = scm.FromContext(c).Status(ctx, u, b, r.GetOrg(), r.GetName())
	if err != nil {
		logger.Errorf("unable to set commit status for build %s: %v", entry, err)
	}

	// determine queue route
	route, err := queue.FromContext(c).Route(&p.Worker)
	if err != nil {
		logrus.Errorf("unable to set route for build %d for %s: %v", b.GetNumber(), r.GetFullName(), err)

		// error out the build
		CleanBuild(ctx, database.FromContext(c), b, nil, nil, err)

		return
	}

	// temporarily set host to the route before it gets picked up by a worker
	b.SetHost(route)

	err = PublishBuildExecutable(ctx, database.FromContext(c), p, b)
	if err != nil {
		retErr := fmt.Errorf("unable to publish build executable for %s/%d: %w", r.GetFullName(), b.GetNumber(), err)
		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// publish the build to the queue
	go PublishToQueue(
		ctx,
		queue.FromGinContext(c),
		database.FromContext(c),
		b,
		r,
		u,
		route,
	)
}
