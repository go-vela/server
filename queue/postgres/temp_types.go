package postgres

import "database/sql"

// todo: (vader) move to types
const (
	BuildsQueueTable        = "builds_queue"
	QueueQueryTransactionID = 123
)

// QueueBuild is the database representation of a build queue item.
type QueueBuild struct {
	ID      sql.NullInt64  `sql:"id"`
	Channel sql.NullString `sql:"channel"`
	Item    []byte         `sql:"item"`
}
