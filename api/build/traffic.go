// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/queue"
	"github.com/go-vela/server/queue/models"
	"github.com/go-vela/server/scm"
)

// TrafficBuild is a helper function that will determine whether to publish a build to the queue or place it
// in pending approval status.
func TrafficBuild(c *gin.Context, l *logrus.Entry, b *types.Build, r *types.Repo, item *models.Item) error {
	// if the webhook was from a Pull event from a forked repository, verify it is allowed to run
	if b.GetFork() {
		l.Tracef("inside %s workflow for fork PR build %s/%d", r.GetApproveBuild(), r.GetFullName(), b.GetNumber())

		switch r.GetApproveBuild() {
		case constants.ApproveForkAlways:
			err := gatekeepBuild(c, b, r)
			if err != nil {
				return err
			}

			return nil
		case constants.ApproveForkNoWrite:
			// determine if build sender has write access to parent repo. If not, this call will result in an error
			level, err := scm.FromContext(c).RepoAccess(c.Request.Context(), b.GetSender(), r.GetOwner().GetToken(), r.GetOrg(), r.GetName())
			if err != nil || (level != "admin" && level != "write") {
				err = gatekeepBuild(c, b, r)
				if err != nil {
					return err
				}

				return nil
			}

			l.Debugf("fork PR build %s/%d automatically running without approval. sender %s has %s access", r.GetFullName(), b.GetNumber(), b.GetSender(), level)
		case constants.ApproveOnce:
			// determine if build sender is in the contributors list for the repo
			//
			// NOTE: this call is cumbersome for repos with lots of contributors. Potential TODO: improve this if
			// GitHub adds a single-contributor API endpoint.
			contributor, err := scm.FromContext(c).RepoContributor(c.Request.Context(), r.GetOwner(), b.GetSender(), r.GetOrg(), r.GetName())
			if err != nil {
				return err
			}

			if !contributor {
				err = gatekeepBuild(c, b, r)
				if err != nil {
					return err
				}

				return nil
			}

			fallthrough
		case constants.ApproveNever:
			fallthrough
		default:
			l.Debugf("fork PR build %s/%d automatically running without approval", r.GetFullName(), b.GetNumber())
		}
	}

	// send API call to set the status on the commit
	err := scm.FromContext(c).Status(c.Request.Context(), r.GetOwner(), b, r.GetOrg(), r.GetName())
	if err != nil {
		l.Errorf("unable to set commit status for %s/%d: %v", r.GetFullName(), b.GetNumber(), err)
	}

	// publish the build to the queue
	go Enqueue(
		context.WithoutCancel(c.Request.Context()),
		queue.FromGinContext(c),
		database.FromContext(c),
		item,
		b.GetHost(),
	)

	return nil
}

// gatekeepBuild is a helper function that will set the status of a build to 'pending approval' and
// send a status update to the SCM.
func gatekeepBuild(c *gin.Context, b *types.Build, r *types.Repo) error {
	l := c.MustGet("logger").(*logrus.Entry)

	l = l.WithFields(logrus.Fields{
		"org":      r.GetOrg(),
		"repo":     r.GetName(),
		"repo_id":  r.GetID(),
		"build":    b.GetNumber(),
		"build_id": b.GetID(),
	})

	l.Debug("fork PR build waiting for approval")

	b.SetStatus(constants.StatusPendingApproval)

	_, err := database.FromContext(c).UpdateBuild(c, b)
	if err != nil {
		return fmt.Errorf("unable to update build for %s/%d: %w", r.GetFullName(), b.GetNumber(), err)
	}

	l.Info("build updated")

	// update the build components to pending approval status
	err = UpdateComponentStatuses(c, b, constants.StatusPendingApproval)
	if err != nil {
		return fmt.Errorf("unable to update build components for %s/%d: %w", r.GetFullName(), b.GetNumber(), err)
	}

	// send API call to set the status on the commit
	err = scm.FromContext(c).Status(c, r.GetOwner(), b, r.GetOrg(), r.GetName())
	if err != nil {
		l.Errorf("unable to set commit status for %s/%d: %v", r.GetFullName(), b.GetNumber(), err)
	}

	return nil
}
