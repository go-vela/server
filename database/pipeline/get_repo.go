// SPDX-License-Identifier: Apache-2.0

package pipeline

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/types"
	"github.com/go-vela/types/constants"
)

// GetPipelineForRepo gets a pipeline by number and repo ID from the database.
func (e *engine) GetPipelineForRepo(ctx context.Context, commit string, r *api.Repo) (*api.Pipeline, error) {
	e.logger.WithFields(logrus.Fields{
		"org":      r.GetOrg(),
		"pipeline": commit,
		"repo":     r.GetName(),
	}).Tracef("getting pipeline %s/%s", r.GetFullName(), commit)

	// variable to store query results
	p := new(types.Pipeline)

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TablePipeline).
		Where("repo_id = ?", r.GetID()).
		Where("\"commit\" = ?", commit).
		Take(p).
		Error
	if err != nil {
		return nil, err
	}

	err = p.Decompress()
	if err != nil {
		return nil, err
	}

	result := p.ToAPI()
	result.SetRepo(r)

	return result, nil
}
