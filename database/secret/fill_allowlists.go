// SPDX-License-Identifier: Apache-2.0

package secret

import (
	"context"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// FillSecretAllowlists fills the allowlists for a slice of secrets.
func (e *Engine) FillSecretsAllowlists(ctx context.Context, secrets []*api.Secret) ([]*api.Secret, error) {
	e.logger.Trace("getting allowlist for secret list")

	// create list of IDs because processing this is faster than making multiple single result queries
	idList := []int64{}

	allowlistMap := make(map[int64][]string, len(secrets))

	for _, s := range secrets {
		idList = append(idList, s.GetID())

		allowlistMap[s.GetID()] = []string{}
	}

	// variable to store query results
	result := new([]types.SecretAllowlist)

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TableSecretRepoAllowlist).
		Where("secret_id IN (?)", idList).
		Find(&result).
		Error
	if err != nil {
		return nil, err
	}

	// iterate through all query results
	for _, record := range *result {
		tmp := record

		allowlistMap[tmp.SecretID.Int64] = append(allowlistMap[tmp.SecretID.Int64], tmp.Repo.String)
	}

	for _, s := range secrets {
		s.SetRepoAllowlist(allowlistMap[s.GetID()])
	}

	return secrets, nil
}
