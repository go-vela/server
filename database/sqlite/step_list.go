// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

// nolint: dupl // ignore false positive of duplicate code
package sqlite

import (
	"github.com/go-vela/server/database/sqlite/dml"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"

	"github.com/sirupsen/logrus"
)

// GetStepList gets a list of all steps from the database.
func (c *client) GetStepList() ([]*library.Step, error) {
	logrus.Trace("listing steps from the database")

	// variable to store query results
	s := new([]database.Step)

	// send query to the database and store result in variable
	err := c.Sqlite.
		Table(constants.TableStep).
		Raw(dml.ListSteps).
		Scan(s).Error

	// variable we want to return
	steps := []*library.Step{}
	// iterate through all query results
	for _, step := range *s {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := step

		// convert query result to library type
		steps = append(steps, tmp.ToLibrary())
	}

	return steps, err
}

// GetBuildStepList gets a list of steps by build ID from the database.
func (c *client) GetBuildStepList(b *library.Build, page, perPage int) ([]*library.Step, error) {
	logrus.Tracef("listing steps for build %d from the database", b.GetNumber())

	// variable to store query results
	s := new([]database.Step)
	// calculate offset for pagination through results
	offset := (perPage * (page - 1))

	// send query to the database and store result in variable
	err := c.Sqlite.
		Table(constants.TableStep).
		Raw(dml.ListBuildSteps, b.GetID(), perPage, offset).
		Scan(s).Error

	// variable we want to return
	steps := []*library.Step{}
	// iterate through all query results
	for _, step := range *s {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := step

		// convert query result to library type
		steps = append(steps, tmp.ToLibrary())
	}

	return steps, err
}
