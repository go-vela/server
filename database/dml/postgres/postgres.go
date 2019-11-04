// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package postgres

// Service represents the Postgres DML for a table in the database.
type Service struct {
	List   map[string]string
	Select map[string]string
	Delete string
}

// Map represents the Postgres DML services in a struct for lookups.
type Map struct {
	BuildService   *Service
	LogService     *Service
	RepoService    *Service
	SecretService  *Service
	ServiceService *Service
	StepService    *Service
	UserService    *Service
}

// NewMap returns the Postgres Map for DML lookups.
func NewMap() *Map {
	return &Map{
		BuildService:   createBuildService(),
		LogService:     createLogService(),
		RepoService:    createRepoService(),
		SecretService:  createSecretService(),
		ServiceService: createServiceService(),
		StepService:    createStepService(),
		UserService:    createUserService(),
	}
}
