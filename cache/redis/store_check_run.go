// SPDX-License-Identifier: Apache-2.0

package redis

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/go-vela/server/cache/models"
)

// StoreCheckRuns stores the check runs in Redis with a TTL.
func (c *Client) StoreCheckRuns(ctx context.Context, buildID int64, checkRuns []models.CheckRun, timeout int32) error {
	// set TTL based on repo approval timeout (should be deleted at end of build each time)
	ttl := time.Hour * 24 * time.Duration(timeout+1)

	checkRunBytes, err := json.Marshal(checkRuns)
	if err != nil {
		return err
	}

	key := "check_run:" + strconv.FormatInt(buildID, 10)

	// store a small marker value (or metadata JSON if needed)
	err = c.Redis.Set(ctx, key, checkRunBytes, ttl).Err()
	if err != nil {
		return err
	}

	return nil
}

// StoreStepCheckRuns stores the check runs in Redis with a TTL.
func (c *Client) StoreStepCheckRuns(ctx context.Context, stepID int64, checkRuns []models.CheckRun, timeout int32) error {
	// set TTL based on repo approval timeout (should be deleted at end of build each time)
	ttl := time.Hour * 24 * time.Duration(timeout+1)

	checkRunBytes, err := json.Marshal(checkRuns)
	if err != nil {
		return err
	}

	key := "step_check_run:" + strconv.FormatInt(stepID, 10)

	// store a small marker value (or metadata JSON if needed)
	err = c.Redis.Set(ctx, key, checkRunBytes, ttl).Err()
	if err != nil {
		return err
	}

	return nil
}
