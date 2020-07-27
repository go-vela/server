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
				sqlite.CreateRefreshIndex,
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
		Create: sqlite.CreateBuildTable,
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
