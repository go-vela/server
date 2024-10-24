// SPDX-License-Identifier: Apache-2.0

package step

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// ListStepsForBuild gets a list of all steps from the database.
func (e *engine) ListStepsForBuild(ctx context.Context, b *api.Build, filters map[string]interface{}, page int, perPage int) ([]*api.Step, int64, error) {
	e.logger.WithFields(logrus.Fields{
		"build": b.GetNumber(),
	}).Tracef("listing steps for build %d", b.GetNumber())

	// variables to store query results and return value
	count := int64(0)
	s := new([]types.Step)
	steps := []*api.Step{}

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
		WithContext(ctx).
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

		steps = append(steps, tmp.ToAPI())
	}

	return steps, count, nil
}
