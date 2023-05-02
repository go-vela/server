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

	transaction := func(_db *gorm.DB) error {
		c.Logger.Trace("trying postgres queue transaction lock")

		// this ID is arbitrary but shared among distributed servers
		err := c.TryTransactionLock(QueueQueryTransactionID, c.Postgres)
		if err != nil {
			c.Logger.Errorf("unable to obtain queue database lock %d", QueueQueryTransactionID)

			return err
		}

		c.Logger.Trace("getting row from the postgres queue")

		// todo: pop query timeout manually, to retry the query with a backoff
		retries := 3

		var backoff time.Duration = 1

		for retry := 0; retry <= retries; retry++ {
			// send query to the database and store result in variable
			err = c.Postgres.
				Table(BuildsQueueTable).
				Where("channel IN ?", c.config.Channels).
				First(qb).
				Error

			// 'ignore' record not found, and try again until the pop query timeout is met
			if errors.Is(err, gorm.ErrRecordNotFound) && retry != retries {
				time.Sleep(backoff * time.Second)
				backoff *= 2
			} else {
				// actual error occurred
				break
			}
		}

		if err != nil {
			return err
		}

		c.Logger.Trace("removing row from the postgres queue")

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

	timeoutCtx, cancel := context.WithTimeout(ctx, c.config.PopTransactionTimeout*time.Second)
	defer cancel()

	// attempt to execute and commit the transaction
	err := c.Transaction(timeoutCtx, transaction)
	if err != nil {
		logrus.Errorf("unable to complete build queue transaction: %s", err)
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
