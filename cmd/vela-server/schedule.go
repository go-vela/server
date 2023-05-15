// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-vela/server/api"

	"github.com/go-vela/server/compiler"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"
	"github.com/sirupsen/logrus"
)

const baseErr = "unable to schedule build"

func scheduler(compiler compiler.Engine, database database.Service, metadata *types.Metadata, scm scm.Service) error {
	// send API call to capture the list of active schedules
	schedules, err := database.ListActiveSchedules()
	if err != nil {
		return err
	}

	// iterate through the list of active schedules
	for _, schedule := range schedules {
		// send API call to capture the repo for the schedule
		repo, err := database.GetRepo(schedule.GetRepo().GetID())
		if err != nil {
			logrus.WithError(err).Warnf("%s: unable to fetch repo", baseErr)

			continue
		}

		entry := fmt.Sprintf("%s/%s", repo.GetFullName(), schedule.GetName())

		// check if the repo is active
		if !repo.GetActive() {
			logrus.Warnf("%s for %s: repo is not active", baseErr, entry)

			continue
		}

		// check if the repo allows the schedule event type
		if !repo.GetAllowSchedule() {
			logrus.Warnf("%s for %s: repo does not have %s events enabled", baseErr, entry, constants.EventSchedule)

			continue
		}

		// check if the repo has a valid owner
		if repo.GetUserID() == 0 {
			logrus.Warnf("%s for %s: repo does not have a valid owner", baseErr, entry)

			continue
		}

		// send API call to capture the owner for the repo
		user, err := database.GetUser(repo.GetUserID())
		if err != nil {
			logrus.WithError(err).Warnf("%s for %s: unable to get owner", baseErr, entry)

			continue
		}

		// send API call to confirm repo owner has at least write access to repo
		_, err = scm.RepoAccess(user, user.GetToken(), repo.GetOrg(), repo.GetName())
		if err != nil {
			logrus.WithError(err).Warnf("%s for %s: %s does not have at least write access", baseErr, entry, user.GetName())

			continue
		}

		// create SQL filters for querying pending and running builds for repo
		filters := map[string]interface{}{
			"status": []string{constants.StatusPending, constants.StatusRunning},
		}

		// send API call to capture the number of pending or running builds for the repo
		builds, err := database.GetRepoBuildCount(repo, filters)
		if err != nil {
			logrus.WithError(err).Warnf("%s for %s: unable to get count of builds", baseErr, entry)

			continue
		}

		// check if the number of pending and running builds exceeds the limit for the repo
		if builds >= repo.GetBuildLimit() {
			logrus.Warnf("%s for %s: repo has excceded the concurrent build limit of %d", baseErr, entry, repo.GetBuildLimit())

			continue
		}

		var (
			// variable to store the raw pipeline configuration
			config []byte
			// variable to store executable pipeline
			p *pipeline.Build
			// variable to store pipeline configuration
			pipeline *library.Pipeline
			// variable to control number of times to retry processing pipeline
			retryLimit = 3
			// variable to store the pipeline type for the repository
			pipelineType = repo.GetPipelineType()
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
			pipeline, err = database.GetPipelineForRepo(repo.GetBranch(), repo)
			if err != nil { // assume the pipeline doesn't exist in the database yet
				// send API call to capture the pipeline configuration file
				config, err = scm.ConfigBackoff(user, repo, repo.GetBranch())
				if err != nil {
					logrus.WithError(err).Warnf("%s for %s: unable to get pipeline configuration", baseErr, entry)

					continue
				}
			} else {
				config = pipeline.GetData()
			}

			// send API call to capture repo for the counter (grabbing repo again to ensure counter is correct)
			repo, err = database.GetRepoForOrg(repo.GetOrg(), repo.GetName())
			if err != nil {
				logrus.WithError(err).Warnf("%s for %s: unable to get repo", baseErr, entry)

				continue
			}

			// set the parent equal to the current repo counter
			b.SetParent(repo.GetCounter())
			// check if the parent is set to 0
			if b.GetParent() == 0 {
				// parent should be "1" if it's the first build ran
				b.SetParent(1)
			}
			// update the build numbers based off repo counter
			inc := repo.GetCounter() + 1
			repo.SetCounter(inc)
			b.SetNumber(inc)

			// populate the build link if a web address is provided
			if len(metadata.Vela.WebAddress) > 0 {
				b.SetLink(
					fmt.Sprintf("%s/%s/%d", metadata.Vela.WebAddress, repo.GetFullName(), b.GetNumber()),
				)
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
				logrus.WithError(err).Warnf("%s for %s: unable to compile pipeline configuration", baseErr, entry)

				continue
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
				// set build to successful status
				b.SetStatus(constants.StatusSkipped)

				continue
			}

			// check if the pipeline did not already exist in the database
			if pipeline == nil {
				pipeline = compiled
				pipeline.SetRepoID(repo.GetID())
				pipeline.SetCommit(b.GetCommit())
				pipeline.SetRef(b.GetRef())

				// send API call to create the pipeline
				err = database.FromContext(c).CreatePipeline(pipeline)
				if err != nil {
					retErr := fmt.Errorf("%s: failed to create pipeline for %s: %w", baseErr, repo.GetFullName(), err)

					// check if the retry limit has been exceeded
					if i < retryLimit-1 {
						logrus.WithError(retErr).Warningf("retrying #%d", i+1)

						// continue to the next iteration of the loop
						continue
					}

					util.HandleError(c, http.StatusBadRequest, retErr)

					h.SetStatus(constants.StatusFailure)
					h.SetError(retErr.Error())

					return
				}

				// send API call to capture the created pipeline
				pipeline, err = database.GetPipelineForRepo(pipeline.GetCommit(), repo)
				if err != nil {
					//nolint:lll // ignore long line length due to error message
					retErr := fmt.Errorf("%s: failed to get new pipeline %s/%s: %w", baseErr, repo.GetFullName(), pipeline.GetCommit(), err)
					util.HandleError(c, http.StatusInternalServerError, retErr)

					h.SetStatus(constants.StatusFailure)
					h.SetError(retErr.Error())

					return
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
			err = planBuild(database, p, b, repo)
			if err != nil {
				retErr := fmt.Errorf("%s: %w", baseErr, err)

				// check if the retry limit has been exceeded
				if i < retryLimit-1 {
					logrus.WithError(retErr).Warningf("retrying #%d", i+1)

					// reset fields set by cleanBuild for retry
					b.SetError("")
					b.SetStatus(constants.StatusPending)
					b.SetFinished(0)

					// continue to the next iteration of the loop
					continue
				}

				util.HandleError(c, http.StatusInternalServerError, retErr)

				h.SetStatus(constants.StatusFailure)
				h.SetError(retErr.Error())

				return
			}

			// break the loop because everything was successful
			break
		} // end of retry loop

	}

	return nil
}
