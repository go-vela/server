// SPDX-License-Identifier: Apache-2.0

package build

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/scm"
)

// GatekeepBuild is a helper function that will set the status of a build to 'pending approval' and
// send a status update to the SCM.
func GatekeepBuild(c *gin.Context, b *types.Build, r *types.Repo) error {
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
