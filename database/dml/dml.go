// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package dml

import (
	"fmt"

	"github.com/go-vela/server/database/dml/postgres"
	"github.com/go-vela/server/database/dml/sqlite"

	"github.com/go-vela/types/constants"
)

// Service represents the common DML for a table in the database.
type Service struct {
	List   map[string]string
	Select map[string]string
	Delete string
}

// Map represents the common DML services in a struct for lookups.
type Map struct {
	BuildService   *Service
	HookService    *Service
	LogService     *Service
	RepoService    *Service
	SecretService  *Service
	ServiceService *Service
	StepService    *Service
	UserService    *Service
}

// NewMap returns the Map for DML lookups.
func NewMap(name string) (*Map, error) {
	// determine which data manipulation language should be used
	switch name {
	// handle postgres data manipulation language
	case constants.DriverPostgres:
		return mapFromPostgres(postgres.NewMap()), nil
	// handle sqlite data manipulation language
	case constants.DriverSqlite:
		return mapFromSqlite(sqlite.NewMap()), nil
	default:
		return nil, fmt.Errorf("Unrecognized database driver: %s", name)
	}
}
