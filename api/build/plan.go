// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/api/service"
	"github.com/go-vela/server/api/step"
	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/scm"
)

// PlanBuild is a helper function to plan the build for
// execution. This creates all resources, like steps
// and services, for the build.
// TODO:
// - return build and error.
func PlanBuild(ctx context.Context, database database.Interface, scm scm.Service, p *pipeline.Build, b *types.Build, r *types.Repo) error {
	// update fields in build object
	b.SetCreated(time.Now().UTC().Unix())

	// send API call to create the build
	// TODO: return created build and error instead of just error
	b, err := database.CreateBuild(ctx, b)
	if err != nil {
		// clean up the objects from the pipeline in the database
		// TODO:
		// - even if it was created, we need to get the new build id
		//   otherwise it will be 0, which attempts to INSERT instead
		//   of UPDATE-ing the existing build - which results in
		//   a constraint error (repo_id, number)
		// - do we want to update the build or just delete it?
		CleanBuild(ctx, database, b, nil, nil, err)

		return fmt.Errorf("unable to create new build for %s: %w", r.GetFullName(), err)
	}

	logrus.WithFields(logrus.Fields{
		"build":    b.GetNumber(),
		"build_id": b.GetID(),
	}).Info("build created")

	// plan all services for the build
	services, err := service.PlanServices(ctx, database, p, b)
	if err != nil {
		// clean up the objects from the pipeline in the database
		CleanBuild(ctx, database, b, services, nil, err)

		return err
	}

	// plan all steps for the build
	steps, err := step.PlanSteps(ctx, database, scm, p, b)
	if err != nil {
		// clean up the objects from the pipeline in the database
		CleanBuild(ctx, database, b, services, steps, err)

		return err
	}

	return nil
}
