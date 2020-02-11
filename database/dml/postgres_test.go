// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package dml

import (
	"reflect"
	"testing"

	"github.com/go-vela/server/database/dml/postgres"
)

func TestDML_mapFromPostgres(t *testing.T) {
	// setup types
	want := &Map{
		BuildService: &Service{
			List: map[string]string{
				"all":         postgres.ListBuilds,
				"repo":        postgres.ListRepoBuilds,
				"repoByEvent": postgres.ListRepoBuildsByEvent,
			},
			Select: map[string]string{
				"repo":                postgres.SelectRepoBuild,
				"last":                postgres.SelectLastRepoBuild,
				"count":               postgres.SelectBuildsCount,
				"countByStatus":       postgres.SelectBuildsCountByStatus,
				"countByRepo":         postgres.SelectRepoBuildCount,
				"countByRepoAndEvent": postgres.SelectRepoBuildCountByEvent,
			},
			Delete: postgres.DeleteBuild,
		},
		HookService: &Service{
			List: map[string]string{
				"all":  postgres.ListHooks,
				"repo": postgres.ListRepoHooks,
			},
			Select: map[string]string{
				"count": postgres.SelectRepoHookCount,
				"repo":  postgres.SelectRepoHook,
				"last":  postgres.SelectLastRepoHook,
			},
			Delete: postgres.DeleteHook,
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
				"org":  postgres.ListOrgRepos,
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
				"build":        postgres.SelectBuildService,
				"count":        postgres.SelectBuildServicesCount,
				"count-images": postgres.SelectServiceImagesCount,
			},
			Delete: postgres.DeleteService,
		},
		StepService: &Service{
			List: map[string]string{
				"all":   postgres.ListSteps,
				"build": postgres.ListBuildSteps,
			},
			Select: map[string]string{
				"build":          postgres.SelectBuildStep,
				"count":          postgres.SelectBuildStepsCount,
				"count-images":   postgres.SelectStepImagesCount,
				"count-statuses": postgres.SelectStepStatusesCount,
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
	got := mapFromPostgres(postgres.NewMap())

	if !reflect.DeepEqual(got, want) {
		t.Errorf("mapFromPostgres is %v, want %v", got, want)
	}
}

func TestDML_serviceFromPostgres(t *testing.T) {
	// setup types
	want := &Service{
		List: map[string]string{
			"all":         postgres.ListBuilds,
			"repo":        postgres.ListRepoBuilds,
			"repoByEvent": postgres.ListRepoBuildsByEvent,
		},
		Select: map[string]string{
			"repo":                postgres.SelectRepoBuild,
			"last":                postgres.SelectLastRepoBuild,
			"count":               postgres.SelectBuildsCount,
			"countByStatus":       postgres.SelectBuildsCountByStatus,
			"countByRepo":         postgres.SelectRepoBuildCount,
			"countByRepoAndEvent": postgres.SelectRepoBuildCountByEvent,
		},
		Delete: postgres.DeleteBuild,
	}

	// run test
	got := serviceFromPostgres(postgres.NewMap().BuildService)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("serviceFromPostgres is %v, want %v", got, want)
	}
}
