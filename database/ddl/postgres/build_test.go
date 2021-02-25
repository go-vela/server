// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
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
		Create:  CreateBuildTable,
		Indexes: []string{CreateBuildRepoIDIndex, CreateBuildStatusIndex},
	}

	// run test
	got := createBuildService()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("createBuildService is %v, want %v", got, want)
	}
}
