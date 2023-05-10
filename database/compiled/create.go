// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package compiled

import (
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
)

// CreatePipeline creates a new pipeline in the database.
func (e *engine) CreateCompiled(c *library.Compiled) error {
	// e.logger.WithFields(logrus.Fields{
	// 	"pipeline": p.GetCommit(),
	// }).Tracef("creating pipeline %s in the database", p.GetCommit())

	// cast the library type to database type
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#PipelineFromLibrary
	compiled := database.CompiledFromLibrary(c)

	// validate the necessary fields are populated
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Pipeline.Validate
	err := compiled.Validate()
	if err != nil {
		return err
	}

	// compress data for the pipeline
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Pipeline.Compress
	err = compiled.Compress(e.config.CompressionLevel)
	if err != nil {
		return err
	}

	// send query to the database
	return e.client.
		Table(constants.TableCompiled).
		Create(compiled).
		Error
}
