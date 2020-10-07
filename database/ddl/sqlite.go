// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package ddl

import "github.com/go-vela/server/database/ddl/sqlite"

// mapFromSqlite is a helper function that converts
// a Sqlite DDL Map to a common DDL Map.
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
		WorkerService:  serviceFromSqlite(from.WorkerService),
	}
}

// serviceFromSqlite is a helper function that converts
// a Sqlite DDL service to a common DDL service.
func serviceFromSqlite(from *sqlite.Service) *Service {
	return &Service{
		Create:  from.Create,
		Indexes: from.Indexes,
	}
}
