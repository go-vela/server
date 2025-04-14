// SPDX-License-Identifier: Apache-2.0

package service

import (
	"context"

	"github.com/go-vela/server/constants"
)

// CountServices gets the count of all services from the database.
func (e *Engine) CountServices(ctx context.Context) (int64, error) {
	e.logger.Tracef("getting count of all services")

	// variable to store query results
	var s int64

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TableService).
		Count(&s).
		Error

	return s, err
}
