// SPDX-License-Identifier: Apache-2.0

package service

import (
	"context"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// GetService gets a service by ID from the database.
func (e *engine) GetService(ctx context.Context, id int64) (*api.Service, error) {
	e.logger.Tracef("getting service %d", id)

	// variable to store query results
	s := new(types.Service)

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TableService).
		Where("id = ?", id).
		Take(s).
		Error
	if err != nil {
		return nil, err
	}

	return s.ToAPI(), nil
}
