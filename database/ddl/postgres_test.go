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
			postgres.CreateBuildStatusIndex,
		},
	}

	// run test
	got := serviceFromPostgres(postgres.NewMap().BuildService)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("serviceFromPostgres is %v, want %v", got, want)
	}
}
