// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package pipeline

import (
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// DeletePipeline deletes an existing pipeline from the database.
func (e *engine) DeletePipeline(p *library.Pipeline) error {
	e.logger.WithFields(logrus.Fields{
		"pipeline": p.GetNumber(),
	}).Tracef("deleting pipeline %d from the database", p.GetNumber())

	// cast to database type
	pipeline := database.PipelineFromLibrary(p)

	// send query to the database
	return e.client.
		Table(constants.TablePipeline).
		Delete(pipeline).
		Error
}
