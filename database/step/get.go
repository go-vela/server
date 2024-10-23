// SPDX-License-Identifier: Apache-2.0

package step

import (
	"context"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// GetStep gets a step by ID from the database.
func (e *engine) GetStep(ctx context.Context, id int64) (*api.Step, error) {
	e.logger.Tracef("getting step %d", id)

	// variable to store query results
	s := new(types.Step)

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TableStep).
		Where("id = ?", id).
		Take(s).
		Error
	if err != nil {
		return nil, err
	}

	return s.ToAPI(), nil
}
