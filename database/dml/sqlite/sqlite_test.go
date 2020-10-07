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
			List: map[string]string{
				"all":         ListBuilds,
				"repo":        ListRepoBuilds,
				"repoByEvent": ListRepoBuildsByEvent,
				"org":         ListOrgBuilds,
				"orgByEvent":  ListOrgBuildsByEvent,
			},
			Select: map[string]string{
				"repo":                SelectRepoBuild,
				"last":                SelectLastRepoBuild,
				"lastByBranch":        SelectLastRepoBuildByBranch,
				"count":               SelectBuildsCount,
				"countByStatus":       SelectBuildsCountByStatus,
				"countByRepo":         SelectRepoBuildCount,
				"countByRepoAndEvent": SelectRepoBuildCountByEvent,
				"countByOrg":          SelectOrgBuildCount,
				"countByOrgAndEvent":  SelectOrgBuildCountByEvent,
			},
			Delete: DeleteBuild,
		},
		HookService: &Service{
			List: map[string]string{
				"all":  ListHooks,
				"repo": ListRepoHooks,
			},
			Select: map[string]string{
				"count": SelectRepoHookCount,
				"repo":  SelectRepoHook,
				"last":  SelectLastRepoHook,
			},
			Delete: DeleteHook,
		},
		LogService: &Service{
			List: map[string]string{
				"all":   ListLogs,
				"build": ListBuildLogs,
			},
			Select: map[string]string{
				"step":    SelectStepLog,
				"service": SelectServiceLog,
			},
			Delete: DeleteLog,
		},
		RepoService: &Service{
			List: map[string]string{
				"all":  ListRepos,
				"user": ListUserRepos,
				"org":  ListOrgRepos,
			},
			Select: map[string]string{
				"repo":        SelectRepo,
				"count":       SelectReposCount,
				"countByUser": SelectUserReposCount,
			},
			Delete: DeleteRepo,
		},
		SecretService: &Service{
			List: map[string]string{
				"all":    ListSecrets,
				"org":    ListOrgSecrets,
				"repo":   ListRepoSecrets,
				"shared": ListSharedSecrets,
			},
			Select: map[string]string{
				"org":         SelectOrgSecret,
				"repo":        SelectRepoSecret,
				"shared":      SelectSharedSecret,
				"countOrg":    SelectOrgSecretsCount,
				"countRepo":   SelectRepoSecretsCount,
				"countShared": SelectSharedSecretsCount,
			},
			Delete: DeleteSecret,
		},
		ServiceService: &Service{
			List: map[string]string{
				"all":   ListServices,
				"build": ListBuildServices,
			},
			Select: map[string]string{
				"build":          SelectBuildService,
				"count":          SelectBuildServicesCount,
				"count-images":   SelectServiceImagesCount,
				"count-statuses": SelectServiceStatusesCount,
			},
			Delete: DeleteService,
		},
		StepService: &Service{
			List: map[string]string{
				"all":   ListSteps,
				"build": ListBuildSteps,
			},
			Select: map[string]string{
				"build":          SelectBuildStep,
				"count":          SelectBuildStepsCount,
				"count-images":   SelectStepImagesCount,
				"count-statuses": SelectStepStatusesCount,
			},
			Delete: DeleteStep,
		},
		UserService: &Service{
			List: map[string]string{
				"all":  ListUsers,
				"lite": ListLiteUsers,
			},
			Select: map[string]string{
				"user":  SelectUser,
				"name":  SelectUserName,
				"count": SelectUsersCount,
			},
			Delete: DeleteUser,
		},
		WorkerService: &Service{
			List: map[string]string{
				"all": ListWorkers,
			},
			Select: map[string]string{
				"worker":  SelectWorker,
				"address": SelectWorkerByAddress,
				"count":   SelectWorkersCount,
			},
			Delete: DeleteWorker,
		},
	}

	// run test
	got := NewMap()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("NewMap is %v, want %v", got, want)
	}
}
