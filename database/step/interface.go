// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package step

import (
	"github.com/go-vela/types/library"
)

// StepInterface represents the Vela interface for step
// functions with the supported Database backends.
//
//nolint:revive // ignore name stutter
type StepInterface interface {
	// Step Data Definition Language Functions
	//
	// https://en.wikipedia.org/wiki/Data_definition_language

	// CreateStepTable defines a function that creates the steps table.
	CreateStepTable(string) error

	// Step Data Manipulation Language Functions
	//
	// https://en.wikipedia.org/wiki/Data_manipulation_language

	// CountSteps defines a function that gets the count of all steps.
	CountSteps() (int64, error)
	// CountStepsForBuild defines a function that gets the count of steps by build ID.
	CountStepsForBuild(*library.Build, map[string]interface{}) (int64, error)
	// CreateStep defines a function that creates a new step.
	CreateStep(*library.Step) error
	// DeleteStep defines a function that deletes an existing step.
	DeleteStep(*library.Step) error
	// GetStep defines a function that gets a step by ID.
	GetStep(int64) (*library.Step, error)
	// GetStepForBuild defines a function that gets a step by number and build ID.
	GetStepForBuild(*library.Build, int) (*library.Step, error)
	// ListSteps defines a function that gets a list of all steps.
	ListSteps() ([]*library.Step, error)
	// ListStepsForBuild defines a function that gets a list of steps by build ID.
	ListStepsForBuild(*library.Build, map[string]interface{}, int, int) ([]*library.Step, int64, error)
	// ListStepImageCount defines a function that gets a list of all step images and the count of their occurrence.
	ListStepImageCount() (map[string]float64, error)
	// ListStepStatusCount defines a function that gets a list of all step statuses and the count of their occurrence.
	ListStepStatusCount() (map[string]float64, error)
	// UpdateStep defines a function that updates an existing step.
	UpdateStep(*library.Step) error
}
