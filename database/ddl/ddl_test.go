// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package ddl

import (
	"reflect"
	"testing"

	"github.com/go-vela/server/database/ddl/postgres"
	"github.com/go-vela/server/database/ddl/sqlite"

	"github.com/go-vela/types/constants"
)

func TestDDL_NewMap_Postgres(t *testing.T) {
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
	got, err := NewMap(constants.DriverPostgres)

	if err != nil {
		t.Errorf("NewMap returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("NewMap is %v, want %v", got, want)
	}
}

func TestDDL_NewMap_Sqlite(t *testing.T) {
	// setup types
	want := &Map{
		BuildService: &Service{
			Create: []string{sqlite.CreateBuildTable},
			Indexes: []string{
				sqlite.CreateBuildRepoIDIndex,
				sqlite.CreateBuildRepoIDNumberIndex,
				sqlite.CreateBuildStatusIndex,
			},
		},
		HookService: &Service{
			Create: []string{sqlite.CreateHookTable},
			Indexes: []string{
				sqlite.CreateHookRepoIDNumberIndex,
				sqlite.CreateHookRepoIDIndex,
			},
		},
		LogService: &Service{
			Create: []string{sqlite.CreateLogTable},
			Indexes: []string{
				sqlite.CreateLogBuildIDIndex,
				sqlite.CreateLogStepIDIndex,
				sqlite.CreateLogServiceIDIndex,
			},
		},
		RepoService: &Service{
			Create: []string{sqlite.CreateRepoTable},
			Indexes: []string{
				sqlite.CreateRepoOrgNameIndex,
				sqlite.CreateRepoFullNameIndex,
			},
		},
		SecretService: &Service{
			Create: []string{sqlite.CreateSecretTable},
			Indexes: []string{
				sqlite.CreateSecretTypeOrgRepo,
				sqlite.CreateSecretTypeOrgTeam,
				sqlite.CreateSecretTypeOrg,
				sqlite.CreateSecretType,
			},
		},
		ServiceService: &Service{
			Create: []string{sqlite.CreateServiceTable},
			Indexes: []string{
				sqlite.CreateServiceBuildIDNumberIndex,
			},
		},
		StepService: &Service{
			Create: []string{sqlite.CreateStepTable},
			Indexes: []string{
				sqlite.CreateStepBuildIDNumberIndex,
			},
		},
		UserService: &Service{
			Create: []string{sqlite.CreateUserTable},
			Indexes: []string{
				sqlite.CreateUserNameIndex,
			},
		},
		WorkerService: &Service{
			Create: []string{sqlite.CreateWorkerTable},
			Indexes: []string{
				sqlite.CreateWorkersHostnameAddressIndex,
			},
		},
	}

	// run test
	got, err := NewMap(constants.DriverSqlite)

	if err != nil {
		t.Errorf("NewMap returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("NewMap is %v, want %v", got, want)
	}
}

func TestDDL_NewMap_None(t *testing.T) {
	// run test
	got, err := NewMap("")

	if err == nil {
		t.Errorf("NewMap should have returned err")
	}

	if got != nil {
		t.Errorf("NewMap is %v, want nil", got)
	}
}
