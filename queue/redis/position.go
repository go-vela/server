// SPDX-License-Identifier: Apache-2.0

package redis

import (
	"context"
	"strconv"

	api "github.com/go-vela/server/api/types"
)

// Position returns the position of the item in the queue.
func (c *Client) Position(ctx context.Context, build *api.Build) int64 {
	c.Logger.Tracef("returning build %d position in queue %s", build.GetID(), build.GetHost())

	zRankCmd := c.Redis.ZRank(ctx, build.GetHost(), strconv.FormatInt(build.GetID(), 10))
	if zRankCmd.Err() != nil {
		// the item could be out of the queue so just return 0 on errors
		return 0
	}

	return zRankCmd.Val()
}
