// SPDX-License-Identifier: Apache-2.0

package service

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// UpdateService updates an existing service in the database.
func (e *engine) UpdateService(ctx context.Context, s *api.Service) (*api.Service, error) {
	e.logger.WithFields(logrus.Fields{
		"service": s.GetNumber(),
	}).Tracef("updating service %s", s.GetName())

	service := types.ServiceFromAPI(s)

	err := service.Validate()
	if err != nil {
		return nil, err
	}

	// send query to the database
	result := e.client.
		WithContext(ctx).
		Table(constants.TableService).
		Save(service)

	return service.ToAPI(), result.Error
}
