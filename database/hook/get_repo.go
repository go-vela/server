// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package hook

import (
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// GetHookForRepo gets a hook by repo ID and number from the database.
func (e *engine) GetHookForRepo(r *library.Repo, number int) (*library.Hook, error) {
	e.logger.WithFields(logrus.Fields{
		"hook": number,
		"org":  r.GetOrg(),
		"repo": r.GetName(),
	}).Tracef("getting hook %s/%d from the database", r.GetFullName(), number)

	// variable to store query results
	h := new(database.Hook)

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableHook).
		Where("repo_id = ?", r.GetID()).
		Where("number = ?", number).
		Take(h).
		Error
	if err != nil {
		return nil, err
	}

	// return the hook
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Hook.ToLibrary
	return h.ToLibrary(), nil
}
