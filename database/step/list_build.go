// SPDX-License-Identifier: Apache-2.0

package step

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
)

// ListStepsForBuild gets a list of all steps from the database.
func (e *engine) ListStepsForBuild(ctx context.Context, b *library.Build, filters map[string]interface{}, page int, perPage int) ([]*library.Step, int64, error) {
	e.logger.WithFields(logrus.Fields{
		"build": b.GetNumber(),
	}).Tracef("listing steps for build %d from the database", b.GetNumber())

	// variables to store query results and return value
	count := int64(0)
	s := new([]database.Step)
	steps := []*library.Step{}

	// count the results
	count, err := e.CountStepsForBuild(ctx, b, filters)
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
