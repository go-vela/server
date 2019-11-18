// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package dml

import "github.com/go-vela/server/database/dml/sqlite"

// mapFromSqlite is a helper function that converts
// a Sqlite DML Map to a common DML Map.
func mapFromSqlite(from *sqlite.Map) *Map {
	return &Map{
		BuildService:   serviceFromSqlite(from.BuildService),
		HookService:    serviceFromSqlite(from.HookService),
		LogService:     serviceFromSqlite(from.LogService),
		RepoService:    serviceFromSqlite(from.RepoService),
		SecretService:  serviceFromSqlite(from.SecretService),
		ServiceService: serviceFromSqlite(from.ServiceService),
		StepService:    serviceFromSqlite(from.StepService),
		UserService:    serviceFromSqlite(from.UserService),
	}
}

// serviceFromSqlite is a helper function that converts
// a Sqlite DML service to a common DML service.
func serviceFromSqlite(from *sqlite.Service) *Service {
	return &Service{
		List:   from.List,
		Select: from.Select,
		Delete: from.Delete,
	}
}
