// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database"
)

// helper function to clean pending approval builds from the database.
func cleanupPendingApproval(ctx context.Context, db database.Interface) error {
	logrus.Debug("cleaning pending approval builds")

	before := time.Now().Add(-(time.Duration(24*constants.ApprovalTimeoutMin) * time.Hour)).Unix()

	builds, err := db.ListPendingApprovalBuilds(ctx, strconv.FormatInt(before, 10))
	if err != nil {
		return err
	}

	for _, build := range builds {
		threshold := time.Now().Add(-(time.Duration(24*build.GetRepo().GetApprovalTimeout()) * time.Hour)).Unix()

		if build.GetCreated() < threshold {
			_, err := db.PopBuildExecutable(ctx, build.GetID())
			if err != nil {
				return err
			}

			build.SetStatus(constants.StatusError)
			build.SetFinished(time.Now().Unix())
			build.SetError(fmt.Sprintf("build exceeded approval timeout of %d days", build.GetRepo().GetApprovalTimeout()))

			_, err = db.UpdateBuild(ctx, build)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
