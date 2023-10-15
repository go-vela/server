// SPDX-License-Identifier: Apache-2.0

package step

import (
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
)

// ListSteps gets a list of all steps from the database.
func (e *engine) ListSteps() ([]*library.Step, error) {
	e.logger.Trace("listing all steps from the database")

	// variables to store query results and return value
	count := int64(0)
	w := new([]database.Step)
	steps := []*library.Step{}

	// count the results
	count, err := e.CountSteps()
	if err != nil {
		return nil, err
	}

	// short-circuit if there are no results
	if count == 0 {
		return steps, nil
	}

	// send query to the database and store result in variable
	err = e.client.
		Table(constants.TableStep).
		Find(&w).
		Error
	if err != nil {
		return nil, err
	}

	// iterate through all query results
	for _, step := range *w {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := step

		// convert query result to library type
		//
		// https://pkg.go.dev/github.com/go-vela/types/database#Step.ToLibrary
		steps = append(steps, tmp.ToLibrary())
	}

	return steps, nil
}
