// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package secret

import "github.com/go-vela/types/library"

// Service represents the interface for Vela integrating
// with the different supported secret providers.
type Service interface {
	// Get defines a function that captures a secret.
	Get(string, string, string, string) (*library.Secret, error)
	// List defines a function that captures a list of secrets.
	List(string, string, string, int, int) ([]*library.Secret, error)
	// Count defines a function that counts a list of secrets.
	Count(string, string, string) (int64, error)
	// Create defines a function that creates a new secret.
	Create(string, string, string, *library.Secret) error
	// Update defines a function that updates an existing secret.
	Update(string, string, string, *library.Secret) error
	// Delete defines a function that deletes a secret.
	Delete(string, string, string, string) error

	// TODO: Add convert functions to interface?
}
