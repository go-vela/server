// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package dml

import "github.com/go-vela/server/database/dml/postgres"

// mapFromPostgres is a helper function that converts
// a Postgres DML Map to a common DML Map.
func mapFromPostgres(from *postgres.Map) *Map {
	return &Map{
		BuildService:   serviceFromPostgres(from.BuildService),
		HookService:    serviceFromPostgres(from.HookService),
		LogService:     serviceFromPostgres(from.LogService),
		RepoService:    serviceFromPostgres(from.RepoService),
		SecretService:  serviceFromPostgres(from.SecretService),
		ServiceService: serviceFromPostgres(from.ServiceService),
		StepService:    serviceFromPostgres(from.StepService),
		UserService:    serviceFromPostgres(from.UserService),
		WorkerService:  serviceFromPostgres(from.WorkerService),
	}
}

// serviceFromPostgres is a helper function that converts
// a Postgres DML service to a common DML service.
func serviceFromPostgres(from *postgres.Service) *Service {
	return &Service{
		List:   from.List,
		Select: from.Select,
		Delete: from.Delete,
	}
}
