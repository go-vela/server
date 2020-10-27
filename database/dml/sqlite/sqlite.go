// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package sqlite

// Service represents the Sqlite DML for a table in the database.
type Service struct {
	List   map[string]string
	Select map[string]string
	Delete string
}

// Map represents the Sqlite DML services in a struct for lookups.
type Map struct {
	BuildService   *Service
	HookService    *Service
	LogService     *Service
	RepoService    *Service
	SecretService  *Service
	ServiceService *Service
	StepService    *Service
	UserService    *Service
	WorkerService  *Service
}

// NewMap returns the Sqlite Map for DML lookups.
func NewMap() *Map {
	return &Map{
		BuildService:   createBuildService(),
		HookService:    createHookService(),
		LogService:     createLogService(),
		RepoService:    createRepoService(),
		SecretService:  createSecretService(),
		ServiceService: createServiceService(),
		StepService:    createStepService(),
		UserService:    createUserService(),
		WorkerService:  createWorkerService(),
	}
}
