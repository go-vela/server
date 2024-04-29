// SPDX-License-Identifier: Apache-2.0

package jwk

import (
	"context"

	api "github.com/go-vela/server/api/types"
)

// JWKInterface represents the Vela interface for JWK
// functions with the supported Database backends.
//
//nolint:revive // ignore name stutter
type JWKInterface interface {
	// JWK Data Definition Language Functions
	//
	// https://en.wikipedia.org/wiki/Data_definition_language
	CreateJWKTable(context.Context, string) error

	// JWK Data Manipulation Language Functions
	//
	// https://en.wikipedia.org/wiki/Data_manipulation_language

	// CreateJWK defines a function that creates a JWK.
	CreateJWK(context.Context, api.JWK) error
	// RotateKeys defines a function that rotates JWKs.
	RotateKeys(context.Context) error
	// ListJWKs defines a function that lists all JWKs configured.
	ListJWKs(context.Context) ([]api.JWK, error)
	// GetJWK defines a function that gets a JWK by the provided key ID.
	GetActiveJWK(context.Context, string) (api.JWK, error)
}
