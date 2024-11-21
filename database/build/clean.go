// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/constants"
)

// CleanBuilds updates builds to an error with a provided message with a created timestamp prior to a defined moment.
func (e *engine) CleanBuilds(ctx context.Context, msg string, before int64) (int64, error) {
	logrus.Tracef("cleaning pending or running builds created prior to %d", before)

	// send query to the database
	result := e.client.
		WithContext(ctx).
		Table(constants.TableBuild).
		Where("created < ?", before).
		Where("status = 'running' OR status = 'pending'").
		Updates(map[string]interface{}{
			"status":   constants.StatusError,
			"error":    msg,
			"finished": time.Now().UTC().Unix(),
		})

	return result.RowsAffected, result.Error
}
