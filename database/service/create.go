// SPDX-License-Identifier: Apache-2.0

package service

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// CreateService creates a new service in the database.
func (e *Engine) CreateService(ctx context.Context, s *api.Service) (*api.Service, error) {
	e.logger.WithFields(logrus.Fields{
		"service": s.GetNumber(),
	}).Tracef("creating service %s in the database", s.GetName())

	service := types.ServiceFromAPI(s)

	err := service.Validate()
	if err != nil {
		return nil, err
	}

	// send query to the database
	result := e.client.
		WithContext(ctx).
		Table(constants.TableService).
		Create(service)

	return service.ToAPI(), result.Error
}
