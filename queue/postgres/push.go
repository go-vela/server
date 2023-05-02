package postgres

import (
	"context"
	"database/sql"
)

// Push defines a function that publishes an
// item to the specified channel in the queue.
func (c *client) Push(ctx context.Context, channel string, item []byte) error {
	c.Logger.Tracef("pushing item to postgres queue %s", channel)

	// todo: (vader) move to types
	qb := &QueueBuild{
		Channel: sql.NullString{String: channel, Valid: true},
		Item:    item,
	}

	// todo: (vader) add compression

	// send query to the database
	return c.Postgres.
		Table(BuildsQueueTable).
		Create(qb).
		Error
}
