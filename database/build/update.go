// SPDX-License-Identifier: Apache-2.0

//nolint:dupl // ignore similar code with create.go
package build

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// UpdateBuild updates an existing build in the database.
func (e *engine) UpdateBuild(ctx context.Context, b *api.Build) (*api.Build, error) {
	e.logger.WithFields(logrus.Fields{
		"build": b.GetNumber(),
	}).Tracef("updating build %d", b.GetNumber())

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
		Save(build).Error
	if err != nil {
		return nil, err
	}

	result := build.ToAPI()
	result.SetRepo(b.GetRepo())

	return result, nil
}
