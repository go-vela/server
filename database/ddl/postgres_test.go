// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package ddl

import (
	"reflect"
	"testing"

	"github.com/go-vela/server/database/ddl/postgres"
)

func TestDDL_mapFromPostgres(t *testing.T) {
	// setup types
	want := &Map{
		BuildService: &Service{
			Create: []string{postgres.CreateBuildTable, postgres.CreatePayloadColumn},
			Indexes: []string{
				postgres.CreateBuildRepoIDIndex,
				postgres.CreateBuildRepoIDNumberIndex,
				postgres.CreateBuildStatusIndex,
			},
		},
		HookService: &Service{
			Create: []string{postgres.CreateHookTable},
			Indexes: []string{
				postgres.CreateHookRepoIDNumberIndex,
				postgres.CreateHookRepoIDIndex,
			},
		},
		LogService: &Service{
			Create: []string{postgres.CreateLogTable},
			Indexes: []string{
				postgres.CreateLogBuildIDIndex,
				postgres.CreateLogStepIDIndex,
				postgres.CreateLogServiceIDIndex,
			},
		},
		RepoService: &Service{
			Create: []string{postgres.CreateRepoTable},
			Indexes: []string{
				postgres.CreateRepoOrgNameIndex,
				postgres.CreateRepoFullNameIndex,
			},
		},
		SecretService: &Service{
			Create: []string{postgres.CreateSecretTable},
			Indexes: []string{
				postgres.CreateSecretTypeOrgRepo,
				postgres.CreateSecretTypeOrgTeam,
				postgres.CreateSecretTypeOrg,
				postgres.CreateSecretType,
			},
		},
		ServiceService: &Service{
			Create: []string{postgres.CreateServiceTable},
			Indexes: []string{
				postgres.CreateServiceBuildIDNumberIndex,
			},
		},
		StepService: &Service{
			Create: []string{postgres.CreateStepTable},
			Indexes: []string{
				postgres.CreateStepBuildIDNumberIndex,
			},
		},
		UserService: &Service{
			Create: []string{postgres.CreateUserTable},
			Indexes: []string{
				postgres.CreateUserNameIndex,
			},
		},
		WorkerService: &Service{
			Create: []string{postgres.CreateWorkerTable},
			Indexes: []string{
				postgres.CreateWorkersHostnameAddressIndex,
			},
		},
	}

	// run test
	got := mapFromPostgres(postgres.NewMap())

	if !reflect.DeepEqual(got, want) {
		t.Errorf("mapFromPostgres is %v, want %v", got, want)
	}
}

func TestDDL_serviceFromPostgres(t *testing.T) {
	// setup types
	want := &Service{
		Create: []string{postgres.CreateBuildTable, postgres.CreatePayloadColumn},
		Indexes: []string{
			postgres.CreateBuildRepoIDIndex,
			postgres.CreateBuildRepoIDNumberIndex,
			postgres.CreateBuildStatusIndex,
		},
	}

	// run test
	got := serviceFromPostgres(postgres.NewMap().BuildService)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("serviceFromPostgres is %v, want %v", got, want)
	}
}
