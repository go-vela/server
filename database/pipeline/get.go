// SPDX-License-Identifier: Apache-2.0

package pipeline

import (
	"context"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// GetPipeline gets a pipeline by ID from the database.
func (e *engine) GetPipeline(ctx context.Context, id int64) (*api.Pipeline, error) {
	e.logger.Tracef("getting pipeline %d", id)

	// variable to store query results
	p := new(types.Pipeline)

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TablePipeline).
		Preload("Repo").
		Preload("Repo.Owner").
		Where("id = ?", id).
		Take(p).
		Error
	if err != nil {
		return nil, err
	}

	err = p.Decompress()
	if err != nil {
		return nil, err
	}

	err = p.Repo.Decrypt(e.config.EncryptionKey)
	if err != nil {
		e.logger.Errorf("unable to decrypt repo: %v", err)
	}

	return p.ToAPI(), nil
}
