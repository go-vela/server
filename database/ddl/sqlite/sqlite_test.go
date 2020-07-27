// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package sqlite

import (
	"reflect"
	"testing"
)

func TestSqlite_NewMap(t *testing.T) {
	// setup types
	want := &Map{
		BuildService: &Service{
			Create:  CreateBuildTable,
			Indexes: []string{CreateBuildRepoIDIndex, CreateBuildRepoIDNumberIndex, CreateBuildStatusIndex},
		},
		HookService: &Service{
			Create:  CreateHookTable,
			Indexes: []string{CreateHookRepoIDNumberIndex, CreateHookRepoIDIndex},
		},
		LogService: &Service{
			Create:  CreateLogTable,
			Indexes: []string{CreateLogBuildIDIndex, CreateLogStepIDIndex, CreateLogServiceIDIndex},
		},
		RepoService: &Service{
			Create:  CreateRepoTable,
			Indexes: []string{CreateRepoOrgNameIndex, CreateRepoFullNameIndex},
		},
		SecretService: &Service{
			Create: CreateSecretTable,
			Indexes: []string{
				CreateSecretTypeOrgRepo,
				CreateSecretTypeOrgTeam,
				CreateSecretTypeOrg,
				CreateSecretType,
			},
		},
		ServiceService: &Service{
			Create:  CreateServiceTable,
			Indexes: []string{CreateServiceBuildIDNumberIndex},
		},
		StepService: &Service{
			Create:  CreateStepTable,
			Indexes: []string{CreateStepBuildIDNumberIndex},
		},
		UserService: &Service{
			Create:  CreateUserTable,
			Indexes: []string{CreateUserNameIndex, CreateRefreshIndex},
		},
	}

	// run test
	got := NewMap()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("NewMap is %v, want %v", got, want)
	}
}
