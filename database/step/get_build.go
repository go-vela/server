// SPDX-License-Identifier: Apache-2.0

package step

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// GetStepForBuild gets a step by number and build ID from the database.
func (e *engine) GetStepForBuild(ctx context.Context, b *api.Build, number int) (*api.Step, error) {
	e.logger.WithFields(logrus.Fields{
		"build": b.GetNumber(),
		"step":  number,
	}).Tracef("getting step %d", number)

	// variable to store query results
	s := new(types.Step)

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TableStep).
		Where("build_id = ?", b.GetID()).
		Where("number = ?", number).
		Take(s).
		Error
	if err != nil {
		return nil, err
	}

	return s.ToAPI(), nil
}
