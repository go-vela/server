// SPDX-License-Identifier: Apache-2.0

package service

import (
	"context"

	"github.com/go-vela/types/constants"
)

// CountServices gets the count of all services from the database.
func (e *engine) CountServices(ctx context.Context) (int64, error) {
	e.logger.Tracef("getting count of all services from the database")

	// variable to store query results
	var s int64

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableService).
		Count(&s).
		Error

	return s, err
}
