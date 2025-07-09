// SPDX-License-Identifier: Apache-2.0

package user

import (
	"context"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// ListLiteUsers gets a lite (only: id, name) list of users from the database.
//
//nolint:lll // ignore long line length due to variable names
func (e *Engine) ListLiteUsers(ctx context.Context, page, perPage int) ([]*api.User, error) {
	e.logger.Trace("listing lite users")

	// variables to store query results and return values
	u := new([]types.User)
	users := []*api.User{}

	// calculate offset for pagination through results
	offset := perPage * (page - 1)

	err := e.client.
		WithContext(ctx).
		Table(constants.TableUser).
		Select("id", "name").
		Limit(perPage).
		Offset(offset).
		Find(&u).
		Error
	if err != nil {
		return nil, err
	}

	// iterate through all query results
	for _, user := range *u {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := user

		// convert query result to API type
		users = append(users, tmp.ToAPI())
	}

	return users, nil
}
