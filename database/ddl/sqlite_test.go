// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package ddl

import (
	"reflect"
	"testing"

	"github.com/go-vela/server/database/ddl/sqlite"
)

func TestDDL_mapFromSqlite(t *testing.T) {
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
	got := mapFromSqlite(sqlite.NewMap())

	if !reflect.DeepEqual(got, want) {
		t.Errorf("mapFromSqlite is %v, want %v", got, want)
	}
}

func TestDDL_serviceFromSqlite(t *testing.T) {
	// setup types
	want := &Service{
		Create: []string{sqlite.CreateBuildTable},
		Indexes: []string{
			sqlite.CreateBuildRepoIDIndex,
			sqlite.CreateBuildRepoIDNumberIndex,
			sqlite.CreateBuildStatusIndex,
		},
	}

	// run test
	got := serviceFromSqlite(sqlite.NewMap().BuildService)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("serviceFromSqlite is %v, want %v", got, want)
	}
}
