package postgres

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/go-vela/types"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/pipeline"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const (
	BuildsQueueTable = "builds_queue"
)

// QueueBuild is the database representation of a build queue item.
type QueueBuild struct {
	ID    sql.NullInt64  `sql:"id"`
	Route sql.NullString `sql:"route"`
	Item  []byte         `sql:"item"`
}

// Length tallies all items present in the configured channels in the queue.
func (c *client) Length(ctx context.Context) (int64, error) {
	c.Logger.Tracef("reading length of all configured channels in queue")

	return 0, errors.New("postgres queue.Length not implemented")
}

// Pop defines a function that grabs an
// item off the queue.
func (c *client) Pop(ctx context.Context) (*types.Item, error) {
	c.Logger.Tracef("popping item from postgres queue %s", c.config.Channels)

	// variable to store query results
	qb := new(QueueBuild)

	tx := func(_db *gorm.DB) error {
		c.Logger.Trace("tx: trying postgres queue transaction lock")

		// this ID is arbitrary but shared among distributed servers
		lockID := 123
		err := c.TryTransactionLock(lockID, c.Postgres)
		if err != nil {
			c.Logger.Infof("tx: unable to obtain queue database lock %d", lockID)

			return err
		}

		c.Logger.Trace("tx: getting row from the postgres queue")

		// send query to the database and store result in variable
		err = c.Postgres.
			Table(BuildsQueueTable).
			Where("route IN ?", c.config.Channels).
			First(qb).
			Error
		if err != nil {
			return err
		}

		c.Logger.Trace("tx: removing row from the postgres queue")

		// send query to the database and store result in variable
		err = c.Postgres.
			Table(BuildsQueueTable).
			Delete(qb).
			Error
		if err != nil {
			return err
		}

		c.Logger.Trace("tx: build popped, completing transaction")
		return nil
	}

	c.Logger.Trace("commiting transaction")

	// attempt to execute and commit the transaction
	err := c.Transaction(tx)
	if err != nil {
		logrus.Errorf("unable to complete build queue transaction: %s", err)
	}

	c.Logger.Tracef("parsing row from postgres queue into executable item")

	item := new(types.Item)

	// unmarshal result into queue item
	err = json.Unmarshal([]byte(qb.Item), item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

// Push defines a function that publishes an
// item to the specified route in the queue.
func (c *client) Push(ctx context.Context, channel string, item []byte) error {
	c.Logger.Tracef("pushing item to postgres queue %s", channel)

	qb := &QueueBuild{
		Route: sql.NullString{String: channel, Valid: true},
		Item:  item,
	}

	// send query to the database
	return c.Postgres.
		Table(BuildsQueueTable).
		Create(qb).
		Error
}

// Route defines a function that decides which
// channel a build gets placed within the queue.
func (c *client) Route(w *pipeline.Worker) (string, error) {
	c.Logger.Tracef("deciding route from queue channels %s", c.config.Channels)

	// create buffer to store route
	buf := bytes.Buffer{}

	// if pipline does not specify route information return default
	//
	// https://github.com/go-vela/types/blob/main/constants/queue.go#L10
	if w.Empty() {
		return constants.DefaultRoute, nil
	}

	// append flavor to route
	if !strings.EqualFold(strings.ToLower(w.Flavor), "") {
		buf.WriteString(fmt.Sprintf(":%s", w.Flavor))
	}

	// append platform to route
	if !strings.EqualFold(strings.ToLower(w.Platform), "") {
		buf.WriteString(fmt.Sprintf(":%s", w.Platform))
	}

	route := strings.TrimLeft(buf.String(), ":")

	for _, r := range c.config.Channels {
		if strings.EqualFold(route, r) {
			return route, nil
		}
	}

	return "", fmt.Errorf("invalid route %s provided", route)
}

func (c *client) Transaction(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) error {
	return c.Postgres.Transaction(fc, opts...)
}

func (c *client) TryTransactionLock(txID int, tx *gorm.DB) error {
	// use transaction db if provided
	//  allows for nested transactions
	db := c.Postgres
	if tx != nil {
		db = tx
	}

	err := db.Exec("SELECT PG_TRY_ADVISORY_XACT_LOCK(?);", txID).Error
	if err != nil {
		return err
	}

	return nil
}
