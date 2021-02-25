// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
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
				postgres.CreateBuildStatusIndex,
			},
		},
		HookService: &Service{
			Create: postgres.CreateHookTable,
			Indexes: []string{
				postgres.CreateHookRepoIDIndex,
			},
		},
		LogService: &Service{
			Create: postgres.CreateLogTable,
			Indexes: []string{
				postgres.CreateLogBuildIDIndex,
			},
		},
		RepoService: &Service{
			Create: postgres.CreateRepoTable,
			Indexes: []string{
				postgres.CreateRepoOrgNameIndex,
			},
		},
		SecretService: &Service{
			Create: postgres.CreateSecretTable,
			Indexes: []string{
				postgres.CreateSecretTypeOrgRepo,
				postgres.CreateSecretTypeOrgTeam,
				postgres.CreateSecretTypeOrg,
			},
		},
		ServiceService: &Service{
			Create:  postgres.CreateServiceTable,
			Indexes: []string{},
		},
		StepService: &Service{
			Create:  postgres.CreateStepTable,
			Indexes: []string{},
		},
		UserService: &Service{
			Create: postgres.CreateUserTable,
			Indexes: []string{
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
				sqlite.CreateBuildStatusIndex,
			},
		},
		HookService: &Service{
			Create: sqlite.CreateHookTable,
			Indexes: []string{
				sqlite.CreateHookRepoIDIndex,
			},
		},
		LogService: &Service{
			Create: sqlite.CreateLogTable,
			Indexes: []string{
				sqlite.CreateLogBuildIDIndex,
			},
		},
		RepoService: &Service{
			Create: sqlite.CreateRepoTable,
			Indexes: []string{
				sqlite.CreateRepoOrgNameIndex,
			},
		},
		SecretService: &Service{
			Create: sqlite.CreateSecretTable,
			Indexes: []string{
				sqlite.CreateSecretTypeOrgRepo,
				sqlite.CreateSecretTypeOrgTeam,
				sqlite.CreateSecretTypeOrg,
			},
		},
		ServiceService: &Service{
			Create:  sqlite.CreateServiceTable,
			Indexes: []string{},
		},
		StepService: &Service{
			Create:  sqlite.CreateStepTable,
			Indexes: []string{},
		},
		UserService: &Service{
			Create: sqlite.CreateUserTable,
			Indexes: []string{
				sqlite.CreateRefreshIndex,
			},
		},
		WorkerService: &Service{
			Create: sqlite.CreateWorkerTable,
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
