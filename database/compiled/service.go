// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package compiled

import "github.com/go-vela/types/library"

// CompiledService represents the Vela interface for compiled
// functions with the supported Database backends.
//
//nolint:revive // ignore name stutter
type CompiledService interface {
	// Compiled Data Definition Language Functions
	//
	// https://en.wikipedia.org/wiki/Data_definition_language
	CreateCompiledTable(string) error

	// // Compiled Data Manipulation Language Functions
	// //
	// // https://en.wikipedia.org/wiki/Data_manipulation_language

	// CreateCompiled defines a function that creates a compiled build.
	CreateCompiled(*library.Compiled) error
	// PopCompiled defines a function that gets and deletes a compiled build.
	PopCompiled(int64) (*library.Compiled, error)
}
