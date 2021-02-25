// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package postgres

import (
	"reflect"
	"testing"
)

func TestPostgres_NewMap(t *testing.T) {
	// setup types
	want := &Map{
		BuildService: &Service{
			Create:  CreateBuildTable,
			Indexes: []string{CreateBuildRepoIDIndex, CreateBuildStatusIndex},
		},
		HookService: &Service{
			Create:  CreateHookTable,
			Indexes: []string{CreateHookRepoIDIndex},
		},
		LogService: &Service{
			Create:  CreateLogTable,
			Indexes: []string{CreateLogBuildIDIndex},
		},
		RepoService: &Service{
			Create:  CreateRepoTable,
			Indexes: []string{CreateRepoOrgNameIndex},
		},
		SecretService: &Service{
			Create: CreateSecretTable,
			Indexes: []string{
				CreateSecretTypeOrgRepo,
				CreateSecretTypeOrgTeam,
				CreateSecretTypeOrg,
			},
		},
		ServiceService: &Service{
			Create:  CreateServiceTable,
			Indexes: []string{},
		},
		StepService: &Service{
			Create:  CreateStepTable,
			Indexes: []string{},
		},
		UserService: &Service{
			Create:  CreateUserTable,
			Indexes: []string{CreateRefreshIndex},
		},
		WorkerService: &Service{
			Create:  CreateWorkerTable,
			Indexes: []string{CreateWorkersHostnameAddressIndex},
		},
	}

	// run test
	got := NewMap()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("NewMap is %v, want %v", got, want)
	}
}
