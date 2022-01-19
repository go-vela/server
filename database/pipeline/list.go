// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package pipeline

import (
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
)

// ListPipelines gets a list of all pipelines from the database.
func (e *engine) ListPipelines() ([]*library.Pipeline, error) {
	e.logger.Trace("listing all pipelines from the database")

	// variable to store query results
	p := new([]database.Pipeline)

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TablePipeline).
		Find(&p).
		Error

	// variable we want to return
	pipelines := []*library.Pipeline{}
	// iterate through all query results
	for _, pipeline := range *p {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := pipeline

		// decompress data for the pipeline
		//
		// https://pkg.go.dev/github.com/go-vela/types/database#Pipeline.Decompress
		err = tmp.Decompress()
		if err != nil {
			return nil, err
		}

		// convert query result to library type
		pipelines = append(pipelines, tmp.ToLibrary())
	}

	return pipelines, err
}
