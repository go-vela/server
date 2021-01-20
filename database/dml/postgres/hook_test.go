// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
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
	}

	// run test
	got := createHookService()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("createHookService is %v, want %v", got, want)
	}
}
