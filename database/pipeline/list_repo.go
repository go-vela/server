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

// ListPipelinesForRepo gets a list of pipelines by repo ID from the database.
//
//nolint:lll // ignore long line length due to variable names
func (e *engine) ListPipelinesForRepo(ctx context.Context, r *library.Repo, page, perPage int) ([]*library.Pipeline, int64, error) {
	e.logger.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
	}).Tracef("listing pipelines for repo %s from the database", r.GetFullName())

	// variables to store query results and return values
	count := int64(0)
	p := new([]database.Pipeline)
	pipelines := []*library.Pipeline{}

	// count the results
	count, err := e.CountPipelinesForRepo(ctx, r)
	if err != nil {
		return pipelines, 0, err
	}

	// short-circuit if there are no results
	if count == 0 {
		return pipelines, 0, nil
	}

	// calculate offset for pagination through results
	offset := perPage * (page - 1)

	err = e.client.
		Table(constants.TablePipeline).
		Where("repo_id = ?", r.GetID()).
		Limit(perPage).
		Offset(offset).
		Find(&p).
		Error
	if err != nil {
		return nil, count, err
	}

	// iterate through all query results
	for _, pipeline := range *p {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := pipeline

		// decompress data for the pipeline
		//
		// https://pkg.go.dev/github.com/go-vela/types/database#Pipeline.Decompress
		err = tmp.Decompress()
		if err != nil {
			return nil, count, err
		}

		// convert query result to library type
		//
		// https://pkg.go.dev/github.com/go-vela/types/database#Pipeline.ToLibrary
		pipelines = append(pipelines, tmp.ToLibrary())
	}

	return pipelines, count, nil
}
