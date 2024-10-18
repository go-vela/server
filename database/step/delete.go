// SPDX-License-Identifier: Apache-2.0

package step

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// DeleteStep deletes an existing step from the database.
func (e *engine) DeleteStep(ctx context.Context, s *api.Step) error {
	e.logger.WithFields(logrus.Fields{
		"step": s.GetNumber(),
	}).Tracef("deleting step %s", s.GetName())

	// cast the library type to database type
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#StepFromLibrary
	step := types.StepFromAPI(s)

	// send query to the database
	return e.client.
		WithContext(ctx).
		Table(constants.TableStep).
		Delete(step).
		Error
}
