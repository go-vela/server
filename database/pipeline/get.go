// SPDX-License-Identifier: Apache-2.0

package pipeline

import (
	"context"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
)

// GetPipeline gets a pipeline by ID from the database.
func (e *engine) GetPipeline(ctx context.Context, id int64) (*library.Pipeline, error) {
	e.logger.Tracef("getting pipeline %d from the database", id)

	// variable to store query results
	p := new(database.Pipeline)

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TablePipeline).
		Where("id = ?", id).
		Take(p).
		Error
	if err != nil {
		return nil, err
	}

	// decompress data for the pipeline
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Pipeline.Decompress
	err = p.Decompress()
	if err != nil {
		return nil, err
	}

	// return the decompressed pipeline
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Pipeline.ToLibrary
	return p.ToLibrary(), nil
}
