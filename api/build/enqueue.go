// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/queue"
	"github.com/go-vela/server/queue/models"
	"github.com/go-vela/server/scm"
)

// Enqueue is a helper function that pushes a queue item (build, repo, user) to the queue.
func Enqueue(ctx context.Context, queue queue.Service, db database.Interface, item *models.Item, route string) {
	l := logrus.WithFields(logrus.Fields{
		"build":    item.Build.GetNumber(),
		"build_id": item.Build.GetID(),
		"org":      item.Build.GetRepo().GetOrg(),
		"repo":     item.Build.GetRepo().GetName(),
		"repo_id":  item.Build.GetRepo().GetID(),
	})

	l.Debug("adding item to queue")

	// push item on to the queue
	err := queue.Push(ctx, route, item.Build.GetID())
	if err != nil {
		l.Errorf("retrying; failed to publish build: %v", err)

		err = queue.Push(ctx, route, item.Build.GetID())
		if err != nil {
			l.Errorf("failed to publish build: %v", err)

			// error out the build
			CleanBuild(ctx, db, item.Build, nil, nil, err)

			return
		}
	}

	// update fields in build object
	item.Build.SetEnqueued(time.Now().UTC().Unix())

	// update the build in the db to reflect the time it was enqueued
	_, err = db.UpdateBuild(ctx, item.Build)
	if err != nil {
		l.Errorf("failed to update build during publish to queue: %v", err)
	}

	l.Info("updated build as enqueued")
}

// ShouldEnqueue is a helper function that will determine whether to publish a build to the queue or place it
// in pending approval status.
func ShouldEnqueue(c *gin.Context, l *logrus.Entry, b *types.Build, r *types.Repo) (bool, error) {
	// if the webhook was from a Pull event from a forked repository, verify it is allowed to run
	if b.GetFork() {
		l.Tracef("inside %s workflow for fork PR build %s/%d", r.GetApproveBuild(), r.GetFullName(), b.GetNumber())

		switch r.GetApproveBuild() {
		case constants.ApproveForkAlways:
			return false, nil
		case constants.ApproveForkNoWrite:
			// determine if build sender has write access to parent repo. If not, this call will result in an error
			level, err := scm.FromContext(c).RepoAccess(c.Request.Context(), b.GetSender(), r.GetOwner().GetToken(), r.GetOrg(), r.GetName())
			if err != nil || (level != constants.PermissionAdmin && level != constants.PermissionWrite) {
				//nolint:nilerr // an error here is not something we want to return since we are gating it anyway
				return false, nil
			}

			l.Debugf("fork PR build %s/%d automatically running without approval. sender %s has %s access", r.GetFullName(), b.GetNumber(), b.GetSender(), level)
		case constants.ApproveOnce:
			// determine if build sender is in the contributors list for the repo
			//
			// NOTE: this call is cumbersome for repos with lots of contributors. Potential TODO: improve this if
			// GitHub adds a single-contributor API endpoint.
			contributor, err := scm.FromContext(c).RepoContributor(c.Request.Context(), r.GetOwner(), b.GetSender(), r.GetOrg(), r.GetName())
			if err != nil {
				return false, err
			}

			if !contributor {
				return false, nil
			}

			fallthrough
		case constants.ApproveNever:
			fallthrough
		default:
			l.Debugf("fork PR build %s/%d automatically running without approval", r.GetFullName(), b.GetNumber())
		}
	}

	return true, nil
}
