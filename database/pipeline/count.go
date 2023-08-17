// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package pipeline

import (
	"context"

	"github.com/go-vela/types/constants"
)

// CountPipelines gets the count of all pipelines from the database.
func (e *engine) CountPipelines(ctx context.Context) (int64, error) {
	e.logger.Tracef("getting count of all pipelines from the database")

	// variable to store query results
	var p int64

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TablePipeline).
		Count(&p).
		Error

	return p, err
}
