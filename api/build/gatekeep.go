// SPDX-License-Identifier: Apache-2.0

package build

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/cache"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/scm"
)

// GatekeepBuild is a helper function that will set the status of a build to 'pending approval' and
// send a status update to the SCM.
func GatekeepBuild(c *gin.Context, b *types.Build, r *types.Repo, token string) error {
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
	err = UpdateComponentStatuses(c, b, constants.StatusPendingApproval, token)
	if err != nil {
		return fmt.Errorf("unable to update build components for %s/%d: %w", r.GetFullName(), b.GetNumber(), err)
	}

	// send API call to set the status on the commit
	checks, err := scm.FromContext(c).Status(c, b, token, nil)
	if err != nil {
		l.Errorf("unable to set commit status for %s/%d: %v", r.GetFullName(), b.GetNumber(), err)
	}

	err = cache.FromContext(c).StoreCheckRuns(c, b.GetID(), checks, r)
	if err != nil {
		l.Errorf("unable to store check runs for %s/%d: %v", r.GetFullName(), b.GetNumber(), err)
	}

	return nil
}
