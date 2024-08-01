// SPDX-License-Identifier: Apache-2.0

package step

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
)

// DeleteStep deletes an existing step from the database.
func (e *engine) DeleteStep(ctx context.Context, s *library.Step) error {
	e.logger.WithFields(logrus.Fields{
		"step": s.GetNumber(),
	}).Tracef("deleting step %s", s.GetName())

	// cast the library type to database type
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#StepFromLibrary
	step := database.StepFromLibrary(s)

	// send query to the database
	return e.client.
		Table(constants.TableStep).
		Delete(step).
		Error
}
