package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/go-vela/types"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// todo: (vader) test scenarios for load balancing and multiple workers, etc

// Pop defines a function that grabs an
// item off the queue.
func (c *client) Pop(ctx context.Context) (*types.Item, error) {
	c.Logger.Tracef("popping item from postgres queue %s", c.config.Channels)

	// variable to store query results
	qb := new(QueueBuild)

	// create the transaction to combine SELECT and DELETE for concurrency
	transaction := func(_db *gorm.DB) error {
		c.Logger.Trace("trying postgres queue transaction lock")

		// this ID is arbitrary but shared among distributed servers
		err := c.TryTransactionLock(QueueQueryTransactionID, c.Postgres)
		if err != nil {
			c.Logger.Errorf("unable to obtain queue database lock %d", QueueQueryTransactionID)

			return err
		}

		c.Logger.Trace("getting row from the postgres queue")

		// send query to the database and store result in variable

		// todo: (vader) make sure this is performant for a full queue
		err = c.Postgres.
			Table(BuildsQueueTable).
			Where("channel IN ?", c.config.Channels).
			First(qb). // todo: (vader) maybe order by something like... queued_at?
			Error

		if err != nil {
			return err
		}

		c.Logger.Trace("removing row from the postgres queue")

		// todo: (vader) does this many deletions cause bloat that needs to be vacuumed?

		// send query to the database and store result in variable
		err = c.Postgres.
			Table(BuildsQueueTable).
			Delete(qb).
			Error
		if err != nil {
			return err
		}

		c.Logger.Trace("build popped, completing transaction")
		return nil
	}

	c.Logger.Trace("executing postgres queue pop transaction")

	retries := 3
	var err error
	// backoff is timeout / retries
	// var backoff time.Duration = c.config.PopTimeout / time.Duration(retries)
	var backoff time.Duration = 3
	for retry := 0; retry <= retries; retry++ {

		// use a transaction timeout for safety
		timeoutCtx, cancel := context.WithTimeout(ctx, c.config.PopTransactionTimeout*time.Second)
		defer cancel()

		// attempt to execute and commit the transaction
		err = c.Transaction(timeoutCtx, transaction)
		if err != nil {
			// 'ignore' record not found, and try again until the pop query timeout is met
			if errors.Is(err, gorm.ErrRecordNotFound) && retry != retries {
				time.Sleep(backoff * time.Second)
				backoff *= 2
				continue
			}

			logrus.Errorf("error completing build queue transaction retry %d: %s", retry, err)
		}

		// build popped
		if len(qb.Item) > 0 {
			break
		}
	}

	if err != nil {
		// ignore not found errors which indicate the worker has no builds to execute
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.Trace("no build to pop from the queue")
			return nil, nil
		}

		// return any actual errors
		logrus.Errorf("unable to complete build queue transaction after %d retries: %s", retries, err)
		return nil, err
	}

	c.Logger.Tracef("parsing row from postgres builds queue")

	item := new(types.Item)

	// unmarshal result into queue item
	err = json.Unmarshal([]byte(qb.Item), item)
	if err != nil {
		return nil, err
	}

	return item, nil
}
