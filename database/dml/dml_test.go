// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package dml

import (
	"reflect"
	"testing"

	"github.com/go-vela/server/database/dml/postgres"
	"github.com/go-vela/server/database/dml/sqlite"

	"github.com/go-vela/types/constants"
)

func TestDML_NewMap_Postgres(t *testing.T) {
	// setup types
	want := &Map{
		BuildService: &Service{
			List: map[string]string{
				"all":  postgres.ListBuilds,
				"repo": postgres.ListRepoBuilds,
			},
			Select: map[string]string{
				"repo":          postgres.SelectRepoBuild,
				"last":          postgres.SelectLastRepoBuild,
				"count":         postgres.SelectBuildsCount,
				"countByStatus": postgres.SelectBuildsCountByStatus,
				"countByRepo":   postgres.SelectRepoBuildCount,
			},
			Delete: postgres.DeleteBuild,
		},
		LogService: &Service{
			List: map[string]string{
				"all":   postgres.ListLogs,
				"build": postgres.ListBuildLogs,
			},
			Select: map[string]string{
				"step":    postgres.SelectStepLog,
				"service": postgres.SelectServiceLog,
			},
			Delete: postgres.DeleteLog,
		},
		RepoService: &Service{
			List: map[string]string{
				"all":  postgres.ListRepos,
				"user": postgres.ListUserRepos,
			},
			Select: map[string]string{
				"repo":        postgres.SelectRepo,
				"count":       postgres.SelectReposCount,
				"countByUser": postgres.SelectUserReposCount,
			},
			Delete: postgres.DeleteRepo,
		},
		SecretService: &Service{
			List: map[string]string{
				"all":    postgres.ListSecrets,
				"org":    postgres.ListOrgSecrets,
				"repo":   postgres.ListRepoSecrets,
				"shared": postgres.ListSharedSecrets,
			},
			Select: map[string]string{
				"org":         postgres.SelectOrgSecret,
				"repo":        postgres.SelectRepoSecret,
				"shared":      postgres.SelectSharedSecret,
				"countOrg":    postgres.SelectOrgSecretsCount,
				"countRepo":   postgres.SelectRepoSecretsCount,
				"countShared": postgres.SelectSharedSecretsCount,
			},
			Delete: postgres.DeleteSecret,
		},
		ServiceService: &Service{
			List: map[string]string{
				"all":   postgres.ListServices,
				"build": postgres.ListBuildServices,
			},
			Select: map[string]string{
				"build": postgres.SelectBuildService,
				"count": postgres.SelectBuildServicesCount,
			},
			Delete: postgres.DeleteService,
		},
		StepService: &Service{
			List: map[string]string{
				"all":   postgres.ListSteps,
				"build": postgres.ListBuildSteps,
			},
			Select: map[string]string{
				"build": postgres.SelectBuildStep,
				"count": postgres.SelectBuildStepsCount,
			},
			Delete: postgres.DeleteStep,
		},
		UserService: &Service{
			List: map[string]string{
				"all":  postgres.ListUsers,
				"lite": postgres.ListLiteUsers,
			},
			Select: map[string]string{
				"user":  postgres.SelectUser,
				"name":  postgres.SelectUserName,
				"count": postgres.SelectUsersCount,
			},
			Delete: postgres.DeleteUser,
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

func TestDML_NewMap_Sqlite(t *testing.T) {
	// setup types
	want := &Map{
		BuildService: &Service{
			List: map[string]string{
				"all":  sqlite.ListBuilds,
				"repo": sqlite.ListRepoBuilds,
			},
			Select: map[string]string{
				"repo":          sqlite.SelectRepoBuild,
				"last":          sqlite.SelectLastRepoBuild,
				"count":         sqlite.SelectBuildsCount,
				"countByStatus": sqlite.SelectBuildsCountByStatus,
				"countByRepo":   sqlite.SelectRepoBuildCount,
			},
			Delete: sqlite.DeleteBuild,
		},
		LogService: &Service{
			List: map[string]string{
				"all":   sqlite.ListLogs,
				"build": sqlite.ListBuildLogs,
			},
			Select: map[string]string{
				"step":    sqlite.SelectStepLog,
				"service": sqlite.SelectServiceLog,
			},
			Delete: sqlite.DeleteLog,
		},
		RepoService: &Service{
			List: map[string]string{
				"all":  sqlite.ListRepos,
				"user": sqlite.ListUserRepos,
			},
			Select: map[string]string{
				"repo":        sqlite.SelectRepo,
				"count":       sqlite.SelectReposCount,
				"countByUser": sqlite.SelectUserReposCount,
			},
			Delete: sqlite.DeleteRepo,
		},
		SecretService: &Service{
			List: map[string]string{
				"all":    sqlite.ListSecrets,
				"org":    sqlite.ListOrgSecrets,
				"repo":   sqlite.ListRepoSecrets,
				"shared": sqlite.ListSharedSecrets,
			},
			Select: map[string]string{
				"org":         sqlite.SelectOrgSecret,
				"repo":        sqlite.SelectRepoSecret,
				"shared":      sqlite.SelectSharedSecret,
				"countOrg":    sqlite.SelectOrgSecretsCount,
				"countRepo":   sqlite.SelectRepoSecretsCount,
				"countShared": sqlite.SelectSharedSecretsCount,
			},
			Delete: sqlite.DeleteSecret,
		},
		ServiceService: &Service{
			List: map[string]string{
				"all":   sqlite.ListServices,
				"build": sqlite.ListBuildServices,
			},
			Select: map[string]string{
				"build": sqlite.SelectBuildService,
				"count": sqlite.SelectBuildServicesCount,
			},
			Delete: sqlite.DeleteService,
		},
		StepService: &Service{
			List: map[string]string{
				"all":   sqlite.ListSteps,
				"build": sqlite.ListBuildSteps,
			},
			Select: map[string]string{
				"build": sqlite.SelectBuildStep,
				"count": sqlite.SelectBuildStepsCount,
			},
			Delete: sqlite.DeleteStep,
		},
		UserService: &Service{
			List: map[string]string{
				"all":  sqlite.ListUsers,
				"lite": sqlite.ListLiteUsers,
			},
			Select: map[string]string{
				"user":  sqlite.SelectUser,
				"name":  sqlite.SelectUserName,
				"count": sqlite.SelectUsersCount,
			},
			Delete: sqlite.DeleteUser,
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

func TestDML_NewMap_None(t *testing.T) {
	// run test
	got, err := NewMap("")

	if err == nil {
		t.Errorf("NewMap should have returned err")
	}

	if got != nil {
		t.Errorf("NewMap is %v, want nil", got)
	}
}
