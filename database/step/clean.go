// SPDX-License-Identifier: Apache-2.0

package step

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// CleanSteps updates steps to an error with a created timestamp prior to a defined moment.
func (e *Engine) CleanSteps(ctx context.Context, msg string, before int64) (int64, error) {
	logrus.Tracef("cleaning pending or running steps in the database created prior to %d", before)

	s := new(api.Step)
	s.SetStatus(constants.StatusError)
	s.SetError(msg)
	s.SetFinished(time.Now().UTC().Unix())

	step := types.StepFromAPI(s)

	// send query to the database
	result := e.client.
		WithContext(ctx).
		Table(constants.TableStep).
		Where("created < ?", before).
		Where("status = 'running' OR status = 'pending'").
		Updates(step)

	return result.RowsAffected, result.Error
}
