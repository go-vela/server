// SPDX-License-Identifier: Apache-2.0

package secret

import (
	"context"

	api "github.com/go-vela/server/api/types"
)

// Service represents the interface for Vela integrating
// with the different supported secret providers.
type Service interface {
	// Service Interface Functions

	// Driver defines a function that outputs
	// the configured source driver.
	Driver() string

	// Get defines a function that captures a secret.
	Get(context.Context, string, string, string, string) (*api.Secret, error)
	// List defines a function that captures a list of secrets.
	List(context.Context, string, string, string, int, int, []string) ([]*api.Secret, error)
	// Count defines a function that counts a list of secrets.
	Count(context.Context, string, string, string, []string) (int64, error)
	// Create defines a function that creates a new secret.
	Create(context.Context, string, string, string, *api.Secret) (*api.Secret, error)
	// Update defines a function that updates an existing secret.
	Update(context.Context, string, string, string, *api.Secret) (*api.Secret, error)
	// Delete defines a function that deletes a secret.
	Delete(context.Context, string, string, string, string) error

	// TODO: Add convert functions to interface?
}
