// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"
	"encoding/json"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/queue"
	"github.com/go-vela/server/queue/models"
)

// Enqueue is a helper function that pushes a queue item (build, repo, user) to the queue.
func Enqueue(ctx context.Context, queue queue.Service, db database.Interface, item *models.Item, route string) {
	logger := logrus.WithFields(logrus.Fields{
		"build":    item.Build.GetNumber(),
		"build_id": item.Build.GetID(),
		"org":      item.Build.GetRepo().GetOrg(),
		"repo":     item.Build.GetRepo().GetName(),
		"repo_id":  item.Build.GetRepo().GetID(),
	})

	logger.Debug("converting queue item to json")

	byteItem, err := json.Marshal(item)
	if err != nil {
		logger.Errorf("failed to convert item to json: %v", err)

		// error out the build
		CleanBuild(ctx, db, item.Build, nil, nil, err)

		return
	}

	logger.Debugf("pushing item for build to queue route %#q", route)

	// push item on to the queue
	err = queue.Push(context.Background(), route, byteItem)
	if err != nil {
		logger.Errorf("retrying; failed to publish build: %v", err)

		err = queue.Push(context.Background(), route, byteItem)
		if err != nil {
			logger.Errorf("failed to publish build: %v", err)

			// error out the build
			CleanBuild(ctx, db, item.Build, nil, nil, err)

			return
		}
	}

	// update fields in build object
	item.Build.SetEnqueued(time.Now().UTC().Unix())

	// update the build in the db to reflect the time it was enqueued
	_, err = db.UpdateBuild(ctx, item.Build)
	if err != nil {
		logger.Errorf("failed to update build during publish to queue: %v", err)
	}

	logger.Info("updated build as enqueued")
}
