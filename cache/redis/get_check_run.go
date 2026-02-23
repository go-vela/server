// SPDX-License-Identifier: Apache-2.0

package redis

import (
	"context"
	"encoding/json"
	"strconv"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/cache/models"
	"github.com/go-vela/server/constants"
)

// GetCheckRuns retrieves the check run info from Redis for a given build ID.
func (c *Client) GetCheckRuns(ctx context.Context, build *api.Build) ([]models.CheckRun, error) {
	key := "check_run:" + strconv.FormatInt(build.GetID(), 10)

	var (
		checkRunBytes []byte
		err           error
	)

	if isFinished(build.GetStatus()) {
		checkRunBytes, err = c.Redis.GetDel(ctx, key).Bytes()
	} else {
		checkRunBytes, err = c.Redis.Get(ctx, key).Bytes()
	}

	if err != nil {
		return nil, err
	}

	checkRuns := []models.CheckRun{}

	err = json.Unmarshal(checkRunBytes, &checkRuns)
	if err != nil {
		return nil, err
	}

	return checkRuns, nil
}

// GetStepCheckRuns retrieves the check run info from Redis for a given step ID.
func (c *Client) GetStepCheckRuns(ctx context.Context, step *api.Step) ([]models.CheckRun, error) {
	key := "step_check_run:" + strconv.FormatInt(step.GetID(), 10)

	var (
		checkRunBytes []byte
		err           error
	)

	if isFinished(step.GetStatus()) {
		checkRunBytes, err = c.Redis.GetDel(ctx, key).Bytes()
	} else {
		checkRunBytes, err = c.Redis.Get(ctx, key).Bytes()
	}

	if err != nil {
		return nil, err
	}

	checkRuns := []models.CheckRun{}

	err = json.Unmarshal(checkRunBytes, &checkRuns)
	if err != nil {
		return nil, err
	}

	return checkRuns, nil
}

// isFinished is a helper function for determining if a build or step status is in a finished state.
func isFinished(status string) bool {
	return status == constants.StatusSuccess ||
		status == constants.StatusFailure ||
		status == constants.StatusSkipped ||
		status == constants.StatusCanceled ||
		status == constants.StatusError ||
		status == constants.StatusKilled
}
