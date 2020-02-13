// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package postgres

import (
	"reflect"
	"testing"
)

func TestPostgres_createBuildService(t *testing.T) {
	// setup types
	want := &Service{
		List: map[string]string{
			"all":         ListBuilds,
			"repo":        ListRepoBuilds,
			"repoByEvent": ListRepoBuildsByEvent,
		},
		Select: map[string]string{
			"repo":                SelectRepoBuild,
			"last":                SelectLastRepoBuild,
			"count":               SelectBuildsCount,
			"countByStatus":       SelectBuildsCountByStatus,
			"countByRepo":         SelectRepoBuildCount,
			"countByRepoAndEvent": SelectRepoBuildCountByEvent,
		},
		Delete: DeleteBuild,
	}

	// run test
	got := createBuildService()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("createBuildService is %v, want %v", got, want)
	}
}
