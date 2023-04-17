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

// ListStepsForBuild gets a list of all steps from the database.
func (e *engine) ListStepsForBuild(b *library.Build, filters map[string]interface{}, page int, perPage int) ([]*library.Step, int64, error) {
	e.logger.WithFields(logrus.Fields{
		"build": b.GetNumber(),
	}).Tracef("listing steps for build %d from the database", b.GetNumber())

	// variables to store query results and return value
	count := int64(0)
	s := new([]database.Step)
	steps := []*library.Step{}

	// count the results
	count, err := e.CountStepsForBuild(b, filters)
	if err != nil {
		return steps, 0, err
	}

	// short-circuit if there are no results
	if count == 0 {
		return steps, 0, nil
	}

	// calculate offset for pagination through results
	offset := perPage * (page - 1)

	// send query to the database and store result in variable
	err = e.client.
		Table(constants.TableStep).
		Where("build_id = ?", b.GetID()).
		Where(filters).
		Order("id DESC").
		Limit(perPage).
		Offset(offset).
		Find(&s).
		Error
	if err != nil {
		return nil, count, err
	}

	// iterate through all query results
	for _, step := range *s {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := step

		// convert query result to library type
		//
		// https://pkg.go.dev/github.com/go-vela/types/database#Step.ToLibrary
		steps = append(steps, tmp.ToLibrary())
	}

	return steps, count, nil
}
