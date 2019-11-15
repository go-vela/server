// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-vela/compiler/compiler"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/queue"
	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/source"
	"github.com/go-vela/server/util"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// CreateBuild represents the API handler to
// create a build in the configured backend.
func CreateBuild(c *gin.Context) {
	// capture middleware values
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

	// send API call to capture list of files changed for the commit
	files, err := source.FromContext(c).ListChanges(u, r, input.GetCommit())
	if err != nil {
		retErr := fmt.Errorf("unable to get changeset for %s/%d: %w", r.GetFullName(), input.GetNumber(), err)
		util.HandleError(c, http.StatusInternalServerError, retErr)
		return
	}

	// send API call to capture the pipeline configuration file
	config, err := source.FromContext(c).Config(u, r.GetOrg(), r.GetName(), input.GetCommit())
	if err != nil {
		retErr := fmt.Errorf("unable to get pipeline configuration for %s/%d: %w", r.GetFullName(), input.GetNumber(), err)
		util.HandleError(c, http.StatusNotFound, retErr)
		return
	}

	// parse and compile the pipeline configuration file
	pipe, err := compiler.FromContext(c).
		WithBuild(input).
		WithFiles(files).
		WithRepo(r).
		WithUser(u).
		Compile(config)
	if err != nil {
		retErr := fmt.Errorf("unable to compile pipeline configuration for %s/%d: %w", r.GetFullName(), input.GetNumber(), err)
		util.HandleError(c, http.StatusInternalServerError, retErr)
		return
	}

	// create the objects from the pipeline in the database
	err = planBuild(database.FromContext(c), pipe, input, r)
	if err != nil {
		util.HandleError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, input)

	// send API call to set the status on the commit
	err = source.FromContext(c).Status(u, input, r.GetOrg(), r.GetName())
	if err != nil {
		logrus.Errorf("unable to set commit status for build %s/%d: %w", r.GetFullName(), input.GetNumber(), err)
	}

	// publish the build to the queue
	go publishToQueue(
		queue.FromContext(c),
		pipe,
		input,
		r,
		u,
	)
}

// GetBuilds represents the API handler to capture a
// list of builds for a repo from the configured backend.
func GetBuilds(c *gin.Context) {
	// capture middleware values
	r := repo.Retrieve(c)

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

	// send API call to capture the total number of builds for the repo
	t, err := database.FromContext(c).GetRepoBuildCount(r)
	if err != nil {
		retErr := fmt.Errorf("unable to get build count for repo %s: %w", r.GetFullName(), err)
		util.HandleError(c, http.StatusInternalServerError, retErr)
		return
	}

	// send API call to capture the list of builds for the repo
	b, err := database.FromContext(c).GetRepoBuildList(r, page, perPage)
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

// RestartBuild represents the API handler to
// restart an existing build in the configured backend.
func RestartBuild(c *gin.Context) {
	// capture middleware values
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

	// send API call to capture list of files changed for the commit
	files, err := source.FromContext(c).ListChanges(u, r, lastBuild.GetCommit())
	if err != nil {
		retErr := fmt.Errorf("unable to get changeset for %s/%d: %w", r.GetFullName(), b.GetNumber(), err)
		util.HandleError(c, http.StatusInternalServerError, retErr)
		return
	}

	// send API call to capture the pipeline configuration file
	config, err := source.FromContext(c).Config(u, r.GetOrg(), r.GetName(), b.GetCommit())
	if err != nil {
		retErr := fmt.Errorf("unable to get pipeline configuration for %s/%d: %w", r.GetFullName(), b.GetNumber(), err)
		util.HandleError(c, http.StatusNotFound, retErr)
		return
	}

	// parse and compile the pipeline configuration file
	pipe, err := compiler.FromContext(c).
		WithBuild(b).
		WithFiles(files).
		WithRepo(r).
		WithUser(u).
		Compile(config)
	if err != nil {
		retErr := fmt.Errorf("unable to compile pipeline configuration for %s/%d: %w", r.GetFullName(), b.GetNumber(), err)
		util.HandleError(c, http.StatusInternalServerError, retErr)
		return
	}

	// create the objects from the pipeline in the database
	err = planBuild(database.FromContext(c), pipe, b, r)
	if err != nil {
		util.HandleError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, b)

	// send API call to set the status on the commit
	err = source.FromContext(c).Status(u, b, r.GetOrg(), r.GetName())
	if err != nil {
		logrus.Errorf("unable to set commit status for build %s/%d: %w", r.GetFullName(), b.GetNumber(), err)
	}

	// publish the build to the queue
	go publishToQueue(
		queue.FromContext(c),
		pipe,
		b,
		r,
		u,
	)
}

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
			logrus.Errorf("unable to get owner for %s/%d: %w", r.GetFullName(), b.GetNumber(), err)
		}

		// send API call to set the status on the commit
		err = source.FromContext(c).Status(u, b, r.GetOrg(), r.GetName())
		if err != nil {
			logrus.Errorf("unable to set commit status for build %s/%d: %w", r.GetFullName(), b.GetNumber(), err)
		}
	}
}

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
		retErr := fmt.Errorf("unable to delete build %s/%d: %w", r.GetFullName(), b.GetNumber(), err)
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
		return fmt.Errorf("unable to create new build for %s: %w", r.GetFullName(), err)
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
		logrus.Errorf("unable to kill build %d: %w", b.GetNumber(), err)
	}

	for _, s := range services {
		// update fields in service object
		s.SetStatus(constants.StatusKilled)
		s.SetFinished(time.Now().UTC().Unix())

		// send API call to update the service
		err := database.UpdateService(s)
		if err != nil {
			logrus.Errorf("unable to kill service %s for build %d: %w", s.GetName(), b.GetNumber(), err)
		}
	}

	for _, s := range steps {
		// update fields in step object
		s.SetStatus(constants.StatusKilled)
		s.SetFinished(time.Now().UTC().Unix())

		// send API call to update the step
		err := database.UpdateStep(s)
		if err != nil {
			logrus.Errorf("unable to kill step %s for build %d: %w", s.GetName(), b.GetNumber(), err)
		}
	}
}
