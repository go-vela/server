// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package build

import (
	"fmt"
	"time"

	"github.com/go-vela/server/api/service"
	"github.com/go-vela/server/api/step"
	"github.com/go-vela/server/database"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"
	"github.com/sirupsen/logrus"
)

// PlanBuild is a helper function to plan the build for
// execution. This creates all resources, like steps
// and services, for the build in the configured backend.
// TODO:
// - return build and error.
func PlanBuild(database database.Interface, p *pipeline.Build, b *library.Build, r *library.Repo) error {
	// update fields in build object
	b.SetCreated(time.Now().UTC().Unix())

	// send API call to create the build
	// TODO: return created build and error instead of just error
	err := database.CreateBuild(b)
	if err != nil {
		// clean up the objects from the pipeline in the database
		// TODO:
		// - return build in CreateBuild
		// - even if it was created, we need to get the new build id
		//   otherwise it will be 0, which attempts to INSERT instead
		//   of UPDATE-ing the existing build - which results in
		//   a constraint error (repo_id, number)
		// - do we want to update the build or just delete it?
		cleanBuild(database, b, nil, nil, err)

		return fmt.Errorf("unable to create new build for %s: %w", r.GetFullName(), err)
	}

	// send API call to capture the created build
	// TODO: this can be dropped once we return
	// the created build above
	b, err = database.GetBuild(b.GetNumber(), r)
	if err != nil {
		return fmt.Errorf("unable to get new build for %s: %w", r.GetFullName(), err)
	}

	// plan all services for the build
	services, err := service.PlanServices(database, p, b)
	if err != nil {
		// clean up the objects from the pipeline in the database
		cleanBuild(database, b, services, nil, err)

		return err
	}

	// plan all steps for the build
	steps, err := step.PlanSteps(database, p, b)
	if err != nil {
		// clean up the objects from the pipeline in the database
		cleanBuild(database, b, services, steps, err)

		return err
	}

	return nil
}

// cleanBuild is a helper function to kill the build
// without execution. This will kill all resources,
// like steps and services, for the build in the
// configured backend.
func cleanBuild(database database.Interface, b *library.Build, services []*library.Service, steps []*library.Step, e error) {
	// update fields in build object
	b.SetError(fmt.Sprintf("unable to publish to queue: %s", e.Error()))
	b.SetStatus(constants.StatusError)
	b.SetFinished(time.Now().UTC().Unix())

	// send API call to update the build
	err := database.UpdateBuild(b)
	if err != nil {
		logrus.Errorf("unable to kill build %d: %v", b.GetNumber(), err)
	}

	for _, s := range services {
		// update fields in service object
		s.SetStatus(constants.StatusKilled)
		s.SetFinished(time.Now().UTC().Unix())

		// send API call to update the service
		err := database.UpdateService(s)
		if err != nil {
			logrus.Errorf("unable to kill service %s for build %d: %v", s.GetName(), b.GetNumber(), err)
		}
	}

	for _, s := range steps {
		// update fields in step object
		s.SetStatus(constants.StatusKilled)
		s.SetFinished(time.Now().UTC().Unix())

		// send API call to update the step
		err := database.UpdateStep(s)
		if err != nil {
			logrus.Errorf("unable to kill step %s for build %d: %v", s.GetName(), b.GetNumber(), err)
		}
	}
}
