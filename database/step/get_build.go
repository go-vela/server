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

// GetStepForBuild gets a step by number and build ID from the database.
func (e *engine) GetStepForBuild(b *library.Build, number int) (*library.Step, error) {
	e.logger.WithFields(logrus.Fields{
		"build": b.GetNumber(),
		"step":  number,
	}).Tracef("getting step %d from the database", number)

	// variable to store query results
	s := new(database.Step)

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableStep).
		Where("build_id = ?", b.GetID()).
		Where("number = ?", number).
		Take(s).
		Error
	if err != nil {
		return nil, err
	}

	// return the step
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Step.ToLibrary
	return s.ToLibrary(), nil
}
