// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package pipeline

import (
	"errors"

	"github.com/sirupsen/logrus"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"gorm.io/gorm"
)

// GetPipelineForRepo gets a pipeline by number and repo ID from the database.
func (e *engine) GetPipelineForRepo(number int, r *library.Repo) (*library.Pipeline, error) {
	e.logger.WithFields(logrus.Fields{
		"org":      r.GetOrg(),
		"pipeline": number,
		"repo":     r.GetName(),
	}).Tracef("getting pipeline %s/%d from the database", r.GetFullName(), number)

	// variable to store query results
	p := new(database.Pipeline)

	// send query to the database and store result in variable
	result := e.client.
		Table(constants.TablePipeline).
		Where("repo_id = ?", r.GetID()).
		Where("number = ?", number).
		Limit(1).
		Scan(p)

	// check if the query returned a record not found error or no rows were returned
	if errors.Is(result.Error, gorm.ErrRecordNotFound) || result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	// decompress data for the pipeline
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Pipeline.Decompress
	err := p.Decompress()
	if err != nil {
		return nil, err
	}

	// return the decompressed pipeline
	return p.ToLibrary(), result.Error
}
