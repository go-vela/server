// SPDX-License-Identifier: Apache-2.0

package service

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// GetServiceForBuild gets a service by number and build ID from the database.
func (e *Engine) GetServiceForBuild(ctx context.Context, b *api.Build, number int32) (*api.Service, error) {
	e.logger.WithFields(logrus.Fields{
		"build":   b.GetNumber(),
		"service": number,
	}).Tracef("getting service %d", number)

	// variable to store query results
	s := new(types.Service)

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TableService).
		Where("build_id = ?", b.GetID()).
		Where("number = ?", number).
		Take(s).
		Error
	if err != nil {
		return nil, err
	}

	return s.ToAPI(), nil
}
