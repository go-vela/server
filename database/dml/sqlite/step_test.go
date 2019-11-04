// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package sqlite

import (
	"reflect"
	"testing"
)

func TestSqlite_createStepService(t *testing.T) {
	// setup types
	want := &Service{
		List: map[string]string{
			"all":   ListSteps,
			"build": ListBuildSteps,
		},
		Select: map[string]string{
			"build": SelectBuildStep,
			"count": SelectBuildStepsCount,
		},
		Delete: DeleteStep,
	}

	// run test
	got := createStepService()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("createStepService is %v, want %v", got, want)
	}
}
