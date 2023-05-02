package postgres

import "context"

// Length tallies all items present in the configured channels in the queue.
func (c *client) Length(ctx context.Context) (int64, error) {
	c.Logger.Tracef("reading length of all configured channels in queue")

	total := int64(0)

	// send query to the database to capture queue length
	err := c.Postgres.
		Table(BuildsQueueTable).
		Where("channel IN ?", c.config.Channels).
		Count(&total).
		Error

	if err != nil {
		return 0, err
	}

	return total, nil
}
