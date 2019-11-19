// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package postgres

import (
	"reflect"
	"testing"
)

func TestPostgres_createHookService(t *testing.T) {
	// setup types
	want := &Service{
		Create:  CreateHookTable,
		Indexes: []string{CreateHookRepoIDBuildIDIndex},
	}

	// run test
	got := createHookService()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("createHookService is %v, want %v", got, want)
	}
}
