// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package postgres

import (
	"reflect"
	"testing"
)

func TestPostgres_createStepService(t *testing.T) {
	// setup types
	want := &Service{
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
	}

	// run test
	got := createStepService()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("createStepService is %v, want %v", got, want)
	}
}
