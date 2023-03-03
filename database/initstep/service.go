// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package initstep

import (
	"github.com/go-vela/types/library"
)

// InitStepService represents the Vela interface for init step
// functions with the supported Database backends.
//
//nolint:revive // ignore name stutter
type InitStepService interface {
	// InitStep Data Definition Language Functions
	//
	// https://en.wikipedia.org/wiki/Data_definition_language

	// CreateInitStepIndexes defines a function that creates the indexes for the initsteps table.
	CreateInitStepIndexes() error
	// CreateInitStepTable defines a function that creates the initsteps table.
	CreateInitStepTable(string) error

	// InitStep Data Manipulation Language Functions
	//
	// https://en.wikipedia.org/wiki/Data_manipulation_language

	// CountInitSteps defines a function that gets the count of all InitSteps.
	CountInitSteps() (int64, error)
	// CountInitStepsForBuild defines a function that gets the count of InitSteps by build ID.
	CountInitStepsForBuild(*library.Build) (int64, error)
	// CreateInitStep defines a function that creates a new InitStep.
	CreateInitStep(*library.InitStep) error
	// DeleteInitStep defines a function that deletes an existing InitStep.
	DeleteInitStep(*library.InitStep) error
	// GetInitStep defines a function that gets a InitStep by ID.
	GetInitStep(int64) (*library.InitStep, error)
	// GetInitStepForBuild defines a function that gets a InitStep by build ID and number.
	GetInitStepForBuild(*library.Build, int) (*library.InitStep, error)
	// ListInitSteps defines a function that gets a list of all InitSteps.
	ListInitSteps() ([]*library.InitStep, error)
	// ListInitStepsForBuild defines a function that gets a list of InitSteps by build ID.
	ListInitStepsForBuild(*library.Build, int, int) ([]*library.InitStep, int64, error)
	// UpdateInitStep defines a function that updates an existing InitStep.
	UpdateInitStep(*library.InitStep) error
}
