// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/adhocore/gronx"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/go-vela/server/api/build"
	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/api/types/settings"
	"github.com/go-vela/server/compiler"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/internal"
	"github.com/go-vela/server/queue"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/util"
)

const (
	scheduleErr = "unable to trigger build for schedule"

	scheduleWait = "waiting to trigger build for schedule"
)

func processSchedules(ctx context.Context, start time.Time, settings *settings.Platform, compiler compiler.Engine, database database.Interface, metadata *internal.Metadata, queue queue.Service, scm scm.Service) error {
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
			handleError(ctx, database, err, schedule)

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
			handleError(ctx, database, err, schedule)

			continue
		}

		// capture the next occurrence of the entry after the last schedule rounded to the nearest whole interval
		//
		// i.e. if it's 4:02 on five minute intervals, this will be 4:05
		nextTime, err := gronx.NextTickAfter(schedule.GetEntry(), scheduled, true)
		if err != nil {
			handleError(ctx, database, err, schedule)

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
			handleError(ctx, database, err, schedule)

			continue
		}

		logrus.WithFields(logrus.Fields{
			"schedule":    schedule.GetName(),
			"schedule_id": schedule.GetID(),
		}).Info("schedule updated - scheduled at set")

		// process the schedule and trigger a new build
		err = processSchedule(ctx, schedule, settings, compiler, database, metadata, queue, scm)
		if err != nil {
			handleError(ctx, database, err, schedule)

			continue
		}

		// successfully scheduled build so clear error message, if not already cleared
		if schedule.GetError() != "" {
			schedule.SetError("")

			// send API call to update schedule with the error message field cleared
			_, err = database.UpdateSchedule(ctx, schedule, true)
			if err != nil {
				handleError(ctx, database, err, schedule)

				continue
			}

			logrus.WithFields(logrus.Fields{
				"schedule":    schedule.GetName(),
				"schedule_id": schedule.GetID(),
			}).Info("schedule updated - error message cleared")
		}
	}

	return nil
}

// processSchedule will, given a schedule, process it and trigger a new build.
func processSchedule(ctx context.Context, s *api.Schedule, settings *settings.Platform, compiler compiler.Engine, database database.Interface, metadata *internal.Metadata, queue queue.Service, scm scm.Service) error {
	// send API call to capture the repo for the schedule
	r, err := database.GetRepo(ctx, s.GetRepo().GetID())
	if err != nil {
		return fmt.Errorf("unable to fetch repo: %w", err)
	}

	// ensure repo has not been removed from allow list
	if !util.CheckAllowlist(r, settings.GetScheduleAllowlist()) {
		return fmt.Errorf("skipping schedule: repo %s no longer on allow list", r.GetFullName())
	}

	logrus.Tracef("processing schedule %s/%s", r.GetFullName(), s.GetName())

	// check if the repo is active
	if !r.GetActive() {
		return fmt.Errorf("repo %s is not active", r.GetFullName())
	}

	url := strings.TrimSuffix(r.GetClone(), ".git")

	b := new(api.Build)
	b.SetAuthor(s.GetCreatedBy())
	b.SetBranch(s.GetBranch())
	b.SetClone(r.GetClone())
	b.SetDeploy(s.GetName())
	b.SetEvent(constants.EventSchedule)
	b.SetMessage(fmt.Sprintf("triggered for %s schedule with %s entry", s.GetName(), s.GetEntry()))
	b.SetRef(fmt.Sprintf("refs/heads/%s", b.GetBranch()))
	b.SetRepo(r)
	b.SetSender(s.GetUpdatedBy())

	// fetch scm user id
	senderID, err := scm.GetUserID(ctx, s.GetUpdatedBy(), r.GetOwner().GetToken())
	if err != nil {
		return fmt.Errorf("unable to get SCM user id for %s: %w", s.GetUpdatedBy(), err)
	}

	b.SetSenderSCMID(senderID)

	b.SetSource(fmt.Sprintf("%s/tree/%s", url, b.GetBranch()))
	b.SetStatus(constants.StatusPending)
	b.SetTitle(fmt.Sprintf("%s received from %s", constants.EventSchedule, url))

	// schedule form
	config := build.CompileAndPublishConfig{
		Build:    b,
		Metadata: metadata,
		BaseErr:  "unable to schedule build",
		Source:   "schedule",
		Retries:  1,
	}

	_, item, _, err := build.CompileAndPublish(
		ctx,
		config,
		database,
		scm,
		compiler,
		queue,
	)

	if err != nil {
		return err
	}

	// publish the build to the queue
	go build.Enqueue(
		context.WithoutCancel(ctx),
		queue,
		database,
		item,
		item.Build.GetRoute(),
	)

	return nil
}

func handleError(ctx context.Context, database database.Interface, err error, schedule *api.Schedule) {
	// log the error message
	logrus.WithError(err).Warnf("%s %s: %s", scheduleErr, schedule.GetName(), err.Error())

	// format the error message
	msg := fmt.Sprintf("%s %s: %s", scheduleErr, schedule.GetName(), err.Error())

	// update the message field with the error message
	schedule.SetError(msg)

	// send API call to update schedule to ensure message field is set
	_, err = database.UpdateSchedule(ctx, schedule, true)
	if err != nil {
		logrus.WithError(err).Warnf("%s %s: %s", scheduleErr, schedule.GetName(), err.Error())
	}

	logrus.WithFields(logrus.Fields{
		"schedule":    schedule.GetName(),
		"schedule_id": schedule.GetID(),
	}).Info("schedule updated - error message set")
}
