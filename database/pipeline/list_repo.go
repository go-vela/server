// SPDX-License-Identifier: Apache-2.0

package pipeline

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// ListPipelinesForRepo gets a list of pipelines by repo ID from the database.
//
//nolint:lll // ignore long line length due to variable names
func (e *Engine) ListPipelinesForRepo(ctx context.Context, r *api.Repo, page, perPage int) ([]*api.Pipeline, error) {
	e.logger.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
	}).Tracef("listing pipelines for repo %s", r.GetFullName())

	// variables to store query results and return values
	p := new([]types.Pipeline)
	pipelines := []*api.Pipeline{}

	// calculate offset for pagination through results
	offset := perPage * (page - 1)

	err := e.client.
		WithContext(ctx).
		Table(constants.TablePipeline).
		Where("repo_id = ?", r.GetID()).
		Limit(perPage).
		Offset(offset).
		Find(&p).
		Error
	if err != nil {
		return nil, err
	}

	// iterate through all query results
	for _, pipeline := range *p {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := pipeline

		err = tmp.Decompress()
		if err != nil {
			return nil, err
		}

		result := tmp.ToAPI()
		result.SetRepo(r)

		pipelines = append(pipelines, result)
	}

	return pipelines, nil
}
