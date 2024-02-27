// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/adhocore/gronx"
	"github.com/go-vela/server/api/build"
	"github.com/go-vela/server/compiler"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/queue"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"
	"github.com/sirupsen/logrus"

	"k8s.io/apimachinery/pkg/util/wait"
)

const (
	scheduleErr = "unable to trigger build for schedule"

	scheduleWait = "waiting to trigger build for schedule"
)

func processSchedules(ctx context.Context, start time.Time, compiler compiler.Engine, database database.Interface, metadata *types.Metadata, queue queue.Service, scm scm.Service, allowList []string) error {
	logrus.Infof("processing active schedules to create builds")

	// send API call to capture the list of active schedules
	schedules, err := database.ListActiveSchedules(ctx)
	if err != nil {
		return err
	}

	// iterate through the list of active schedules
	for _, s := range schedules {
		// sleep for 1s - 2s before processing the active schedule
		//
		// This should prevent multiple servers from processing a schedule at the same time by
		// leveraging a base duration along with a standard deviation of randomness a.k.a.
		// "jitter". To create the jitter, we use a base duration of 1s with a scale factor of 1.0.
		time.Sleep(wait.Jitter(time.Second, 1.0))

		// send API call to capture the schedule
		//
		// This is needed to ensure we are not dealing with a stale schedule since we fetch
		// all schedules once and iterate through that list which can take a significant
		// amount of time to get to the end of the list.
		schedule, err := database.GetSchedule(ctx, s.GetID())
		if err != nil {
			logrus.WithError(err).Warnf("%s %s", scheduleErr, schedule.GetName())

			continue
		}

		// ignore triggering a build if the schedule is no longer active
		if !schedule.GetActive() {
			logrus.Tracef("skipping to trigger build for inactive schedule %s", schedule.GetName())

			continue
		}

		// capture the last time a build was triggered for the schedule in UTC
		scheduled := time.Unix(schedule.GetScheduledAt(), 0).UTC()

		// capture the previous occurrence of the entry rounded to the nearest whole interval
		//
		// i.e. if it's 4:02 on five minute intervals, this will be 4:00
		prevTime, err := gronx.PrevTick(schedule.GetEntry(), true)
		if err != nil {
			logrus.WithError(err).Warnf("%s %s", scheduleErr, schedule.GetName())

			continue
		}

		// capture the next occurrence of the entry after the last schedule rounded to the nearest whole interval
		//
		// i.e. if it's 4:02 on five minute intervals, this will be 4:05
		nextTime, err := gronx.NextTickAfter(schedule.GetEntry(), scheduled, true)
		if err != nil {
			logrus.WithError(err).Warnf("%s %s", scheduleErr, schedule.GetName())

			continue
		}

		// check if we should wait to trigger a build for the schedule
		//
		// The current time must be after the next occurrence of the schedule.
		if !time.Now().After(nextTime) {
			logrus.Tracef("%s %s: current time not past next occurrence", scheduleWait, schedule.GetName())

			continue
		}

		// check if we should wait to trigger a build for the schedule
		//
		// The previous occurrence of the schedule must be after the starting time of processing schedules.
		if !prevTime.After(start) {
			logrus.Tracef("%s %s: previous occurrence not after starting point", scheduleWait, schedule.GetName())

			continue
		}

		// update the scheduled_at field with the current timestamp
		//
		// This should help prevent multiple servers from processing a schedule at the same time
		// by updating the schedule with a new timestamp to reflect the current state.
		schedule.SetScheduledAt(time.Now().UTC().Unix())

		// send API call to update schedule for ensuring scheduled_at field is set
		_, err = database.UpdateSchedule(ctx, schedule, false)
		if err != nil {
			logrus.WithError(err).Warnf("%s %s", scheduleErr, schedule.GetName())

			continue
		}

		// process the schedule and trigger a new build
		err = processSchedule(ctx, schedule, compiler, database, metadata, queue, scm, allowList)
		if err != nil {
			logrus.WithError(err).Warnf("%s %s", scheduleErr, schedule.GetName())

			continue
		}
	}

	return nil
}

//nolint:funlen // ignore function length and number of statements
func processSchedule(ctx context.Context, s *library.Schedule, compiler compiler.Engine, database database.Interface, metadata *types.Metadata, queue queue.Service, scm scm.Service, allowList []string) error {
	// send API call to capture the repo for the schedule
	r, err := database.GetRepo(ctx, s.GetRepoID())
	if err != nil {
		return fmt.Errorf("unable to fetch repo: %w", err)
	}

	// ensure repo has not been removed from allow list
	if !util.CheckAllowlist(r, allowList) {
		return fmt.Errorf("skipping schedule: repo %s no longer on allow list", r.GetFullName())
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
	u, err := database.GetUser(ctx, r.GetUserID())
	if err != nil {
		return fmt.Errorf("unable to get owner for repo %s: %w", r.GetFullName(), err)
	}

	// send API call to confirm repo owner has at least write access to repo
	_, err = scm.RepoAccess(ctx, u.GetName(), u.GetToken(), r.GetOrg(), r.GetName())
	if err != nil {
		return fmt.Errorf("%s does not have at least write access for repo %s", u.GetName(), r.GetFullName())
	}

	// create SQL filters for querying pending and running builds for repo
	filters := map[string]interface{}{
		"status": []string{constants.StatusPending, constants.StatusRunning},
	}

	// send API call to capture the number of pending or running builds for the repo
	builds, err := database.CountBuildsForRepo(ctx, r, filters)
	if err != nil {
		return fmt.Errorf("unable to get count of builds for repo %s: %w", r.GetFullName(), err)
	}

	// check if the number of pending and running builds exceeds the limit for the repo
	if builds >= r.GetBuildLimit() {
		return fmt.Errorf("repo %s has excceded the concurrent build limit of %d", r.GetFullName(), r.GetBuildLimit())
	}

	// send API call to capture the commit sha for the branch
	_, commit, err := scm.GetBranch(ctx, u, r, s.GetBranch())
	if err != nil {
		return fmt.Errorf("failed to get commit for repo %s on %s branch: %w", r.GetFullName(), r.GetBranch(), err)
	}

	url := strings.TrimSuffix(r.GetClone(), ".git")

	b := new(library.Build)
	b.SetAuthor(s.GetCreatedBy())
	b.SetBranch(s.GetBranch())
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
		pipeline, err = database.GetPipelineForRepo(ctx, b.GetCommit(), r)
		if err != nil { // assume the pipeline doesn't exist in the database yet
			// send API call to capture the pipeline configuration file
			config, err = scm.ConfigBackoff(ctx, u, r, b.GetCommit())
			if err != nil {
				return fmt.Errorf("unable to get pipeline config for %s/%s: %w", r.GetFullName(), b.GetCommit(), err)
			}
		} else {
			config = pipeline.GetData()
		}

		// send API call to capture repo for the counter (grabbing repo again to ensure counter is correct)
		r, err = database.GetRepoForOrg(ctx, r.GetOrg(), r.GetName())
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
			pipeline, err = database.CreatePipeline(ctx, pipeline)
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
		}

		b.SetPipelineID(pipeline.GetID())

		// create the objects from the pipeline in the database
		// TODO:
		// - if a build gets created and something else fails midway,
		//   the next loop will attempt to create the same build,
		//   using the same Number and thus create a constraint
		//   conflict; consider deleting the partially created
		//   build object in the database
		err = build.PlanBuild(ctx, database, p, b, r)
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

	// send API call to update repo for ensuring counter is incremented
	r, err = database.UpdateRepo(ctx, r)
	if err != nil {
		return fmt.Errorf("unable to update repo %s: %w", r.GetFullName(), err)
	}

	// send API call to capture the triggered build
	b, err = database.GetBuildForRepo(ctx, r, b.GetNumber())
	if err != nil {
		return fmt.Errorf("unable to get new build %s/%d: %w", r.GetFullName(), b.GetNumber(), err)
	}

	// determine queue route
	route, err := queue.Route(&p.Worker)
	if err != nil {
		logrus.Errorf("unable to set route for build %d for %s: %v", b.GetNumber(), r.GetFullName(), err)

		// error out the build
		build.CleanBuild(ctx, database, b, nil, nil, err)

		return err
	}

	// temporarily set host to the route before it gets picked up by a worker
	b.SetHost(route)

	err = build.PublishBuildExecutable(ctx, database, p, b)
	if err != nil {
		retErr := fmt.Errorf("unable to publish build executable for %s/%d: %w", r.GetFullName(), b.GetNumber(), err)

		return retErr
	}

	// publish the build to the queue
	go build.PublishToQueue(
		ctx,
		queue,
		database,
		b,
		r,
		u,
		route,
	)

	return nil
}
