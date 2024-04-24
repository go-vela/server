// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/types"
	"github.com/go-vela/types/constants"
)

// CleanBuilds updates builds to an error with a provided message with a created timestamp prior to a defined moment.
func (e *engine) CleanBuilds(ctx context.Context, msg string, before int64) (int64, error) {
	logrus.Tracef("cleaning pending or running builds in the database created prior to %d", before)

	b := new(api.Build)
	b.SetStatus(constants.StatusError)
	b.SetError(msg)
	b.SetFinished(time.Now().UTC().Unix())

	build := types.BuildFromAPI(b)

	// send query to the database
	result := e.client.
		Table(constants.TableBuild).
		Where("created < ?", before).
		Where("status = 'running' OR status = 'pending'").
		Updates(build)

	return result.RowsAffected, result.Error
}
