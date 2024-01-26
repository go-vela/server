// SPDX-License-Identifier: Apache-2.0

package step

import (
	"context"
	"time"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// CleanSteps updates steps to an error with a created timestamp prior to a defined moment.
func (e *engine) CleanSteps(ctx context.Context, msg string, before int64) (int64, error) {
	logrus.Tracef("cleaning pending or running steps in the database created prior to %d", before)

	s := new(library.Step)
	s.SetStatus(constants.StatusError)
	s.SetError(msg)
	s.SetFinished(time.Now().UTC().Unix())

	step := database.StepFromLibrary(s)

	// send query to the database
	result := e.client.
		Table(constants.TableStep).
		Where("created < ?", before).
		Where("status = 'running' OR status = 'pending'").
		Updates(step)

	return result.RowsAffected, result.Error
}
