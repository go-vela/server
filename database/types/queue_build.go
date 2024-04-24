// SPDX-License-Identifier: Apache-2.0

package types

import (
	"database/sql"

	api "github.com/go-vela/server/api/types"
)

// QueueBuild is the database representation of the builds in the queue.
type QueueBuild struct {
	Status   sql.NullString `sql:"status"`
	Number   sql.NullInt32  `sql:"number"`
	Created  sql.NullInt64  `sql:"created"`
	FullName sql.NullString `sql:"full_name"`
}

// ToAPI converts the QueueBuild type
// to a API QueueBuild type.
func (b *QueueBuild) ToAPI() *api.QueueBuild {
	buildQueue := new(api.QueueBuild)

	buildQueue.SetStatus(b.Status.String)
	buildQueue.SetNumber(b.Number.Int32)
	buildQueue.SetCreated(b.Created.Int64)
	buildQueue.SetFullName(b.FullName.String)

	return buildQueue
}

// QueueBuildFromAPI converts the API QueueBuild type
// to a database build queue type.
func QueueBuildFromAPI(b *api.QueueBuild) *QueueBuild {
	buildQueue := &QueueBuild{
		Status:   sql.NullString{String: b.GetStatus(), Valid: true},
		Number:   sql.NullInt32{Int32: b.GetNumber(), Valid: true},
		Created:  sql.NullInt64{Int64: b.GetCreated(), Valid: true},
		FullName: sql.NullString{String: b.GetFullName(), Valid: true},
	}

	return buildQueue
}
