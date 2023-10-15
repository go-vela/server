// SPDX-License-Identifier: Apache-2.0

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
