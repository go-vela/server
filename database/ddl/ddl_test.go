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
			Create: sqlite.CreateBuildTable,
			Indexes: []string{
				sqlite.CreateBuildRepoIDIndex,
				sqlite.CreateBuildRepoIDNumberIndex,
				sqlite.CreateBuildStatusIndex,
			},
		},
		HookService: &Service{
			Create: sqlite.CreateHookTable,
			Indexes: []string{
				sqlite.CreateHookRepoIDNumberIndex,
				sqlite.CreateHookRepoIDIndex,
			},
		},
		LogService: &Service{
			Create: sqlite.CreateLogTable,
			Indexes: []string{
				sqlite.CreateLogBuildIDIndex,
				sqlite.CreateLogStepIDIndex,
				sqlite.CreateLogServiceIDIndex,
			},
		},
		RepoService: &Service{
			Create: sqlite.CreateRepoTable,
			Indexes: []string{
				sqlite.CreateRepoOrgNameIndex,
				sqlite.CreateRepoFullNameIndex,
			},
		},
		SecretService: &Service{
			Create: sqlite.CreateSecretTable,
			Indexes: []string{
				sqlite.CreateSecretTypeOrgRepo,
				sqlite.CreateSecretTypeOrgTeam,
				sqlite.CreateSecretTypeOrg,
				sqlite.CreateSecretType,
			},
		},
		ServiceService: &Service{
			Create: sqlite.CreateServiceTable,
			Indexes: []string{
				sqlite.CreateServiceBuildIDNumberIndex,
			},
		},
		StepService: &Service{
			Create: sqlite.CreateStepTable,
			Indexes: []string{
				sqlite.CreateStepBuildIDNumberIndex,
			},
		},
		UserService: &Service{
			Create: sqlite.CreateUserTable,
			Indexes: []string{
				sqlite.CreateUserNameIndex,
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
