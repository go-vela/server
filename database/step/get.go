// SPDX-License-Identifier: Apache-2.0

package step

import (
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
)

// GetStep gets a step by ID from the database.
func (e *engine) GetStep(id int64) (*library.Step, error) {
	e.logger.Tracef("getting step %d from the database", id)

	// variable to store query results
	s := new(database.Step)

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableStep).
		Where("id = ?", id).
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
