// SPDX-License-Identifier: Apache-2.0

package service

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// DeleteService deletes an existing service from the database.
func (e *engine) DeleteService(ctx context.Context, s *api.Service) error {
	e.logger.WithFields(logrus.Fields{
		"service": s.GetNumber(),
	}).Tracef("deleting service %s", s.GetName())

	service := types.ServiceFromAPI(s)

	// send query to the database
	return e.client.
		WithContext(ctx).
		Table(constants.TableService).
		Delete(service).
		Error
}
