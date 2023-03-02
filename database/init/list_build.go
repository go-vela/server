// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package init

import (
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// ListInitsForBuild gets a list of inits by build ID from the database.
func (e *engine) ListInitsForBuild(b *library.Build, page, perPage int) ([]*library.Init, int64, error) {
	e.logger.WithFields(logrus.Fields{
		"build": b.GetNumber(),
	}).Tracef("listing inits for build %d from the database", b.GetNumber())

	// variables to store query results and return value
	count := int64(0)
	h := new([]database.Init)
	inits := []*library.Init{}

	// count the results
	count, err := e.CountInitsForBuild(b)
	if err != nil {
		return nil, 0, err
	}

	// short-circuit if there are no results
	if count == 0 {
		return inits, 0, nil
	}

	// calculate offset for pagination through results
	offset := perPage * (page - 1)

	// send query to the database and store result in variable
	err = e.client.
		Table(constants.TableInit).
		Where("build_id = ?", b.GetID()).
		Order("id DESC").
		Limit(perPage).
		Offset(offset).
		Find(&h).
		Error
	if err != nil {
		return nil, count, err
	}

	// iterate through all query results
	for _, init := range *h {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := init

		// convert query result to library type
		//
		// https://pkg.go.dev/github.com/go-vela/types/database#Init.ToLibrary
		inits = append(inits, tmp.ToLibrary())
	}

	return inits, count, nil
}
