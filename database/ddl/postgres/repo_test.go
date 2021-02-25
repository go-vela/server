// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package postgres

import (
	"reflect"
	"testing"
)

func TestPostgres_createRepoService(t *testing.T) {
	// setup types
	want := &Service{
		Create:  CreateRepoTable,
		Indexes: []string{CreateRepoOrgNameIndex},
	}

	// run test
	got := createRepoService()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("createRepoService is %v, want %v", got, want)
	}
}
