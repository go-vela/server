// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
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

// LastPipelineForRepo gets the last pipeline by repo ID from the database.
func (e *engine) LastPipelineForRepo(r *library.Repo) (*library.Pipeline, error) {
	e.logger.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
	}).Tracef("getting last pipeline for repo %s from the database", r.GetFullName())

	// variable to store query results
	p := new(database.Pipeline)

	// send query to the database and store result in variable
	result := e.client.
		Table(constants.TablePipeline).
		Where("repo_id = ?", r.GetID()).
		Order("number DESC").
		Limit(1).
		Scan(p)

	// check if the query returned a record not found error or no rows were returned
	if errors.Is(result.Error, gorm.ErrRecordNotFound) || result.RowsAffected == 0 {
		// the record will not exist if it's a new repo
		return nil, nil
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
