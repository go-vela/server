// SPDX-License-Identifier: Apache-2.0

package service

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// ListServicesForBuild gets a list of all services from the database.
func (e *engine) ListServicesForBuild(ctx context.Context, b *api.Build, filters map[string]interface{}, page int, perPage int) ([]*api.Service, int64, error) {
	e.logger.WithFields(logrus.Fields{
		"build": b.GetNumber(),
	}).Tracef("listing services for build %d", b.GetNumber())

	// variables to store query results and return value
	count := int64(0)
	s := new([]types.Service)
	services := []*api.Service{}

	// count the results
	count, err := e.CountServicesForBuild(ctx, b, filters)
	if err != nil {
		return services, 0, err
	}

	// short-circuit if there are no results
	if count == 0 {
		return services, 0, nil
	}

	// calculate offset for pagination through results
	offset := perPage * (page - 1)

	// send query to the database and store result in variable
	err = e.client.
		WithContext(ctx).
		Table(constants.TableService).
		Where("build_id = ?", b.GetID()).
		Where(filters).
		Order("id DESC").
		Limit(perPage).
		Offset(offset).
		Find(&s).
		Error
	if err != nil {
		return nil, count, err
	}

	// iterate through all query results
	for _, service := range *s {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := service

		services = append(services, tmp.ToAPI())
	}

	return services, count, nil
}
