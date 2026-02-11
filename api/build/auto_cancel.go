// SPDX-License-Identifier: Apache-2.0

package build

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database"
)

// AutoCancel is a helper function that checks to see if any pending or running
// builds for the repo can be replaced by the current build.
func AutoCancel(c *gin.Context, b *types.Build, rB *types.Build, cancelOpts *pipeline.CancelOptions) (bool, error) {
	l := c.MustGet("logger").(*logrus.Entry)

	// in this path, the middleware doesn't inject build,
	// so we need to set it manually
	l = l.WithFields(logrus.Fields{
		"build":    b.GetNumber(),
		"build_id": b.GetID(),
	})

	l.Debug("checking if builds should be auto canceled")

	// if build is the current build, continue
	if rB.GetID() == b.GetID() {
		return false, nil
	}

	status := rB.GetStatus()

	// ensure criteria is met
	if isCancelable(rB, b) {
		switch {
		case strings.EqualFold(status, constants.StatusPendingApproval) ||
			(strings.EqualFold(status, constants.StatusPending) &&
				cancelOpts.Pending):
			// pending build will be handled gracefully by worker once pulled off queue
			rB.SetStatus(constants.StatusCanceled)

			_, err := database.FromContext(c).UpdateBuild(c, rB)
			if err != nil {
				return false, err
			}

			l.WithFields(logrus.Fields{
				"build":    rB.GetNumber(),
				"build_id": rB.GetID(),
			}).Info("build updated - build canceled")

			// remove executable from table
			_, err = database.FromContext(c).PopBuildExecutable(c, rB.GetID())
			if err != nil {
				return true, err
			}
		case strings.EqualFold(status, constants.StatusRunning) && cancelOpts.Running:
			// call cancelRunning routine for builds already running on worker
			_, err := CancelRunning(c, rB)
			if err != nil {
				return false, err
			}
		default:
			return false, nil
		}

		// set error message that references current build
		rB.SetError(fmt.Sprintf("%s build was auto canceled in favor of build %d", status, b.GetNumber()))

		_, err := database.FromContext(c).UpdateBuild(c, rB)
		if err != nil {
			// if this call fails, we still canceled the build, so return true
			return true, err
		}

		l.WithFields(logrus.Fields{
			"build":    rB.GetNumber(),
			"build_id": rB.GetID(),
		}).Info("build updated - build canceled")
	}

	return false, nil
}

// isCancelable is a helper function that determines whether a `target` build should be auto-canceled
// given a current build that intends to supersede it.
func isCancelable(target *types.Build, current *types.Build) bool {
	switch target.GetEvent() {
	case constants.EventPush:
		// target is cancelable if current build is also a push event and the branches are the same
		return strings.EqualFold(current.GetEvent(), constants.EventPush) && strings.EqualFold(current.GetBranch(), target.GetBranch())
	case constants.EventPull:
		cancelableAction := strings.EqualFold(target.GetEventAction(), constants.ActionOpened) ||
			strings.EqualFold(target.GetEventAction(), constants.ActionSynchronize)

		// target is cancelable if current build is also a pull event, target is an opened / synchronize action, and the current head ref matches target head ref
		return strings.EqualFold(current.GetEvent(), constants.EventPull) && cancelableAction && strings.EqualFold(current.GetHeadRef(), target.GetHeadRef())
	default:
		return false
	}
}

// ShouldAutoCancel is a helper function that determines whether or not a build should be eligible to
// auto cancel currently running / pending builds.
func ShouldAutoCancel(opts *pipeline.CancelOptions, b *types.Build, defaultBranch string) bool {
	// if the build is pending approval, it should always be eligible to auto cancel
	if strings.EqualFold(b.GetStatus(), constants.StatusPendingApproval) {
		return true
	}

	// if anything is provided in the auto_cancel metadata, then we start with true
	runAutoCancel := opts.Running || opts.Pending || opts.DefaultBranch

	switch b.GetEvent() {
	case constants.EventPush:
		// pushes to the default branch should only auto cancel if pipeline specifies default_branch: true
		if !opts.DefaultBranch && strings.EqualFold(b.GetBranch(), defaultBranch) {
			runAutoCancel = false
		}

		return runAutoCancel

	case constants.EventPull:
		// only synchronize actions of the pull_request event are eligible to auto cancel
		return runAutoCancel && (strings.EqualFold(b.GetEventAction(), constants.ActionSynchronize))
	default:
		return false
	}
}
