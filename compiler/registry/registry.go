// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package registry

import "github.com/go-vela/types/library"

// Service represents the interface for Vela integrating
// with the different supported template registries.
type Service interface {
	// Parse defines a function that creates the
	// registry source object from a template path.
	Parse(string) (*Source, error)

	// Template defines a function that captures the
	// templated pipeline configuration from a repo.
	Template(*library.User, *Source) ([]byte, error)
}
