// SPDX-License-Identifier: Apache-2.0

package jwk

import (
	"context"

	"github.com/lestrrat-go/jwx/v2/jwk"

	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// ListJWKs gets a list of all configured JWKs from the database.
func (e *engine) ListJWKs(_ context.Context) (jwk.Set, error) {
	e.logger.Trace("listing all jwks from the database")

	k := new([]types.JWK)
	keySet := jwk.NewSet()

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableJWK).
		Find(&k).
		Error
	if err != nil {
		return nil, err
	}

	// iterate through all query results
	for _, key := range *k {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := key

		// convert query result to API type
		keySet.AddKey(tmp.ToAPI())
	}

	return keySet, nil
}
