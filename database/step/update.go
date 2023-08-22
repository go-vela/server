// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package step

import (
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// UpdateStep updates an existing step in the database.
func (e *engine) UpdateStep(s *library.Step) (*library.Step, error) {
	e.logger.WithFields(logrus.Fields{
		"step": s.GetNumber(),
	}).Tracef("updating step %s in the database", s.GetName())

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
	result := e.client.Table(constants.TableStep).Save(step)

	return step.ToLibrary(), result.Error
}
