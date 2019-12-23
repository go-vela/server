// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package ddl

import (
	"fmt"

	"github.com/go-vela/server/database/ddl/postgres"
	"github.com/go-vela/server/database/ddl/sqlite"

	"github.com/go-vela/types/constants"
)

// Service represents the common DDL for a table in the database.
type Service struct {
	Create  string
	Indexes []string
}

// Map represents the common DDL services in a struct for lookups.
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

// NewMap returns the Map for DDL lookups.
func NewMap(name string) (*Map, error) {
	// determine which data definition language should be used
	switch name {
	// handle postgres data definition language
	case constants.DriverPostgres:
		return mapFromPostgres(postgres.NewMap()), nil
	// handle sqlite data definition language
	case constants.DriverSqlite:
		return mapFromSqlite(sqlite.NewMap()), nil
	default:
		return nil, fmt.Errorf("unrecognized database driver: %s", name)
	}
}
