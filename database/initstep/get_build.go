// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package initstep

import (
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// GetInitStepForBuild gets an init step by build ID and number from the database.
func (e *engine) GetInitStepForBuild(b *library.Build, number int) (*library.InitStep, error) {
	e.logger.WithFields(logrus.Fields{
		"build":    b.GetNumber(),
		"initstep": number,
	}).Tracef("getting init step #%d for build %d from the database", number, b.GetID())

	// variable to store query results
	h := new(database.InitStep)

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableInitStep).
		Where("build_id = ?", b.GetID()).
		Where("number = ?", number).
		Take(h).
		Error
	if err != nil {
		return nil, err
	}

	// return the InitStep
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#InitStep.ToLibrary
	return h.ToLibrary(), nil
}
