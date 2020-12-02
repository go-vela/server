// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
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
			Create:  []string{CreateBuildTable, CreatePayloadColumn},
			Indexes: []string{CreateBuildRepoIDIndex, CreateBuildRepoIDNumberIndex, CreateBuildStatusIndex},
		},
		HookService: &Service{
			Create:  []string{CreateHookTable},
			Indexes: []string{CreateHookRepoIDNumberIndex, CreateHookRepoIDIndex},
		},
		LogService: &Service{
			Create:  []string{CreateLogTable},
			Indexes: []string{CreateLogBuildIDIndex, CreateLogStepIDIndex, CreateLogServiceIDIndex},
		},
		RepoService: &Service{
			Create:  []string{CreateRepoTable},
			Indexes: []string{CreateRepoOrgNameIndex, CreateRepoFullNameIndex},
		},
		SecretService: &Service{
			Create: []string{CreateSecretTable},
			Indexes: []string{
				CreateSecretTypeOrgRepo,
				CreateSecretTypeOrgTeam,
				CreateSecretTypeOrg,
				CreateSecretType,
			},
		},
		ServiceService: &Service{
			Create:  []string{CreateServiceTable},
			Indexes: []string{CreateServiceBuildIDNumberIndex},
		},
		StepService: &Service{
			Create:  []string{CreateStepTable},
			Indexes: []string{CreateStepBuildIDNumberIndex},
		},
		UserService: &Service{
			Create:  []string{CreateUserTable},
			Indexes: []string{CreateUserNameIndex},
		},
		WorkerService: &Service{
			Create:  []string{CreateWorkerTable},
			Indexes: []string{CreateWorkersHostnameAddressIndex},
		},
	}

	// run test
	got := NewMap()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("NewMap is %v, want %v", got, want)
	}
}
