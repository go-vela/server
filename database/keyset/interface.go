// SPDX-License-Identifier: Apache-2.0

package keyset

import (
	"context"

	api "github.com/go-vela/server/api/types"
)

// KeySetInterface represents the Vela interface for key set
// functions with the supported Database backends.
type KeySetInterface interface {
	// KeySet Data Definition Language Functions
	//
	// https://en.wikipedia.org/wiki/Data_definition_language
	CreateKeySetTable(context.Context, string) error

	// KeySet Data Manipulation Language Functions
	//
	// https://en.wikipedia.org/wiki/Data_manipulation_language

	// CreateKeySet defines a function that creates a key set.
	CreateKeySet(context.Context, api.JWK) error
	// DeleteKeySet defines a function that gets and deletes a key set.
	RotateKeys(context.Context) error
	// ListKeySets defines a function that lists all key sets configured.
	ListKeySets(context.Context) ([]api.JWK, error)
	// GetKeySet defines a function that gets a key set by the provided key ID.
	GetActiveKeySet(context.Context, string) (api.JWK, error)
}
