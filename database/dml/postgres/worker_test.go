// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package postgres

import (
	"reflect"
	"testing"
)

func TestPostgres_createWorkerService(t *testing.T) {
	// setup types
	want := &Service{
		List: map[string]string{
			"all": ListWorkers,
		},
		Select: map[string]string{
			"worker":  SelectWorker,
			"count":   SelectWorkersCount,
			"address": SelectWorkerByAddress,
		},
		Delete: DeleteWorker,
	}

	// run test
	got := createWorkerService()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("createWorkerService is %v, want %v", got, want)
	}
}
