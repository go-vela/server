// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package init

import (
	"github.com/go-vela/types/library"
)

// InitService represents the Vela interface for init
// functions with the supported Database backends.
//
//nolint:revive // ignore name stutter
type InitService interface {
	// Init Data Definition Language Functions
	//
	// https://en.wikipedia.org/wiki/Data_definition_language

	// CreateInitsIndexes defines a function that creates the indexes for the inits table.
	CreateInitsIndexes() error
	// CreateInitsTable defines a function that creates the inits table.
	CreateInitsTable(string) error

	// Init Data Manipulation Language Functions
	//
	// https://en.wikipedia.org/wiki/Data_manipulation_language

	// CountInits defines a function that gets the count of all inits.
	CountInits() (int64, error)
	// CountInitsForBuild defines a function that gets the count of inits by build ID.
	CountInitsForBuild(*library.Build) (int64, error)
	// CreateInit defines a function that creates a new init.
	CreateInit(*library.Init) error
	// DeleteInit defines a function that deletes an existing init.
	DeleteInit(*library.Init) error
	// GetInit defines a function that gets a init by ID.
	GetInit(int64) (*library.Init, error)
	// GetInitForBuild defines a function that gets a init by build ID and number.
	GetInitForBuild(*library.Build, int) (*library.Init, error)
	// ListInits defines a function that gets a list of all inits.
	ListInits() ([]*library.Init, error)
	// ListInitsForBuild defines a function that gets a list of inits by build ID.
	ListInitsForBuild(*library.Build, int, int) ([]*library.Init, int64, error)
	// UpdateInit defines a function that updates an existing init.
	UpdateInit(*library.Init) error
}
