// SPDX-License-Identifier: Apache-2.0

package step

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// CreateStep creates a new step in the database.
func (e *engine) CreateStep(ctx context.Context, s *api.Step) (*api.Step, error) {
	e.logger.WithFields(logrus.Fields{
		"step": s.GetNumber(),
	}).Tracef("creating step %s in the database", s.GetName())

	step := types.StepFromAPI(s)

	err := step.Validate()
	if err != nil {
		return nil, err
	}

	// send query to the database
	result := e.client.
		WithContext(ctx).
		Table(constants.TableStep).
		Create(step)

	return step.ToAPI(), result.Error
}
