// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-vela/server/api"
	"github.com/go-vela/server/compiler"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/types"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"
	"github.com/sirupsen/logrus"
)

const baseErr = "unable to schedule build"

func processSchedules(compiler compiler.Engine, database database.Interface, metadata *types.Metadata, scm scm.Service) error {
	// send API call to capture the list of active schedules
	schedules, err := database.ListActiveSchedules()
	if err != nil {
		return err
	}

	// iterate through the list of active schedules
	for _, schedule := range schedules {
		err = processSchedule(schedule, compiler, database, metadata, scm)
		if err != nil {
			logrus.WithError(err).Warnf("%s for %s", baseErr, schedule.GetName())

			continue
		}
	}

	return nil
}

func processSchedule(s *library.Schedule, compiler compiler.Engine, database database.Interface, metadata *types.Metadata, scm scm.Service) error {
	// send API call to capture the repo for the schedule
	r, err := database.GetRepo(s.GetRepoID())
	if err != nil {
		return fmt.Errorf("unable to fetch repo: %w", err)
	}

	// check if the repo is active
	if !r.GetActive() {
		return fmt.Errorf("repo %s is not active", r.GetFullName())
	}

	// check if the repo allows the schedule event type
	if !r.GetAllowSchedule() {
		return fmt.Errorf("repo %s does not have %s events enabled", r.GetFullName(), constants.EventSchedule)
	}

	// check if the repo has a valid owner
	if r.GetUserID() == 0 {
		return fmt.Errorf("repo %s does not have a valid owner", r.GetFullName())
	}

	// send API call to capture the owner for the repo
	user, err := database.GetUser(r.GetUserID())
	if err != nil {
		return fmt.Errorf("unable to get owner for repo %s: %w", r.GetFullName(), err)
	}

	// send API call to confirm repo owner has at least write access to repo
	_, err = scm.RepoAccess(user, user.GetToken(), r.GetOrg(), r.GetName())
	if err != nil {
		return fmt.Errorf("%s does not have at least write access for repo %s", user.GetName(), r.GetFullName())
	}

	// create SQL filters for querying pending and running builds for repo
	filters := map[string]interface{}{
		"status": []string{constants.StatusPending, constants.StatusRunning},
	}

	// send API call to capture the number of pending or running builds for the repo
	builds, err := database.GetRepoBuildCount(r, filters)
	if err != nil {
		return fmt.Errorf("unable to get count of builds for repo %s: %w", r.GetFullName(), err)
	}

	// check if the number of pending and running builds exceeds the limit for the repo
	if builds >= r.GetBuildLimit() {
		return fmt.Errorf("repo %s has excceded the concurrent build limit of %d", r.GetFullName(), r.GetBuildLimit())
	}

	url := strings.TrimSuffix(r.GetClone(), ".git")

	b := new(library.Build)
	b.SetAuthor(s.GetCreatedBy())
	b.SetBranch(r.GetBranch())
	b.SetClone(r.GetClone())
	b.SetCommit(r.GetBranch())
	b.SetDeploy(s.GetName())
	b.SetEvent(constants.EventSchedule)
	b.SetMessage(fmt.Sprintf("triggered for %s schedule with %s entry", s.GetName(), s.GetEntry()))
	b.SetRef(fmt.Sprintf("refs/heads/%s", b.GetBranch()))
	b.SetRepoID(r.GetID())
	b.SetSender(s.GetUpdatedBy())
	b.SetSource(fmt.Sprintf("%s/tree/%s", url, b.GetBranch()))
	b.SetStatus(constants.StatusPending)
	b.SetTitle(fmt.Sprintf("%s received from %s", constants.EventSchedule, url))

	// populate the build link if a web address is provided
	if len(metadata.Vela.WebAddress) > 0 {
		b.SetLink(fmt.Sprintf("%s/%s/%d", metadata.Vela.WebAddress, r.GetFullName(), b.GetNumber()))
	}

	var (
		// variable to store the raw pipeline configuration
		config []byte
		// variable to store executable pipeline
		p *pipeline.Build
		// variable to store pipeline configuration
		pipeline *library.Pipeline
		// variable to control number of times to retry processing pipeline
		retryLimit = 5
		// variable to store the pipeline type for the repository
		pipelineType = r.GetPipelineType()
	)

	// implement a loop to process asynchronous operations with a retry limit
	//
	// Some operations taken during this workflow can lead to race conditions failing to successfully process
	// the request. This logic ensures we attempt our best efforts to handle these cases gracefully.
	for i := 0; i < retryLimit; i++ {
		logrus.Debugf("compilation loop - attempt %d", i+1)
		// check if we're on the first iteration of the loop
		if i > 0 {
			// incrementally sleep in between retries
			time.Sleep(time.Duration(i) * time.Second)
		}

		// send API call to attempt to capture the pipeline
		pipeline, err = database.GetPipelineForRepo(b.GetCommit(), r)
		if err != nil { // assume the pipeline doesn't exist in the database yet
			// send API call to capture the pipeline configuration file
			config, err = scm.ConfigBackoff(user, r, b.GetCommit())
			if err != nil {
				return fmt.Errorf("unable to get pipeline config for %s/%s: %w", r.GetFullName(), b.GetCommit(), err)
			}
		} else {
			config = pipeline.GetData()
		}

		// send API call to capture repo for the counter (grabbing repo again to ensure counter is correct)
		repo, err := database.GetRepoForOrg(r.GetOrg(), r.GetName())
		if err != nil {
			err = fmt.Errorf("unable to get repo %s: %w", r.GetFullName(), err)

			// check if the retry limit has been exceeded
			if i < retryLimit-1 {
				logrus.WithError(err).Warningf("retrying #%d", i+1)

				// continue to the next iteration of the loop
				continue
			}

			return err
		}

		// update repo fields with any changes from SCM process
		repo.SetTopics(r.GetTopics())
		repo.SetBranch(r.GetBranch())

		// set the build numbers based off repo counter
		repo.SetCounter(repo.GetCounter() + 1)
		b.SetNumber(repo.GetCounter() + 1)
		// set the parent equal to the current repo counter
		b.SetParent(repo.GetCounter())
		// check if the parent is set to 0
		if b.GetParent() == 0 {
			// parent should be "1" if it's the first build ran
			b.SetParent(1)
		}

		// set the build link if a web address is provided
		if len(metadata.Vela.WebAddress) > 0 {
			b.SetLink(fmt.Sprintf("%s/%s/%d", metadata.Vela.WebAddress, repo.GetFullName(), b.GetNumber()))
		}

		// ensure we use the expected pipeline type when compiling
		//
		// The pipeline type for a repo can change at any time which can break compiling
		// existing pipelines in the system for that repo. To account for this, we update
		// the repo pipeline type to match what was defined for the existing pipeline
		// before compiling. After we're done compiling, we reset the pipeline type.
		if len(pipeline.GetType()) > 0 {
			repo.SetPipelineType(pipeline.GetType())
		}

		var compiled *library.Pipeline
		// parse and compile the pipeline configuration file
		p, compiled, err = compiler.
			Duplicate().
			WithBuild(b).
			WithMetadata(metadata).
			WithRepo(repo).
			WithUser(user).
			Compile(config)
		if err != nil {
			return fmt.Errorf("unable to compile pipeline config for %s/%s: %w", r.GetFullName(), b.GetCommit(), err)
		}

		// reset the pipeline type for the repo
		//
		// The pipeline type for a repo can change at any time which can break compiling
		// existing pipelines in the system for that repo. To account for this, we update
		// the repo pipeline type to match what was defined for the existing pipeline
		// before compiling. After we're done compiling, we reset the pipeline type.
		repo.SetPipelineType(pipelineType)

		// skip the build if only the init or clone steps are found
		skip := api.SkipEmptyBuild(p)
		if skip != "" {
			return nil
		}

		// check if the pipeline did not already exist in the database
		if pipeline == nil {
			pipeline = compiled
			pipeline.SetRepoID(repo.GetID())
			pipeline.SetCommit(b.GetCommit())
			pipeline.SetRef(b.GetRef())

			// send API call to create the pipeline
			err = database.CreatePipeline(pipeline)
			if err != nil {
				err = fmt.Errorf("failed to create pipeline for %s: %w", repo.GetFullName(), err)

				// check if the retry limit has been exceeded
				if i < retryLimit-1 {
					logrus.WithError(err).Warningf("retrying #%d", i+1)

					// continue to the next iteration of the loop
					continue
				}

				return err
			}

			// send API call to capture the created pipeline
			pipeline, err = database.GetPipelineForRepo(pipeline.GetCommit(), repo)
			if err != nil {
				return fmt.Errorf("unable to get new pipeline %s/%s: %w", repo.GetFullName(), pipeline.GetCommit(), err)
			}
		}

		b.SetPipelineID(pipeline.GetID())

		// create the objects from the pipeline in the database
		// TODO:
		// - if a build gets created and something else fails midway,
		//   the next loop will attempt to create the same build,
		//   using the same Number and thus create a constraint
		//   conflict; consider deleting the partially created
		//   build object in the database
		err = api.PlanBuild(database, p, b, repo)
		if err != nil {
			// check if the retry limit has been exceeded
			if i < retryLimit-1 {
				logrus.WithError(err).Warningf("retrying #%d", i+1)

				// reset fields set by cleanBuild for retry
				b.SetError("")
				b.SetStatus(constants.StatusPending)
				b.SetFinished(0)

				// continue to the next iteration of the loop
				continue
			}

			return err
		}

		// break the loop because everything was successful
		break
	} // end of retry loop

	return nil
}
