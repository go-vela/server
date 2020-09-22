// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package sqlite

import (
	"reflect"
	"testing"
)

func TestSqlite_createServiceService(t *testing.T) {
	// setup types
	want := &Service{
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
	}

	// run test
	got := createServiceService()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("createServiceService is %v, want %v", got, want)
	}
}
