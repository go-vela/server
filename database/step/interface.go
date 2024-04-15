// SPDX-License-Identifier: Apache-2.0

package step

import (
	"context"

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
	CreateStepTable(context.Context, string) error

	// Step Data Manipulation Language Functions
	//
	// https://en.wikipedia.org/wiki/Data_manipulation_language

	// CleanSteps defines a function that sets running or pending steps to error status before a given created time.
	CleanSteps(context.Context, string, int64) (int64, error)
	// CountSteps defines a function that gets the count of all steps.
	CountSteps(context.Context) (int64, error)
	// CountStepsForBuild defines a function that gets the count of steps by build ID.
	CountStepsForBuild(context.Context, *api.Build, map[string]interface{}) (int64, error)
	// CreateStep defines a function that creates a new step.
	CreateStep(context.Context, *library.Step) (*library.Step, error)
	// DeleteStep defines a function that deletes an existing step.
	DeleteStep(context.Context, *library.Step) error
	// GetStep defines a function that gets a step by ID.
	GetStep(context.Context, int64) (*library.Step, error)
	// GetStepForBuild defines a function that gets a step by number and build ID.
	GetStepForBuild(context.Context, *api.Build, int) (*library.Step, error)
	// ListSteps defines a function that gets a list of all steps.
	ListSteps(ctx context.Context) ([]*library.Step, error)
	// ListStepsForBuild defines a function that gets a list of steps by build ID.
	ListStepsForBuild(context.Context, *api.Build, map[string]interface{}, int, int) ([]*library.Step, int64, error)
	// ListStepImageCount defines a function that gets a list of all step images and the count of their occurrence.
	ListStepImageCount(context.Context) (map[string]float64, error)
	// ListStepStatusCount defines a function that gets a list of all step statuses and the count of their occurrence.
	ListStepStatusCount(context.Context) (map[string]float64, error)
	// UpdateStep defines a function that updates an existing step.
	UpdateStep(context.Context, *library.Step) (*library.Step, error)
}
