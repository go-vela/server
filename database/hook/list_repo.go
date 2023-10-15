// SPDX-License-Identifier: Apache-2.0

package hook

import (
	"context"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// ListHooksForRepo gets a list of hooks by repo ID from the database.
func (e *engine) ListHooksForRepo(ctx context.Context, r *library.Repo, page, perPage int) ([]*library.Hook, int64, error) {
	e.logger.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
	}).Tracef("listing hooks for repo %s from the database", r.GetFullName())

	// variables to store query results and return value
	count := int64(0)
	h := new([]database.Hook)
	hooks := []*library.Hook{}

	// count the results
	count, err := e.CountHooksForRepo(ctx, r)
	if err != nil {
		return nil, 0, err
	}

	// short-circuit if there are no results
	if count == 0 {
		return hooks, 0, nil
	}

	// calculate offset for pagination through results
	offset := perPage * (page - 1)

	// send query to the database and store result in variable
	err = e.client.
		Table(constants.TableHook).
		Where("repo_id = ?", r.GetID()).
		Order("id DESC").
		Limit(perPage).
		Offset(offset).
		Find(&h).
		Error
	if err != nil {
		return nil, count, err
	}

	// iterate through all query results
	for _, hook := range *h {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := hook

		// convert query result to library type
		//
		// https://pkg.go.dev/github.com/go-vela/types/database#Hook.ToLibrary
		hooks = append(hooks, tmp.ToLibrary())
	}

	return hooks, count, nil
}
