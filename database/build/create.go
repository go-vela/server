// SPDX-License-Identifier: Apache-2.0

//nolint:dupl // ignore similar code with update.go
package build

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/types"
	"github.com/go-vela/types/constants"
)

// CreateBuild creates a new build in the database.
func (e *engine) CreateBuild(ctx context.Context, b *api.Build) (*api.Build, error) {
	e.logger.WithFields(logrus.Fields{
		"build": b.GetNumber(),
	}).Tracef("creating build %d", b.GetNumber())

	build := types.BuildFromAPI(b)

	err := build.Validate()
	if err != nil {
		return nil, err
	}

	// crop build if any columns are too large
	build = build.Crop()

	// send query to the database
	err = e.client.
		WithContext(ctx).
		Table(constants.TableBuild).
		Create(build).Error
	if err != nil {
		return nil, err
	}

	result := build.ToAPI()
	result.SetRepo(b.GetRepo())

	return result, nil
}
