// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/queue"
	"github.com/go-vela/types"
	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"
	"github.com/sirupsen/logrus"
)

// PublishToQueue is a helper function that pushes the build executable to the database
// and publishes a queue item (build, repo, user) to the queue.
func PublishToQueue(ctx context.Context, queue queue.Service, db database.Interface, p *pipeline.Build, b *library.Build, r *library.Repo, u *library.User) {
	// marshal pipeline build into byte data to add to the build executable object
	byteExecutable, err := json.Marshal(p)
	if err != nil {
		logrus.Errorf("Failed to marshal build executable %d for %s: %v", b.GetNumber(), r.GetFullName(), err)

		// error out the build
		CleanBuild(ctx, db, b, nil, nil, err)

		return
	}

	// create build executable to push to database
	bExecutable := new(library.BuildExecutable)
	bExecutable.SetBuildID(b.GetID())
	bExecutable.SetData(byteExecutable)

	// send database call to create a build executable
	err = db.CreateBuildExecutable(ctx, bExecutable)
	if err != nil {
		logrus.Errorf("Failed to publish build executable to database %d for %s: %v", b.GetNumber(), r.GetFullName(), err)

		// error out the build
		CleanBuild(ctx, db, b, nil, nil, err)

		return
	}

	// convert build, repo, and user into queue item
	item := types.ToItem(b, r, u)

	logrus.Infof("Converting queue item to json for build %d for %s", b.GetNumber(), r.GetFullName())

	byteItem, err := json.Marshal(item)
	if err != nil {
		logrus.Errorf("Failed to convert item to json for build %d for %s: %v", b.GetNumber(), r.GetFullName(), err)

		// error out the build
		CleanBuild(ctx, db, b, nil, nil, err)

		return
	}

	logrus.Infof("Establishing route for build %d for %s", b.GetNumber(), r.GetFullName())

	// determine the route on which to publish the queue item
	route, err := queue.Route(&p.Worker)
	if err != nil {
		logrus.Errorf("unable to set route for build %d for %s: %v", b.GetNumber(), r.GetFullName(), err)

		// error out the build
		CleanBuild(ctx, db, b, nil, nil, err)

		return
	}

	logrus.Infof("Publishing item for build %d for %s to queue %s", b.GetNumber(), r.GetFullName(), route)

	// push item on to the queue
	err = queue.Push(context.Background(), route, byteItem)
	if err != nil {
		logrus.Errorf("Retrying; Failed to publish build %d for %s: %v", b.GetNumber(), r.GetFullName(), err)

		err = queue.Push(context.Background(), route, byteItem)
		if err != nil {
			logrus.Errorf("Failed to publish build %d for %s: %v", b.GetNumber(), r.GetFullName(), err)

			// error out the build
			CleanBuild(ctx, db, b, nil, nil, err)

			return
		}
	}

	// update fields in build object
	b.SetEnqueued(time.Now().UTC().Unix())

	// update the build in the db to reflect the time it was enqueued
	_, err = db.UpdateBuild(ctx, b)
	if err != nil {
		logrus.Errorf("Failed to update build %d during publish to queue for %s: %v", b.GetNumber(), r.GetFullName(), err)
	}
}
