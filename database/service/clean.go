// SPDX-License-Identifier: Apache-2.0

package service

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// CleanServices updates services to an error with a created timestamp prior to a defined moment.
func (e *engine) CleanServices(ctx context.Context, msg string, before int64) (int64, error) {
	logrus.Tracef("cleaning pending or running steps in the database created prior to %d", before)

	s := new(api.Service)
	s.SetStatus(constants.StatusError)
	s.SetError(msg)
	s.SetFinished(time.Now().UTC().Unix())

	service := types.ServiceFromAPI(s)

	// send query to the database
	result := e.client.
		WithContext(ctx).
		Table(constants.TableService).
		Where("created < ?", before).
		Where("status = 'running' OR status = 'pending'").
		Updates(service)

	return result.RowsAffected, result.Error
}
