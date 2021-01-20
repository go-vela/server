// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
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
			Create: postgres.CreateBuildTable,
			Indexes: []string{
				postgres.CreateBuildRepoIDIndex,
				postgres.CreateBuildRepoIDNumberIndex,
				postgres.CreateBuildStatusIndex,
			},
		},
		HookService: &Service{
			Create: postgres.CreateHookTable,
			Indexes: []string{
				postgres.CreateHookRepoIDNumberIndex,
				postgres.CreateHookRepoIDIndex,
			},
		},
		LogService: &Service{
			Create: postgres.CreateLogTable,
			Indexes: []string{
				postgres.CreateLogBuildIDIndex,
				postgres.CreateLogStepIDIndex,
				postgres.CreateLogServiceIDIndex,
			},
		},
		RepoService: &Service{
			Create: postgres.CreateRepoTable,
			Indexes: []string{
				postgres.CreateRepoOrgNameIndex,
				postgres.CreateRepoFullNameIndex,
			},
		},
		SecretService: &Service{
			Create: postgres.CreateSecretTable,
			Indexes: []string{
				postgres.CreateSecretTypeOrgRepo,
				postgres.CreateSecretTypeOrgTeam,
				postgres.CreateSecretTypeOrg,
				postgres.CreateSecretType,
			},
		},
		ServiceService: &Service{
			Create: postgres.CreateServiceTable,
			Indexes: []string{
				postgres.CreateServiceBuildIDNumberIndex,
			},
		},
		StepService: &Service{
			Create: postgres.CreateStepTable,
			Indexes: []string{
				postgres.CreateStepBuildIDNumberIndex,
			},
		},
		UserService: &Service{
			Create: postgres.CreateUserTable,
			Indexes: []string{
				postgres.CreateUserNameIndex,
				postgres.CreateRefreshIndex,
			},
		},
		WorkerService: &Service{
			Create: postgres.CreateWorkerTable,
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
		Create: postgres.CreateBuildTable,
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
