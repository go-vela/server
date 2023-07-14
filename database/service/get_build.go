// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package service

import (
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// GetServiceForBuild gets a service by number and build ID from the database.
func (e *engine) GetServiceForBuild(b *library.Build, number int) (*library.Service, error) {
	e.logger.WithFields(logrus.Fields{
		"build":   b.GetNumber(),
		"service": number,
	}).Tracef("getting service %d from the database", number)

	// variable to store query results
	s := new(database.Service)

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableService).
		Where("build_id = ?", b.GetID()).
		Where("number = ?", number).
		Take(s).
		Error
	if err != nil {
		return nil, err
	}

	// return the service
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Service.ToLibrary
	return s.ToLibrary(), nil
}
