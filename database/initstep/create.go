// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package initstep

import (
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// CreateInitStep creates a new init step in the database.
func (e *engine) CreateInitStep(i *library.InitStep) error {
	e.logger.WithFields(logrus.Fields{
		"repo_id":  i.GetRepoID(),
		"build_id": i.GetBuildID(),
		"initstep": i.GetNumber(),
	}).Tracef("creating init step %d in the database", i.GetNumber())

	// cast the library type to database type
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#InitStepFromLibrary
	initStep := database.InitStepFromLibrary(i)

	// validate the necessary fields are populated
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#InitStep.Validate
	err := initStep.Validate()
	if err != nil {
		return err
	}

	// send query to the database
	return e.client.
		Table(constants.TableInitStep).
		Create(initStep).
		Error
}
