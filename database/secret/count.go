// SPDX-License-Identifier: Apache-2.0

package secret

import (
	"context"

	"github.com/go-vela/types/constants"
)

// CountSecrets gets the count of all secrets from the database.
func (e *engine) CountSecrets(ctx context.Context) (int64, error) {
	e.logger.Tracef("getting count of all secrets from the database")

	// variable to store query results
	var s int64

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableSecret).
		Count(&s).
		Error

	return s, err
}
