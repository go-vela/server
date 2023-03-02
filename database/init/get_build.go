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

// GetInitForBuild gets a init by build ID and number from the database.
func (e *engine) GetInitForBuild(b *library.Build, number int) (*library.Init, error) {
	e.logger.WithFields(logrus.Fields{
		"init":  number,
		"build": b.GetNumber(),
	}).Tracef("getting init %d/%d from the database", b.GetNumber(), number)

	// variable to store query results
	h := new(database.Init)

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableInit).
		Where("build_id = ?", b.GetID()).
		Where("number = ?", number).
		Take(h).
		Error
	if err != nil {
		return nil, err
	}

	// return the init
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Init.ToLibrary
	return h.ToLibrary(), nil
}
