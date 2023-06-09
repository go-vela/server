// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/adhocore/gronx"
	"github.com/go-vela/server/api/build"
	"github.com/go-vela/server/compiler"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/queue"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/types"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"
	"github.com/sirupsen/logrus"

	"k8s.io/apimachinery/pkg/util/wait"
)

const baseErr = "unable to schedule build"

func processSchedules(compiler compiler.Engine, database database.Interface, metadata *types.Metadata, queue queue.Service, scm scm.Service) error {
	logrus.Infof("processing active schedules to create builds")

	// send API call to capture the list of active schedules
	schedules, err := database.ListActiveSchedules()
	if err != nil {
		return err
	}

	// iterate through the list of active schedules
	for _, s := range schedules {
		// send API call to capture the schedule
		//
		// This is needed to ensure we are not dealing with a stale schedule since we fetch
		// all schedules once and iterate through that list which can take a significant
		// amount of time to get to the end of the list.
		schedule, err := database.GetSchedule(s.GetID())
		if err != nil {
			logrus.WithError(err).Warnf("%s for %s", baseErr, schedule.GetName())

			continue
		}

		// create a variable to track if a build should be triggered based off the schedule
		trigger := false

		// check if a build has already been triggered for the schedule
		if schedule.GetScheduledAt() == 0 {
			// trigger a build for the schedule since one has not already been scheduled
			trigger = true
		} else {
			// parse the previous occurrence of the entry for the schedule
			prevTime, err := gronx.PrevTick(schedule.GetEntry(), true)
			if err != nil {
				logrus.WithError(err).Warnf("%s for %s", baseErr, schedule.GetName())

				continue
			}

			// parse the next occurrence of the entry for the schedule
			nextTime, err := gronx.NextTick(schedule.GetEntry(), true)
			if err != nil {
				logrus.WithError(err).Warnf("%s for %s", baseErr, schedule.GetName())

				continue
			}

			// parse the UNIX timestamp from when the last build was triggered for the schedule
			t := time.Unix(schedule.GetScheduledAt(), 0).UTC()

			// check if the time since the last triggered build is greater than the entry duration for the schedule
			if time.Since(t) > nextTime.Sub(prevTime) {
				// trigger a build for the schedule since it has not previously ran
				trigger = true
			}
		}

		if trigger && schedule.GetActive() {
			err = processSchedule(schedule, compiler, database, metadata, queue, scm)
			if err != nil {
				logrus.WithError(err).Warnf("%s for %s", baseErr, schedule.GetName())

				continue
			}
		}
	}

	return nil
}

//nolint:funlen // ignore function length and number of statements
func processSchedule(s *library.Schedule, compiler compiler.Engine, database database.Interface, metadata *types.Metadata, queue queue.Service, scm scm.Service) error {
	// sleep for 1s - 3s before processing the schedule
	//
	// This should prevent multiple servers from processing a schedule at the same time by
	// leveraging a base duration along with a standard deviation of randomness a.k.a.
	// "jitter". To create the jitter, we use a base duration of 1s with a scale factor of 3.0.
	time.Sleep(wait.Jitter(time.Second, 3.0))

	// send API call to capture the repo for the schedule
	r, err := database.GetRepo(s.GetRepoID())
	if err != nil {
		return fmt.Errorf("unable to fetch repo: %w", err)
	}

	logrus.Tracef("processing schedule %s/%s", r.GetFullName(), s.GetName())

	// check if the repo is active
	if !r.GetActive() {
		return fmt.Errorf("repo %s is not active", r.GetFullName())
	}

	// check if the repo has a valid owner
	if r.GetUserID() == 0 {
		return fmt.Errorf("repo %s does not have a valid owner", r.GetFullName())
	}

	// send API call to capture the owner for the repo
	u, err := database.GetUser(r.GetUserID())
	if err != nil {
		return fmt.Errorf("unable to get owner for repo %s: %w", r.GetFullName(), err)
	}

	// send API call to confirm repo owner has at least write access to repo
	_, err = scm.RepoAccess(u, u.GetToken(), r.GetOrg(), r.GetName())
	if err != nil {
		return fmt.Errorf("%s does not have at least write access for repo %s", u.GetName(), r.GetFullName())
	}

	// create SQL filters for querying pending and running builds for repo
	filters := map[string]interface{}{
		"status": []string{constants.StatusPending, constants.StatusRunning},
	}

	// send API call to capture the number of pending or running builds for the repo
	builds, err := database.CountBuildsForRepo(r, filters)
	if err != nil {
		return fmt.Errorf("unable to get count of builds for repo %s: %w", r.GetFullName(), err)
	}

	// check if the number of pending and running builds exceeds the limit for the repo
	if builds >= r.GetBuildLimit() {
		return fmt.Errorf("repo %s has excceded the concurrent build limit of %d", r.GetFullName(), r.GetBuildLimit())
	}

	// send API call to capture the commit sha for the branch
	_, commit, err := scm.GetBranch(u, r)
	if err != nil {
		return fmt.Errorf("failed to get commit for repo %s on %s branch: %w", r.GetFullName(), r.GetBranch(), err)
	}

	url := strings.TrimSuffix(r.GetClone(), ".git")

	b := new(library.Build)
	b.SetAuthor(s.GetCreatedBy())
	b.SetBranch(r.GetBranch())
	b.SetClone(r.GetClone())
	b.SetCommit(commit)
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
			config, err = scm.ConfigBackoff(u, r, b.GetCommit())
			if err != nil {
				return fmt.Errorf("unable to get pipeline config for %s/%s: %w", r.GetFullName(), b.GetCommit(), err)
			}
		} else {
			config = pipeline.GetData()
		}

		// send API call to capture repo for the counter (grabbing repo again to ensure counter is correct)
		r, err = database.GetRepoForOrg(r.GetOrg(), r.GetName())
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

		// set the build numbers based off repo counter
		b.SetNumber(r.GetCounter() + 1)
		// set the parent equal to the current repo counter
		b.SetParent(r.GetCounter())
		// check if the parent is set to 0
		if b.GetParent() == 0 {
			// parent should be "1" if it's the first build ran
			b.SetParent(1)
		}
		r.SetCounter(r.GetCounter() + 1)

		// set the build link if a web address is provided
		if len(metadata.Vela.WebAddress) > 0 {
			b.SetLink(fmt.Sprintf("%s/%s/%d", metadata.Vela.WebAddress, r.GetFullName(), b.GetNumber()))
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
		p, compiled, err = compiler.
			Duplicate().
			WithBuild(b).
			WithCommit(b.GetCommit()).
			WithMetadata(metadata).
			WithRepo(r).
			WithUser(u).
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
		r.SetPipelineType(pipelineType)

		// skip the build if only the init or clone steps are found
		skip := build.SkipEmptyBuild(p)
		if skip != "" {
			return nil
		}

		// check if the pipeline did not already exist in the database
		if pipeline == nil {
			pipeline = compiled
			pipeline.SetRepoID(r.GetID())
			pipeline.SetCommit(b.GetCommit())
			pipeline.SetRef(b.GetRef())

			// send API call to create the pipeline
			err = database.CreatePipeline(pipeline)
			if err != nil {
				err = fmt.Errorf("failed to create pipeline for %s: %w", r.GetFullName(), err)

				// check if the retry limit has been exceeded
				if i < retryLimit-1 {
					logrus.WithError(err).Warningf("retrying #%d", i+1)

					// continue to the next iteration of the loop
					continue
				}

				return err
			}

			// send API call to capture the created pipeline
			pipeline, err = database.GetPipelineForRepo(pipeline.GetCommit(), r)
			if err != nil {
				return fmt.Errorf("unable to get new pipeline %s/%s: %w", r.GetFullName(), pipeline.GetCommit(), err)
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
		err = build.PlanBuild(database, p, b, r)
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

		s.SetScheduledAt(time.Now().UTC().Unix())

		// break the loop because everything was successful
		break
	} // end of retry loop

	// send API call to update repo for ensuring counter is incremented
	err = database.UpdateRepo(r)
	if err != nil {
		return fmt.Errorf("unable to update repo %s: %w", r.GetFullName(), err)
	}

	// send API call to update schedule for ensuring scheduled_at field is set
	err = database.UpdateSchedule(s)
	if err != nil {
		return fmt.Errorf("unable to update schedule %s/%s: %w", r.GetFullName(), s.GetName(), err)
	}

	// send API call to capture the triggered build
	b, err = database.GetBuildForRepo(r, b.GetNumber())
	if err != nil {
		return fmt.Errorf("unable to get new build %s/%d: %w", r.GetFullName(), b.GetNumber(), err)
	}

	// publish the build to the queue
	go build.PublishToQueue(
		queue,
		database,
		p,
		b,
		r,
		u,
	)

	return nil
}
