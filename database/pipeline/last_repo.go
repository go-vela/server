// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package pipeline

import (
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
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
	err := e.client.
		Table(constants.TablePipeline).
		Where("repo_id = ?", r.GetID()).
		Order("number DESC").
		Limit(1).
		Scan(p).
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
