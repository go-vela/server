// SPDX-License-Identifier: Apache-2.0

package user

import (
	"context"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/types"
	"github.com/go-vela/types/constants"
)

// ListLiteUsers gets a lite (only: id, name) list of users from the database.
//
//nolint:lll // ignore long line length due to variable names
func (e *engine) ListLiteUsers(ctx context.Context, page, perPage int) ([]*api.User, int64, error) {
	e.logger.Trace("listing lite users")

	// variables to store query results and return values
	count := int64(0)
	u := new([]types.User)
	users := []*api.User{}

	// count the results
	count, err := e.CountUsers(ctx)
	if err != nil {
		return users, 0, err
	}

	// short-circuit if there are no results
	if count == 0 {
		return users, 0, nil
	}

	// calculate offset for pagination through results
	offset := perPage * (page - 1)

	err = e.client.
		WithContext(ctx).
		Table(constants.TableUser).
		Select("id", "name").
		Limit(perPage).
		Offset(offset).
		Find(&u).
		Error
	if err != nil {
		return nil, count, err
	}

	// iterate through all query results
	for _, user := range *u {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := user

		// convert query result to API type
		users = append(users, tmp.ToAPI())
	}

	return users, count, nil
}
