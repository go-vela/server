// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package sqlite

import (
	"reflect"
	"testing"
)

func TestSqlite_createHookService(t *testing.T) {
	// setup types
	want := &Service{
		List: map[string]string{
			"all":  ListHooks,
			"repo": ListRepoHooks,
		},
		Select: map[string]string{
			"repo": SelectRepoHook,
		},
		Delete: DeleteHook,
	}

	// run test
	got := createHookService()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("createHookService is %v, want %v", got, want)
	}
}
