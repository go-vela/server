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
func (e *Engine) ListStepsForBuild(ctx context.Context, b *api.Build, filters map[string]any, page int, perPage int) ([]*api.Step, error) {
	e.logger.WithFields(logrus.Fields{
		"build": b.GetNumber(),
	}).Tracef("listing steps for build %d", b.GetNumber())

	// variables to store query results and return value
	s := new([]types.Step)
	steps := []*api.Step{}

	// calculate offset for pagination through results
	offset := perPage * (page - 1)

	// send query to the database and store result in variable
	err := e.client.
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
		return nil, err
	}

	// iterate through all query results
	for _, step := range *s {
		steps = append(steps, step.ToAPI())
	}

	return steps, nil
}
