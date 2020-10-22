// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package dml

import (
	"reflect"
	"testing"

	"github.com/go-vela/server/database/dml/sqlite"
)

func TestDML_mapFromSqlite(t *testing.T) {
	// setup types
	want := &Map{
		BuildService: &Service{
			List: map[string]string{
				"all":         sqlite.ListBuilds,
				"repo":        sqlite.ListRepoBuilds,
				"repoByEvent": sqlite.ListRepoBuildsByEvent,
				"org":         sqlite.ListOrgBuilds,
				"orgByEvent":  sqlite.ListOrgBuildsByEvent,
			},
			Select: map[string]string{
				"repo":                sqlite.SelectRepoBuild,
				"last":                sqlite.SelectLastRepoBuild,
				"lastByBranch":        sqlite.SelectLastRepoBuildByBranch,
				"count":               sqlite.SelectBuildsCount,
				"countByStatus":       sqlite.SelectBuildsCountByStatus,
				"countByRepo":         sqlite.SelectRepoBuildCount,
				"countByRepoAndEvent": sqlite.SelectRepoBuildCountByEvent,
				"countByOrg":          sqlite.SelectOrgBuildCount,
				"countByOrgAndEvent":  sqlite.SelectOrgBuildCountByEvent,
			},
			Delete: sqlite.DeleteBuild,
		},
		HookService: &Service{
			List: map[string]string{
				"all":  sqlite.ListHooks,
				"repo": sqlite.ListRepoHooks,
			},
			Select: map[string]string{
				"count": sqlite.SelectRepoHookCount,
				"repo":  sqlite.SelectRepoHook,
				"last":  sqlite.SelectLastRepoHook,
			},
			Delete: sqlite.DeleteHook,
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
				"org":  sqlite.ListOrgRepos,
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
				"build":          sqlite.SelectBuildService,
				"count":          sqlite.SelectBuildServicesCount,
				"count-images":   sqlite.SelectServiceImagesCount,
				"count-statuses": sqlite.SelectServiceStatusesCount,
			},
			Delete: sqlite.DeleteService,
		},
		StepService: &Service{
			List: map[string]string{
				"all":   sqlite.ListSteps,
				"build": sqlite.ListBuildSteps,
			},
			Select: map[string]string{
				"build":          sqlite.SelectBuildStep,
				"count":          sqlite.SelectBuildStepsCount,
				"count-images":   sqlite.SelectStepImagesCount,
				"count-statuses": sqlite.SelectStepStatusesCount,
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
		WorkerService: &Service{
			List: map[string]string{
				"all": sqlite.ListWorkers,
			},
			Select: map[string]string{
				"worker":  sqlite.SelectWorker,
				"address": sqlite.SelectWorkerByAddress,
				"count":   sqlite.SelectWorkersCount,
			},
			Delete: sqlite.DeleteWorker,
		},
	}

	// run test
	got := mapFromSqlite(sqlite.NewMap())

	if !reflect.DeepEqual(got, want) {
		t.Errorf("mapFromSqlite is %v, want %v", got, want)
	}
}

func TestDML_serviceFromSqlite(t *testing.T) {
	// setup types
	want := &Service{
		List: map[string]string{
			"all":         sqlite.ListBuilds,
			"repo":        sqlite.ListRepoBuilds,
			"repoByEvent": sqlite.ListRepoBuildsByEvent,
			"org":         sqlite.ListOrgBuilds,
			"orgByEvent":  sqlite.ListOrgBuildsByEvent,
		},
		Select: map[string]string{
			"repo":                sqlite.SelectRepoBuild,
			"last":                sqlite.SelectLastRepoBuild,
			"lastByBranch":        sqlite.SelectLastRepoBuildByBranch,
			"count":               sqlite.SelectBuildsCount,
			"countByStatus":       sqlite.SelectBuildsCountByStatus,
			"countByRepo":         sqlite.SelectRepoBuildCount,
			"countByRepoAndEvent": sqlite.SelectRepoBuildCountByEvent,
			"countByOrg":          sqlite.SelectOrgBuildCount,
			"countByOrgAndEvent":  sqlite.SelectOrgBuildCountByEvent,
		},
		Delete: sqlite.DeleteBuild,
	}

	// run test
	got := serviceFromSqlite(sqlite.NewMap().BuildService)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("serviceFromSqlite is %v, want %v", got, want)
	}
}
