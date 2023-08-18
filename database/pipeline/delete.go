// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package pipeline

import (
	"context"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// DeletePipeline deletes an existing pipeline from the database.
func (e *engine) DeletePipeline(ctx context.Context, p *library.Pipeline) error {
	e.logger.WithFields(logrus.Fields{
		"pipeline": p.GetCommit(),
	}).Tracef("deleting pipeline %s from the database", p.GetCommit())

	// cast the library type to database type
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#PipelineFromLibrary
	pipeline := database.PipelineFromLibrary(p)

	// send query to the database
	return e.client.
		Table(constants.TablePipeline).
		Delete(pipeline).
		Error
}
