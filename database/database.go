// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package database

import (
	"fmt"

	"github.com/go-vela/types/constants"

	"github.com/sirupsen/logrus"
)

// New creates and returns a Vela service capable of
// integrating with the configured database provider.
//
// Currently the following database providers are supported:
//
// * Postgres
// * Sqlite
// .
func New(s *Setup) (Service, error) {
	// validate the setup being provided
	//
	// https://pkg.go.dev/github.com/go-vela/server/database?tab=doc#Setup.Validate
	err := s.Validate()
	if err != nil {
		return nil, err
	}

	logrus.Debug("creating database service from setup")
	// process the database driver being provided
	switch s.Driver {
	case constants.DriverPostgres:
		// handle the Postgres database driver being provided
		//
		// https://pkg.go.dev/github.com/go-vela/server/database?tab=doc#Setup.Postgres
		return s.Postgres()
	case constants.DriverSqlite:
		// handle the Sqlite database driver being provided
		//
		// https://pkg.go.dev/github.com/go-vela/server/database?tab=doc#Setup.Sqlite
		return s.Sqlite()
	default:
		// handle an invalid database driver being provided
		return nil, fmt.Errorf("invalid database driver provided: %s", s.Driver)
	}
}
