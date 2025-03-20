// SPDX-License-Identifier: Apache-2.0

package step

import (
	"context"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// ListSteps gets a list of all steps from the database.
func (e *engine) ListSteps(ctx context.Context) ([]*api.Step, error) {
	e.logger.Trace("listing all steps")

	// variables to store query results and return value
	w := new([]types.Step)
	steps := []*api.Step{}

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
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

		steps = append(steps, tmp.ToAPI())
	}

	return steps, nil
}
