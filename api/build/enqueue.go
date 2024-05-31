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
	logrus.Debugf("Converting queue item to json for build %d for %s", item.Build.GetNumber(), item.Build.GetRepo().GetFullName())

	byteItem, err := json.Marshal(item)
	if err != nil {
		logrus.Errorf("Failed to convert item to json for build %d for %s: %v", item.Build.GetNumber(), item.Build.GetRepo().GetFullName(), err)

		// error out the build
		CleanBuild(ctx, db, item.Build, nil, nil, err)

		return
	}

	logrus.Debugf("Pushing item for build %d for %s to queue route %s", item.Build.GetNumber(), item.Build.GetRepo().GetFullName(), route)

	// push item on to the queue
	err = queue.Push(context.Background(), route, byteItem)
	if err != nil {
		logrus.Errorf("Retrying; Failed to publish build %d for %s: %v", item.Build.GetNumber(), item.Build.GetRepo().GetFullName(), err)

		err = queue.Push(context.Background(), route, byteItem)
		if err != nil {
			logrus.Errorf("Failed to publish build %d for %s: %v", item.Build.GetNumber(), item.Build.GetRepo().GetFullName(), err)

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
		logrus.Errorf("Failed to update build %d during publish to queue for %s: %v", item.Build.GetNumber(), item.Build.GetRepo().GetFullName(), err)
	}
}
