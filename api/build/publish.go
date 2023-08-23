// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

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

// PublishToQueue is a helper function that creates
// a build item and publishes it to the queue.
func PublishToQueue(ctx context.Context, queue queue.Service, db database.Interface, p *pipeline.Build, b *library.Build, r *library.Repo, u *library.User) {
	byteExecutable, err := json.Marshal(p)
	if err != nil {
		logrus.Errorf("Failed to marshal build executable %d for %s: %v", b.GetNumber(), r.GetFullName(), err)

		// error out the build
		CleanBuild(ctx, db, b, nil, nil, err)

		return
	}

	bExecutable := new(library.BuildExecutable)
	bExecutable.SetBuildID(b.GetID())
	bExecutable.SetData(byteExecutable)

	err = db.CreateBuildExecutable(bExecutable)
	if err != nil {
		logrus.Errorf("Failed to publish build executable to database %d for %s: %v", b.GetNumber(), r.GetFullName(), err)

		// error out the build
		CleanBuild(ctx, db, b, nil, nil, err)

		return
	}

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

	route, err := queue.Route(&p.Worker)
	if err != nil {
		logrus.Errorf("unable to set route for build %d for %s: %v", b.GetNumber(), r.GetFullName(), err)

		// error out the build
		CleanBuild(ctx, db, b, nil, nil, err)

		return
	}

	logrus.Infof("Publishing item for build %d for %s to queue %s", b.GetNumber(), r.GetFullName(), route)

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
