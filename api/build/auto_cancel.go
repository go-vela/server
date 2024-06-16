// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/internal/token"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/pipeline"
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
			err := cancelRunning(c, rB)
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

	return true, nil
}

// cancelRunning is a helper function that determines the executor currently running a build and sends an API call
// to that executor's worker to cancel the build.
func cancelRunning(c *gin.Context, b *types.Build) error {
	l := c.MustGet("logger").(*logrus.Entry)

	e := new([]types.Executor)
	// retrieve the worker
	w, err := database.FromContext(c).GetWorkerForHostname(c, b.GetHost())
	if err != nil {
		return err
	}

	// prepare the request to the worker to retrieve executors
	client := http.DefaultClient
	client.Timeout = 30 * time.Second
	endpoint := fmt.Sprintf("%s/api/v1/executors", w.GetAddress())

	req, err := http.NewRequestWithContext(context.Background(), "GET", endpoint, nil)
	if err != nil {
		return err
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
		return err
	}

	// add the token to authenticate to the worker
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tkn))

	// make the request to the worker and check the response
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	// Read Response Body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// parse response and validate at least one item was returned
	err = json.Unmarshal(respBody, e)
	if err != nil {
		return err
	}

	for _, executor := range *e {
		// check each executor on the worker running the build to see if it's running the build we want to cancel
		if executor.Build.GetID() == b.GetID() {
			// prepare the request to the worker
			client := http.DefaultClient
			client.Timeout = 30 * time.Second

			// set the API endpoint path we send the request to
			u := fmt.Sprintf("%s/api/v1/executors/%d/build/cancel", w.GetAddress(), executor.GetID())

			req, err := http.NewRequestWithContext(context.Background(), "DELETE", u, nil)
			if err != nil {
				return err
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
				return err
			}

			// add the token to authenticate to the worker
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tkn))

			// perform the request to the worker
			resp, err := client.Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			l.Debugf("sent cancel request to worker %s (executor %d) for build %d", w.GetHostname(), executor.GetID(), b.GetID())

			// Read Response Body
			respBody, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			err = json.Unmarshal(respBody, b)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// isCancelable is a helper function that determines whether a `target` build should be auto-canceled
// given a current build that intends to supersede it.
func isCancelable(target *types.Build, current *types.Build) bool {
	switch target.GetEvent() {
	case constants.EventPush:
		// target is cancelable if current build is also a push event and the branches are the same
		return strings.EqualFold(current.GetEvent(), constants.EventPush) && strings.EqualFold(current.GetBranch(), target.GetBranch())
	case constants.EventPull:
		cancelableAction := strings.EqualFold(target.GetEventAction(), constants.ActionOpened) || strings.EqualFold(target.GetEventAction(), constants.ActionSynchronize)

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
