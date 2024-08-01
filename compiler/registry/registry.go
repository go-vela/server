// SPDX-License-Identifier: Apache-2.0

package registry

import (
	api "github.com/go-vela/server/api/types"
)

// Service represents the interface for Vela integrating
// with the different supported template registries.
type Service interface {
	// Parse defines a function that creates the
	// registry source object from a template path.
	Parse(string) (*Source, error)

	// Template defines a function that captures the
	// templated pipeline configuration from a repo.
	Template(*api.User, *Source) ([]byte, error)
}
