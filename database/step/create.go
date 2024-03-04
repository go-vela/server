// SPDX-License-Identifier: Apache-2.0

package step

import (
	"context"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// CreateStep creates a new step in the database.
func (e *engine) CreateStep(ctx context.Context, s *library.Step) (*library.Step, error) {
	e.logger.WithFields(logrus.Fields{
		"step": s.GetNumber(),
	}).Tracef("creating step %s in the database", s.GetName())

	// cast the library type to database type
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#StepFromLibrary
	step := database.StepFromLibrary(s)

	// validate the necessary fields are populated
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Step.Validate
	err := step.Validate()
	if err != nil {
		return nil, err
	}

	// send query to the database
	result := e.client.Table(constants.TableStep).Create(step)

	return step.ToLibrary(), result.Error
}
